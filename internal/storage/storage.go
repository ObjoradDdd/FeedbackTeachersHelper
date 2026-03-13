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

	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY,
		api_key VARCHAR(255)
	);

	DO $$
	BEGIN
		IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'teachers') THEN
			INSERT INTO users (id, api_key)
			SELECT id, api_key FROM teachers
			ON CONFLICT (id) DO UPDATE SET api_key = EXCLUDED.api_key;
		END IF;
	END $$;

	CREATE TABLE IF NOT EXISTS tags (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,	
		meaning VARCHAR(255) NOT NULL,
		user_id INTEGER REFERENCES users(id) ON DELETE CASCADE
	);

	ALTER TABLE tags ADD COLUMN IF NOT EXISTS user_id INTEGER;

	DO $$
	BEGIN
		IF EXISTS (
			SELECT 1
			FROM information_schema.columns
			WHERE table_name = 'tags' AND column_name = 'teacher_id'
		) THEN
			UPDATE tags SET user_id = teacher_id WHERE user_id IS NULL;
		END IF;
	END $$;

	CREATE TABLE IF NOT EXISTS groups (
		id SERIAL PRIMARY KEY,
		user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
		name VARCHAR(255) NOT NULL
	);

	ALTER TABLE groups ADD COLUMN IF NOT EXISTS user_id INTEGER;

	DO $$
	BEGIN
		IF EXISTS (
			SELECT 1
			FROM information_schema.columns
			WHERE table_name = 'groups' AND column_name = 'teacher_id'
		) THEN
			UPDATE groups SET user_id = teacher_id WHERE user_id IS NULL;
		END IF;
	END $$;

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
