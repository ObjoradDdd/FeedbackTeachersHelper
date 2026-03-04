package storage

import (
	"fmt"

	"github.com/ObjoradDdd/FeedbackTeachersHelper/internal/models"
)

func (s *Storage) AddStudent(student *models.Student) (int, error) {
	query := `INSERT INTO students (name, group_id) VALUES ($1, $2) RETURNING id`
	var studentID int

	if err := s.db.QueryRow(query, student.Name, student.GroupID).Scan(&studentID); err != nil {
		return 0, fmt.Errorf("Error adding student %s: %w", student.Name, err)
	}

	return studentID, nil
}

func (s *Storage) GetGroupStudents(groupId int) ([]models.Student, error) {
	query := `SELECT id, name FROM students WHERE group_id = $1`

	rows, err := s.db.Query(query, groupId)
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

func (s *Storage) DeleteStudent(id int) error {
	query := `DELETE FROM students WHERE id = $1`

	if _, err := s.db.Exec(query, id); err != nil {
		return fmt.Errorf("Error deleting student: %w", err)
	}

	return nil
}

func (s *Storage) UpdateStudent(student *models.Student) error {
	query := `UPDATE students SET name = $1, group_id = $2 WHERE id = $3`

	if _, err := s.db.Exec(query, student.Name, student.GroupID, student.ID); err != nil {
		return fmt.Errorf("Error updating student: %w", err)
	}

	return nil
}
