package storage

import (
	"fmt"

	"github.com/ObjoradDdd/FeedbackTeachersHelper/internal/models"
)

func (s *Storage) AddTag(tag *models.Tag) (int, error) {
	query := `INSERT INTO tags (name, meaning, is_bad, teacher_id) VALUES ($1, $2, $3, $4) RETURNING id`
	var tagID int

	if err := s.db.QueryRow(query, tag.Name, tag.Meaning, tag.IsBad, tag.TeacherID).Scan(&tagID); err != nil {
		return 0, fmt.Errorf("Error adding tag %s: %w", tag.Name, err)
	}

	return tagID, nil
}

func (s *Storage) GetTeachersTags(teacherId int) ([]models.Tag, error) {
	query := `SELECT id, name, meaning, is_bad FROM tags WHERE teacher_id = $1`

	rows, err := s.db.Query(query, teacherId)
	if err != nil {
		return nil, fmt.Errorf("Error fetching tags: %w", err)
	}
	defer rows.Close()

	var tags []models.Tag
	for rows.Next() {
		var tag models.Tag
		if err := rows.Scan(&tag.ID, &tag.Name, &tag.Meaning, &tag.IsBad); err != nil {
			return nil, fmt.Errorf("Error scanning tag: %w", err)
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

func (s *Storage) DeleteTag(id int) error {
	query := `DELETE FROM tags WHERE id = $1`

	if _, err := s.db.Exec(query, id); err != nil {
		return fmt.Errorf("Error deleting tag: %w", err)
	}

	return nil
}

func (s *Storage) UpdateTag(tag *models.Tag) error {
	query := `UPDATE tags SET name = $1, meaning = $2, is_bad = $3 WHERE id = $4`

	if _, err := s.db.Exec(query, tag.Name, tag.Meaning, tag.IsBad, tag.ID); err != nil {
		return fmt.Errorf("Error updating tag: %w", err)
	}

	return nil
}
