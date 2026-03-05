package storage

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/ObjoradDdd/FeedbackTeachersHelper/internal/models"
)

func (s *Storage) CreateStudent(student *models.Student, teacherID int) (int, error) {
	query := `
		INSERT INTO students (name, group_id) 
		SELECT $1, $2 
		WHERE EXISTS (
			SELECT 1 FROM groups WHERE id = $2 AND teacher_id = $3
		)
		RETURNING id
	`
	var studentID int

	err := s.db.QueryRow(query, student.Name, student.GroupID, teacherID).Scan(&studentID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, errors.New("forbidden: target group does not exist or access denied")
		}
		return 0, fmt.Errorf("failed to insert student %s: %w", student.Name, err)
	}

	return studentID, nil
}

func (s *Storage) GetGroupStudents(groupId int, teacherID int) ([]models.Student, error) {
	query := `SELECT id, name FROM students WHERE group_id = $1 AND group_id IN (SELECT id FROM groups WHERE teacher_id = $2)`

	rows, err := s.db.Query(query, groupId, teacherID)
	if err != nil {
		return nil, fmt.Errorf("Error fetching students: %w", err)
	}
	defer rows.Close()

	var students []models.Student
	for rows.Next() {
		student := models.Student{GroupID: groupId}

		if err := rows.Scan(&student.ID, &student.Name); err != nil {
			return nil, fmt.Errorf("Error scanning student: %w", err)
		}
		students = append(students, student)
	}

	return students, nil
}

func (s *Storage) DeleteStudent(id int, teacherID int) error {
	query := `DELETE FROM students WHERE id = $1 AND group_id IN (SELECT id FROM groups WHERE teacher_id = $2)`

	if _, err := s.db.Exec(query, id, teacherID); err != nil {
		return fmt.Errorf("Error deleting student: %w", err)
	}

	return nil
}

func (s *Storage) UpdateStudent(student *models.Student, teacherID int) error {
	query := `
		UPDATE students 
		SET name = $1, group_id = $2 
		WHERE id = $3 
		AND EXISTS (SELECT 1 FROM groups WHERE id = students.group_id AND teacher_id = $4)
		AND EXISTS (SELECT 1 FROM groups WHERE id = $2 AND teacher_id = $4)
	`

	result, err := s.db.Exec(query, student.Name, student.GroupID, student.ID, teacherID)
	if err != nil {
		return fmt.Errorf("Error updating student: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("forbidden: student not found or access denied for current/target group")
	}

	return nil
}
