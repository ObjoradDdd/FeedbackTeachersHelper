package storage

import (
	"fmt"

	"github.com/ObjoradDdd/FeedbackTeachersHelper/internal/models"
)

func (s *Storage) CreateGroup(group *models.Group) (int, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return 0, fmt.Errorf("Error starting transaction: %w", err)
	}
	defer tx.Rollback()

	var groupID int
	query := `INSERT INTO groups (name) VALUES ($1) RETURNING id`

	if err := tx.QueryRow(query, group.Name).Scan(&groupID); err != nil {
		return 0, fmt.Errorf("Error creating group: %w", err)
	}

	if err := s.addStudentsTx(tx, groupID, group.Students); err != nil {
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("Error committing transaction: %w", err)
	}

	return groupID, nil
}
