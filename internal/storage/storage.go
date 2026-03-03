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

	slog.Info("DB succes", "driver", "postgres")

	return &Storage{db: db}, nil
}

func (s *Storage) Close() error {
	slog.Info("DB close")
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
		name VARCHAR(255) NOT NULL,	
		hash VARCHAR(255) NOT NULL
	);

	CREATE TABLE IF NOT EXISTS groups (
		id SERIAL PRIMARY KEY,
		teacher_id INTEGER REFERENCES teachers(id) ON DELETE SET NULL,
		name VARCHAR(255) NOT NULL
	);

	CREATE TABLE IF NOT EXISTS students (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		group_id INTEGER REFERENCES groups(id) ON DELETE CASCADE
	);
	

	CREATE TABLE IF NOT EXISTS tags (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		is_bad BOOLEAN NOT NULL DEFAULT false
	);
	`

	_, err := s.db.Exec(query)
	if err != nil {
		return fmt.Errorf("Error initializing tables: %w", err)
	}

	return nil
}
