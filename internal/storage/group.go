package storage

import (
	"fmt"

	"github.com/ObjoradDdd/FeedbackTeachersHelper/internal/models"
)

func (s *Storage) CreateGroup(group *models.Group, teacherId int) (int, error) {
	query := `INSERT INTO groups (teacher_id, name) VALUES ($1, $2) RETURNING id`
	var groupId int

	if err := s.db.QueryRow(query, teacherId, group.Name).Scan(&groupId); err != nil {
		return 0, fmt.Errorf("Error creating group: %w", err)
	}

	return groupId, nil
}

func (s *Storage) GetGroup(id int, teacherId int) (*models.Group, error) {
	var group models.Group

	query := `SELECT id, name FROM groups WHERE id = $1 AND teacher_id = $2`

	if err := s.db.QueryRow(query, id, teacherId).Scan(&group.Id, &group.Name); err != nil {
		return nil, fmt.Errorf("Error fetching group: %w", err)
	}

	return &group, nil
}

func (s *Storage) DeleteGroup(id int, teacherId int) error {
	query := `DELETE FROM groups WHERE id = $1 AND teacher_id = $2`

	if _, err := s.db.Exec(query, id, teacherId); err != nil {
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
		if err := rows.Scan(&group.Id, &group.Name); err != nil {
			return nil, fmt.Errorf("Error scanning group: %w", err)
		}
		groups = append(groups, group)
	}

	return groups, nil
}

func (s *Storage) UpdateGroup(group *models.Group, teacherId int) error {
	query := `UPDATE groups SET name = $1 WHERE id = $2 AND teacher_id = $3`

	if _, err := s.db.Exec(query, group.Name, group.Id, teacherId); err != nil {
		return fmt.Errorf("Error updating group: %w", err)
	}

	return nil

}
