package storage

import (
	"context"
	"fmt"

	"github.com/ObjoradDdd/FeedbackTeachersHelper/internal/models"
)

func (s *Storage) CreateTag(ctx context.Context, tag *models.Tag, userID int) (int, error) {
	if err := s.ensureUserExists(ctx, userID); err != nil {
		return 0, err
	}

	query := `
		INSERT INTO tags (name, meaning, user_id) VALUES ($1, $2, $3) RETURNING id
	`
	var tagId int

	if err := s.db.QueryRowContext(ctx, query, tag.Name, tag.Meaning, userID).Scan(&tagId); err != nil {
		return 0, fmt.Errorf("Error adding tag %s: %w", tag.Name, err)
	}

	return tagId, nil
}

func (s *Storage) GetUserTags(ctx context.Context, userID int) ([]models.Tag, error) {
	if err := s.ensureUserExists(ctx, userID); err != nil {
		return nil, err
	}

	query := `SELECT id, name, meaning FROM tags WHERE user_id = $1`

	rows, err := s.db.QueryContext(ctx, query, userID)
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

func (s *Storage) DeleteTag(ctx context.Context, id int, userID int) error {
	if err := s.ensureUserExists(ctx, userID); err != nil {
		return err
	}

	query := `DELETE FROM tags WHERE id = $1 AND user_id = $2`

	if _, err := s.db.ExecContext(ctx, query, id, userID); err != nil {
		return fmt.Errorf("Error deleting tag: %w", err)
	}

	return nil
}

func (s *Storage) UpdateTag(ctx context.Context, tag *models.Tag, userID int) error {
	if err := s.ensureUserExists(ctx, userID); err != nil {
		return err
	}

	query := `UPDATE tags SET name = $1, meaning = $2 WHERE id = $3 AND user_id = $4`

	result, err := s.db.ExecContext(ctx, query, tag.Name, tag.Meaning, tag.Id, userID)

	if err != nil {
		return fmt.Errorf("Error updating tag: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Error checking update result: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("No tag updated, check if tag exists and belongs to user ID %d", userID)
	}

	return nil
}
