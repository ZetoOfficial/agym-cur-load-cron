package main

import (
	"context"
	"github.com/ZetoOfficial/agym-cur-load-cron/config"
	"github.com/ZetoOfficial/agym-cur-load-cron/internal/clients"
	"github.com/ZetoOfficial/agym-cur-load-cron/internal/services"
	"github.com/ZetoOfficial/agym-cur-load-cron/internal/storage/sqlite"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		cancel()
	}()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("could not load config: %v", err)
	}

	mobiFitnessClient := clients.NewMobiFitnessApi(cfg.MobiFitness)

	db, err := sqlite.NewStorage("test.db")
	if err != nil {
		log.Fatalf("could not open database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("error closing database: %v", err)
		}
	}()

	service := services.NewClub(mobiFitnessClient, db)

	if err := service.SaveClubsLoad(ctx); err != nil {
		log.Fatalf("error saving clubs load: %v", err)
	}
}
