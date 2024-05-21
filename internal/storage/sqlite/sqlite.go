package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/ZetoOfficial/agym-cur-load-cron/internal/models"
	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(dbPath string) (*Storage, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := initTables(db); err != nil {
		return nil, fmt.Errorf("failed to initialize tables: %w", err)
	}

	log.Printf("database opened at %s", dbPath)
	return &Storage{db: db}, nil
}

func (s *Storage) Close() error {
	if err := s.db.Close(); err != nil {
		return fmt.Errorf("failed to close database: %w", err)
	}
	log.Println("database closed")
	return nil
}

func initTables(db *sql.DB) error {
	createClubLoadsTable := `
	CREATE TABLE IF NOT EXISTS club_loads (
		club_id INTEGER,
		club_title TEXT,
		load INTEGER,
		created_at TIMESTAMP
	)`
	_, err := db.Exec(createClubLoadsTable)
	if err != nil {
		return fmt.Errorf("failed to create club_loads table: %w", err)
	}
	log.Println("table club_loads created or already exists")
	return nil
}

func (s *Storage) SaveClubLoad(ctx context.Context, load *models.ClubInfoResponse) error {
	loc, _ := time.LoadLocation("Asia/Yekaterinburg")
	utcTime := time.Now().In(loc)

	query := `INSERT INTO club_loads (club_id, club_title, load, created_at) VALUES (?, ?, ?, ?)`
	stmt, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			log.Printf("failed to close statement: %v", err)
		}
	}(stmt)

	if _, err := stmt.ExecContext(ctx, load.Id, load.Title, load.CurrentLoad, utcTime); err != nil {
		return fmt.Errorf("failed to execute statement: %w", err)
	}

	log.Printf("club %v (%v) with load %v saved at %v", load.Title, load.Id, load.CurrentLoad, utcTime)
	return nil
}
