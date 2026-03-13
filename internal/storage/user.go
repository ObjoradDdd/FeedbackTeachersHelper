package storage

import (
	"database/sql"
	"fmt"
)

func (s *Storage) ensureUserExists(userID int) error {
	query := `INSERT INTO users (id) VALUES ($1) ON CONFLICT (id) DO NOTHING`
	_, err := s.db.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("failed to ensure user exists: %w", err)
	}
	return nil
}

func (s *Storage) DeleteUserById(id int) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := s.db.Exec(query, id)
	return err
}

func (s *Storage) AddApiKey(userID int, apiKey string) error {
	if err := s.ensureUserExists(userID); err != nil {
		return err
	}

	query := `
		INSERT INTO users (id, api_key)
		VALUES ($1, $2)
		ON CONFLICT (id)
		DO UPDATE SET api_key = EXCLUDED.api_key
	`
	_, err := s.db.Exec(query, userID, apiKey)
	return err
}

func (s *Storage) DeleteApiKey(userID int) error {
	if err := s.ensureUserExists(userID); err != nil {
		return err
	}

	query := `UPDATE users SET api_key = NULL WHERE id = $1`
	_, err := s.db.Exec(query, userID)
	return err
}

func (s *Storage) GetApiKey(userID int) (string, error) {
	if err := s.ensureUserExists(userID); err != nil {
		return "", err
	}

	query := `SELECT api_key FROM users WHERE id = $1`
	var apiKey sql.NullString

	if err := s.db.QueryRow(query, userID).Scan(&apiKey); err != nil {
		return "", err
	}

	if !apiKey.Valid || apiKey.String == "" {
		return "", fmt.Errorf("api key is not set for user %d", userID)
	}

	return apiKey.String, nil
}
