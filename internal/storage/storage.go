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

func (s *Storage) InitTables() error {
	query := `

	DROP TABLE IF EXISTS students CASCADE;
	DROP TABLE IF EXISTS groups CASCADE;
	DROP TABLE IF EXISTS teachers CASCADE;
	DROP TABLE IF EXISTS tags CASCADE;

	CREATE TABLE IF NOT EXISTS teachers (
		id SERIAL PRIMARY KEY,
		api_key VARCHAR(255),
		login VARCHAR(255) NOT NULL UNIQUE,	
		hash VARCHAR(255) NOT NULL
	);

	CREATE TABLE IF NOT EXISTS tags (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,	
		meaning VARCHAR(255) NOT NULL,
		teacher_id INTEGER REFERENCES teachers(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS groups (
		id SERIAL PRIMARY KEY,
		teacher_id INTEGER REFERENCES teachers(id) ON DELETE CASCADE,
		name VARCHAR(255) NOT NULL
	);

	CREATE TABLE IF NOT EXISTS students (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		group_id INTEGER REFERENCES groups(id) ON DELETE CASCADE
	);
	`

	_, err := s.db.Exec(query)
	if err != nil {
		return fmt.Errorf("Error initializing tables: %w", err)
	}

	return nil
}
