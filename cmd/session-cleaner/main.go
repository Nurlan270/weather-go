package main

import (
	"context"
	"log"
	"time"

	"github.com/Nurlan270/weather-go/internal/config"
	database "github.com/Nurlan270/weather-go/internal/db"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	dbCfg := cfg.DB

	db, err := database.Connect(dbCfg)
	if err != nil {
		log.Fatal(err)
	}

	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Println("Starting cleaner...")

	for {
		select {
		case <-ticker.C:
			r, err := db.Exec("DELETE FROM sessions WHERE expires_at < now()")
			if err != nil {
				log.Printf("failed to delete expired sesssions: %v", err)
			}

			count, _ := r.RowsAffected()
			log.Printf("%v row(s) was deleted", count)
		case <-ctx.Done():
			log.Println("Cleaner terminated")
			return
		}
	}
}
