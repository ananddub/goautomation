package monitor

import (
	"context"
	"fmt"
	"goautomation/client"
	"goautomation/db/sqlc"
	"goautomation/process"
	"goautomation/sound"
)

const (
	// width  = 0.72
	// height = 0.92

	width  = 0.55
	height = 0.70
	url    = "amqp://34.56.24.250:5672"
)

func UpdateTable(name string, s *sqlc.Queries, ctx context.Context) {
	rmq, err := ClientCreate(name)
	if err != nil {
		return
	}
	data, err := rmq.GetPixelColor(width, height)
	if err != nil {
		fmt.Printf("❌ Could not get pixel color: %v", err)
		for i := 0; i < 5; i++ {
			sound.PlayBeep()
		}
		return
	}
	if len(data.Windows)%4 != 0 {
		fmt.Printf("invalid number of windows received")
		for i := 0; i < 5; i++ {
			sound.PlayBeep()
		}

		return
	}
	err = process.InsertWindow(rmq, s, ctx)
	if err != nil {
		return
	}
}

func Watch(ctx context.Context, s *sqlc.Queries, name string) {
	UpdateTable(name, s, ctx)
	changes, err := s.CheckRecentChanges5min(ctx)
	if err != nil {
		return
	}
	flag := false
	for _, v := range changes {
		if v.TotalPixel > 7 && v.TotalUniquePixel <= 2 {
			fmt.Println(v.Name, v.Title)
			flag = true
		}
	}
	if flag {
		for i := 0; i < 3; i++ {
			sound.PlayBeep()
		}
	}
}
func ClientCreate(user string) (*client.AutomationRPCClient, error) {
	rmq, err := client.NewAutomationRPCClient(url, user)
	if err != nil {
		return nil, err
	}
	return rmq, nil
}

func ClickEvent(name string, width float64, height float64, x *float64, y *float64) {
	rmq, err := ClientCreate(name)
	if err != nil || rmq == nil {
		return
	}
	w := width
	h := height
	emptyString := "20"
	db, err := rmq.MoveMouse(&emptyString, x, y, &w, &h)
	fmt.Println(
		db)
	if err != nil {
		fmt.Printf("❌ Could not trigger click event: %v", err)
	}
}
