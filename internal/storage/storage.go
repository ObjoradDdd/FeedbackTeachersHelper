package storage

import (
	"database/sql"
	"fmt"
	"log/slog"

	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func New(connString string) (*Storage, error) {
	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, fmt.Errorf("DB connection error: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("DB ping error: %w", err)
	}

	slog.Info("DB connection successful", "driver", "postgres")

	return &Storage{db: db}, nil
}

func (s *Storage) Close() error {
	slog.Info("Closing DB connection")
	return s.db.Close()
}
