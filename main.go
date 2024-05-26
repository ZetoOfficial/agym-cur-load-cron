package main

import (
	"context"
	"github.com/ZetoOfficial/agym-cur-load-cron/config"
	"github.com/ZetoOfficial/agym-cur-load-cron/internal/clients"
	"github.com/ZetoOfficial/agym-cur-load-cron/internal/services"
	"github.com/ZetoOfficial/agym-cur-load-cron/internal/storage/postgres"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	config.ConfigureLogger()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		logrus.Print("Received shutdown signal")
		cancel()
	}()

	cfg, err := config.Load()
	if err != nil {
		logrus.Fatalf("could not load config: %v", err)
	}

	mobiFitnessClient := clients.NewMobiFitnessApi(cfg.MobiFitness)

	db, err := postgres.NewStorage(cfg.Postgres.ConnStr())
	if err != nil {
		logrus.Fatalf("could not open database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			logrus.Printf("error closing database: %v", err)
		}
	}()

	service := services.NewClub(mobiFitnessClient, db)

	if err := service.StartCron(ctx); err != nil {
		logrus.Fatalf("error saving clubs load: %v", err)
	}

	<-ctx.Done()
	logrus.Print("Service shut down gracefully")
}
