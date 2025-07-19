package main

import (
	"context"
	"fmt"
	"goautomation/db"
	"goautomation/monitor"
	"log"
	"time"
)

func mains() {
	ctx := context.Background()
	database, err := db.NewDatabase("./data/db.sqlite")

	s := database.Queries
	if err != nil {
		log.Fatalf("❌ Could not initialize database: %v", err)
		return
	}
	defer func(database *db.Database) {
		err := database.Close()
		if err != nil {
			log.Fatalf("❌ Could not close database: %v", err)
			return
		}
	}(database)
	count := 0
	for {
		monitor.Watch(ctx, s, "raiboss1")
		fmt.Println("checked raiboss1")
		monitor.Watch(ctx, s, "raiboss2")
		fmt.Println("checked raiboss2")
		// monitor.Watch(ctx, s, "mylaptop")
		// fmt.Println("checked mylaptop")
		_ = s.DeleteAfter30min(ctx)
		fmt.Println("🔄 Monitoring... ", count)
		time.Sleep(time.Second * 30)
		count++
	}

}

func main() {
	// rmq, err := monitor.ClientCreate("mypc")
	// if err != nil {
	// 	fmt.Println("Error creating client:", err)
	// 	return
	// }
	// fmt.Println(
	// 	rmq.GetPixelColor(0.72, 0.80),
	// )
	mains()
	// monitor.ClickEvent("raiboss1", 0.55, 0.70, nil, nil)
}
