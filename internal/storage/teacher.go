package storage

import "github.com/ObjoradDdd/FeedbackTeachersHelper/internal/models"

func (s *Storage) CreateTeacher(teacher *models.Teacher, hash string) (int, error) {
	query := `INSERT INTO teachers (login, hash) VALUES ($1, $2) RETURNING id`
	var teacherID int

	if err := s.db.QueryRow(query, teacher.Login, hash).Scan(&teacherID); err != nil {
		return 0, err
	}

	return teacherID, nil
}

func (s *Storage) GetTeacherHash(teacherId int) (string, error) {
	query := `SELECT hash FROM teachers WHERE id = $1`
	var hash string

	if err := s.db.QueryRow(query, teacherId).Scan(&hash); err != nil {
		return "", err
	}
	return hash, nil
}

func (s *Storage) GetTeacherByLogin(login string) (*models.Teacher, error) {
	query := `SELECT id, login, hash FROM teachers WHERE login = $1`
	var teacher models.Teacher

	if err := s.db.QueryRow(query, login).Scan(&teacher.ID, &teacher.Login, &teacher.Hash); err != nil {
		return nil, err
	}
	return &teacher, nil
}
