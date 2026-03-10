package storage

import (
	"fmt"

	"github.com/ObjoradDdd/FeedbackTeachersHelper/internal/models"
)

func (s *Storage) CreateTag(tag *models.Tag, teacherId int) (int, error) {
	query := `INSERT INTO tags (name, meaning, teacher_id) VALUES ($1, $2, $3) RETURNING id`
	var tagId int

	if err := s.db.QueryRow(query, tag.Name, tag.Meaning, teacherId).Scan(&tagId); err != nil {
		return 0, fmt.Errorf("Error adding tag %s: %w", tag.Name, err)
	}

	return tagId, nil
}

func (s *Storage) GetTeachersTags(teacherId int) ([]models.Tag, error) {
	query := `SELECT id, name, meaning FROM tags WHERE teacher_id = $1`

	rows, err := s.db.Query(query, teacherId)
	if err != nil {
		return nil, fmt.Errorf("Error fetching tags: %w", err)
	}
	defer rows.Close()

	var tags []models.Tag
	for rows.Next() {
		var tag models.Tag
		if err := rows.Scan(&tag.Id, &tag.Name, &tag.Meaning); err != nil {
			return nil, fmt.Errorf("Error scanning tag: %w", err)
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

func (s *Storage) DeleteTag(id int, teacherId int) error {
	query := `DELETE FROM tags WHERE id = $1 AND teacher_id = $2`

	if _, err := s.db.Exec(query, id, teacherId); err != nil {
		return fmt.Errorf("Error deleting tag: %w", err)
	}

	return nil
}

func (s *Storage) UpdateTag(tag *models.Tag, teacherId int) error {
	query := `UPDATE tags SET name = $1, meaning = $2 WHERE id = $3 AND teacher_id = $4`

	result, err := s.db.Exec(query, tag.Name, tag.Meaning, tag.Id, teacherId)

	if err != nil {
		return fmt.Errorf("Error updating tag: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Error checking update result: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("No tag updated, check if tag exists and belongs to teacher Id %d", teacherId)
	}

	return nil
}
