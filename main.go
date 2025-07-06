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
		fmt.Printf("🔄 Monitoring... %d\n", count)
		err = monitor.Watch(ctx, s, "raiboss1")
		err = s.DeleteAfter30min(ctx)
		if err != nil {
			log.Fatalf("❌ Could not delete old records: %v", err)
		}
		time.Sleep(time.Second * 30)
		count++
	}

}
