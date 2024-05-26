package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/sirupsen/logrus"
	"time"

	"github.com/ZetoOfficial/agym-cur-load-cron/internal/models"
	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(connStr string) (*Storage, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := initTables(db); err != nil {
		return nil, fmt.Errorf("failed to initialize tables: %w", err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Close() error {
	if err := s.db.Close(); err != nil {
		return fmt.Errorf("failed to close database: %w", err)
	}
	logrus.Println("database closed")
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
	logrus.Println("table club_loads created or already exists")
	return nil
}

func (s *Storage) SaveClubLoad(ctx context.Context, load *models.ClubInfoResponse) error {
	loc, _ := time.LoadLocation("Asia/Yekaterinburg")
	utcTime := time.Now().In(loc)

	query := `INSERT INTO club_loads (club_id, club_title, load, created_at) VALUES ($1, $2, $3, $4)`
	stmt, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			logrus.Printf("failed to close statement: %v", err)
		}
	}(stmt)

	if _, err := stmt.ExecContext(ctx, load.Id, load.Title, load.CurrentLoad, utcTime); err != nil {
		return fmt.Errorf("failed to execute statement: %w", err)
	}

	logrus.Debugf("club %v (%v) with load %v saved at %v", load.Title, load.Id, load.CurrentLoad, utcTime)
	return nil
}
