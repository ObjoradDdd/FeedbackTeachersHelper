package storage

import (
	"database/sql"
	"fmt"

	"github.com/ObjoradDdd/FeedbackTeachersHelper/internal/models"
)

func (s *Storage) addStudentsTx(tx *sql.Tx, groupID int, students []models.Student) error {
	query := `INSERT INTO students (name, group_id) VALUES ($1, $2)`

	for _, student := range students {
		_, err := tx.Exec(query, student.Name, groupID)
		if err != nil {
			return fmt.Errorf("Error adding student %s: %w", student.Name, err)
		}
	}
	return nil
}

func (s *Storage) AddStudentsToExistingGroup(groupID int, students []models.Student) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("Error starting transaction: %w", err)
	}
	defer tx.Rollback()

	if err := s.addStudentsTx(tx, groupID, students); err != nil {
		return err
	}

	return tx.Commit()
}
