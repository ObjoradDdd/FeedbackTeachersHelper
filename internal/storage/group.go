package storage

import (
	"fmt"

	"github.com/ObjoradDdd/FeedbackTeachersHelper/internal/models"
)

func (s *Storage) CreateGroup(group *models.Group) (int, error) {
	query := `INSERT INTO groups (teacher_id, name) VALUES ($1, $2) RETURNING id`
	var groupID int

	if err := s.db.QueryRow(query, group.TeacherID, group.Name).Scan(&groupID); err != nil {
		return 0, fmt.Errorf("Error creating group: %w", err)
	}

	return groupID, nil
}

func (s *Storage) GetGroup(id int) (*models.Group, error) {
	var group models.Group

	query := `SELECT id, name FROM groups WHERE id = $1`

	if err := s.db.QueryRow(query, id).Scan(&group.ID, &group.Name); err != nil {
		return nil, fmt.Errorf("Error fetching group: %w", err)
	}

	return &group, nil
}

func (s *Storage) DeleteGroup(id int) error {
	query := `DELETE FROM groups WHERE id = $1`

	if _, err := s.db.Exec(query, id); err != nil {
		return fmt.Errorf("Error deleting group: %w", err)
	}

	return nil
}

func (s *Storage) GetTeachersGroups(teacherId int) ([]models.Group, error) {
	query := `SELECT id, name FROM groups WHERE teacher_id = $1`

	rows, err := s.db.Query(query, teacherId)
	if err != nil {
		return nil, fmt.Errorf("Error fetching groups: %w", err)
	}
	defer rows.Close()

	var groups []models.Group
	for rows.Next() {
		var group models.Group
		if err := rows.Scan(&group.ID, &group.Name); err != nil {
			return nil, fmt.Errorf("Error scanning group: %w", err)
		}
		groups = append(groups, group)
	}

	return groups, nil
}

func (s *Storage) UpdateGroup(group *models.Group) error {
	query := `UPDATE groups SET name = $1 WHERE id = $2`

	if _, err := s.db.Exec(query, group.Name, group.ID); err != nil {
		return fmt.Errorf("Error updating group: %w", err)
	}

	return nil

}
