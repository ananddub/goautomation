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
	width  = 0.62
	height = 0.92
	url    = "amqp://34.56.24.250:5672"
)

func Watch(ctx context.Context, s *sqlc.Queries, name string) error {
	rmq, err := clientCreate(name)
	if err != nil {
		return fmt.Errorf("❌ Could not create client: %v", err)
	}
	data, err := rmq.GetPixelColor(width, height)
	if err != nil {
		fmt.Printf("❌ Could not get pixel color: %v", err)
		for i := 0; i < 5; i++ {
			sound.PlayBeep()
		}
		return fmt.Errorf("❌ Could not get pixel color: %v", err)
	}
	if len(data.Windows)%4 != 0 {
		fmt.Printf("invalid number of windows received")
		for i := 0; i < 5; i++ {
			sound.PlayBeep()
		}

		return fmt.Errorf("invalid number of windows received")
	}
	err = process.InsertWindow(rmq, s, ctx)
	if err != nil {
		return err
	}
	changes, err := s.CheckRecentChanges5min(ctx)
	if err != nil {
		return fmt.Errorf("❌ Could not check recent changes: %v", err)
	}
	for _, v := range changes {
		if v.TotalPixel > 8 && v.TotalUniquePixel == 1 {
			fmt.Printf(" %s", v.Title)
			for i := 0; i < 3; i++ {
				sound.PlayBeep()
			}
			return fmt.Errorf("garbar hai bhaiya garbar hai %v", v.TotalUniquePixel)
		}
	}
	return nil
}
func clientCreate(user string) (*client.AutomationRPCClient, error) {
	rmq, err := client.NewAutomationRPCClient(url, user)
	if err != nil {
		return nil, err
	}
	return rmq, nil
}
