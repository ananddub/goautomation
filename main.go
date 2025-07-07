package main

import (
	"context"
	"fmt"
	"goautomation/db"
	"goautomation/monitor"
	"log"
	"time"
)

func main() {
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
		monitor.Watch(ctx, s, "mylaptop")
		fmt.Println("checked mylaptop")
		err = s.DeleteAfter30min(ctx)
		fmt.Println("🔄 Monitoring... ", count)
		time.Sleep(time.Second * 30)
		count++
	}

}
