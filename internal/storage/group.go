package storage

import (
	"fmt"

	"github.com/ObjoradDdd/FeedbackTeachersHelper/internal/models"
)

func (s *Storage) CreateGroup(group *models.Group, userID int) (int, error) {
	if err := s.ensureUserExists(userID); err != nil {
		return 0, err
	}

	query := `
		INSERT INTO groups (user_id, name) VALUES ($1, $2) RETURNING id
	`
	var groupId int

	if err := s.db.QueryRow(query, userID, group.Name).Scan(&groupId); err != nil {
		return 0, fmt.Errorf("Error creating group: %w", err)
	}

	return groupId, nil
}

func (s *Storage) GetGroup(id int, userID int) (*models.Group, error) {
	if err := s.ensureUserExists(userID); err != nil {
		return nil, err
	}

	var group models.Group

	query := `SELECT id, name FROM groups WHERE id = $1 AND user_id = $2`

	if err := s.db.QueryRow(query, id, userID).Scan(&group.Id, &group.Name); err != nil {
		return nil, fmt.Errorf("Error fetching group: %w", err)
	}

	return &group, nil
}

func (s *Storage) DeleteGroup(id int, userID int) error {
	if err := s.ensureUserExists(userID); err != nil {
		return err
	}

	query := `DELETE FROM groups WHERE id = $1 AND user_id = $2`

	if _, err := s.db.Exec(query, id, userID); err != nil {
		return fmt.Errorf("Error deleting group: %w", err)
	}

	return nil
}

func (s *Storage) GetUserGroups(userID int) ([]models.Group, error) {
	if err := s.ensureUserExists(userID); err != nil {
		return nil, err
	}

	query := `SELECT id, name FROM groups WHERE user_id = $1`

	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("Error fetching groups: %w", err)
	}
	defer rows.Close()

	var groups []models.Group
	for rows.Next() {
		var group models.Group
		if err := rows.Scan(&group.Id, &group.Name); err != nil {
			return nil, fmt.Errorf("Error scanning group: %w", err)
		}
		groups = append(groups, group)
	}

	return groups, nil
}

func (s *Storage) UpdateGroup(group *models.Group, userID int) error {
	if err := s.ensureUserExists(userID); err != nil {
		return err
	}

	query := `UPDATE groups SET name = $1 WHERE id = $2 AND user_id = $3`

	if _, err := s.db.Exec(query, group.Name, group.Id, userID); err != nil {
		return fmt.Errorf("Error updating group: %w", err)
	}

	return nil

}
