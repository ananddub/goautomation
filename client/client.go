package client

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/streadway/amqp"
)

const (
	DefaultPath = "amqp://34.56.24.250:5672"
	DefaultUser = "raiboss1"
)

type AutomationRPCClient struct {
	connection    *amqp.Connection
	channel       *amqp.Channel
	User          string
	callbackQueue string
	responses     map[string]chan Response
}

type Request struct {
	Action string                 `json:"action"`
	Params map[string]interface{} `json:"params"`
}

type Response struct {
	Success bool                  `json:"success"`
	Windows map[string]WindowData `json:"windows,omitempty"`
	Error   string                `json:"error,omitempty"`
}

type WindowData struct {
	X   float64 `json:"x"`
	Y   float64 `json:"y"`
	RGB RGB     `json:"rgb"`
}

type RGB struct {
	R int `json:"r"`
	G int `json:"g"`
	B int `json:"b"`
}

func NewAutomationRPCClient(rabbitmqURL string, user string) (*AutomationRPCClient, error) {
	if rabbitmqURL == "" {
		rabbitmqURL = DefaultPath
	}

	conn, err := amqp.Dial(rabbitmqURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %v", err)
	}

	// Declare callback queue
	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare callback queue: %v", err)
	}

	client := &AutomationRPCClient{
		connection:    conn,
		channel:       ch,
		User:          user,
		callbackQueue: q.Name,
		responses:     make(map[string]chan Response),
	}

	// Start consuming responses
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		client.Close()
		return nil, fmt.Errorf("failed to register consumer: %v", err)
	}

	go client.handleResponses(msgs)

	return client, nil
}

func (c *AutomationRPCClient) handleResponses(msgs <-chan amqp.Delivery) {
	for d := range msgs {
		corrID := d.CorrelationId
		if responseChan, exists := c.responses[corrID]; exists {
			var response Response
			if err := json.Unmarshal(d.Body, &response); err != nil {
				log.Printf("Error unmarshaling response: %v", err)
				continue
			}
			responseChan <- response
			delete(c.responses, corrID)
		}
	}
}

func (c *AutomationRPCClient) Call(action string, params map[string]interface{}) (Response, error) {
	if params == nil {
		params = make(map[string]interface{})
	}

	corrID := uuid.New().String()
	responseChan := make(chan Response, 1)
	c.responses[corrID] = responseChan

	request := Request{
		Action: action,
		Params: params,
	}

	body, err := json.Marshal(request)
	if err != nil {
		delete(c.responses, corrID)
		return Response{}, fmt.Errorf("failed to marshal request: %v", err)
	}

	err = c.channel.Publish(
		"",     // exchange
		c.User, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType:   "application/json",
			CorrelationId: corrID,
			ReplyTo:       c.callbackQueue,
			Body:          body,
		},
	)
	if err != nil {
		delete(c.responses, corrID)
		return Response{}, fmt.Errorf("failed to publish message: %v", err)
	}

	// Wait for response with timeout
	select {
	case response := <-responseChan:
		return response, nil
	case <-time.After(30 * time.Second):
		delete(c.responses, corrID)
		return Response{}, fmt.Errorf("timeout waiting for response")
	}
}

func (c *AutomationRPCClient) MoveMouse(idParam *string, x *float64, y *float64, width *float64, height *float64) (Response, error) {
	params := make(map[string]interface{})
	if idParam != nil {
		params["id"] = *idParam
	}
	if x != nil {
		params["x"] = *x
	}
	if y != nil {
		params["y"] = *y
	}
	if width != nil {
		params["width"] = *width
	}
	if height != nil {
		params["height"] = *height
	}

	return c.Call("move", params)
}

func (c *AutomationRPCClient) ClickMouse(idParam *string, x *float64, y *float64, width *float64, height *float64) (Response, error) {
	params := make(map[string]interface{})
	if idParam != nil {
		params["id"] = *idParam
	}
	if x != nil {
		params["x"] = *x
	}
	if y != nil {
		params["y"] = *y
	}
	if width != nil {
		params["width"] = *width
	}
	if height != nil {
		params["height"] = *height
	}

	return c.Call("click", params)
}

func (c *AutomationRPCClient) GetPixelColor(width float64, height float64) (Response, error) {
	params := make(map[string]interface{})
	params["id"] = uuid.New().String()
	params["width"] = &width
	params["height"] = &height

	return c.Call("color", params)
}

func (c *AutomationRPCClient) Close() {
	if c.channel != nil {
		c.channel.Close()
	}
	if c.connection != nil {
		c.connection.Close()
	}
}

// Helper functions for easier parameter passing
func StringPtr(s string) *string {
	return &s
}

func Float64Ptr(f float64) *float64 {
	return &f
}

func (s *AutomationRPCClient) GetClientUser() string {
	return s.User
}
