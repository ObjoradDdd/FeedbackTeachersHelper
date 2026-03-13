package services

import (
	"fmt"

	"github.com/ObjoradDdd/FeedbackTeachersHelper/internal/models"
)

type StudentStorage interface {
	CreateStudent(student *models.Student, userID int, groupID int) (int, error)
	GetGroupStudents(id int, userID int) ([]models.Student, error)
	UpdateStudent(student *models.Student, userID int, groupID int) error
	DeleteStudent(id int, userID int) error
}

type StudentService struct {
	db StudentStorage
}

func NewStudentService(db StudentStorage) *StudentService {
	return &StudentService{db: db}
}

func (s *StudentService) CreateStudent(name string, groupID int, userID int) (int, error) {
	student := &models.Student{
		Name: name,
	}

	studentId, err := s.db.CreateStudent(student, userID, groupID)
	if err != nil {
		return 0, fmt.Errorf("error in DB while creating student: %w", err)
	}

	return studentId, nil
}

func (s *StudentService) GetGroupStudents(groupID int, userID int) ([]models.Student, error) {
	students, err := s.db.GetGroupStudents(groupID, userID)

	if err != nil {
		return nil, fmt.Errorf("error in DB while fetching students: %w", err)
	}

	return students, nil
}

func (s *StudentService) UpdateStudent(id int, name string, groupID int, userID int) error {
	student := &models.Student{
		Id:   id,
		Name: name,
	}

	if err := s.db.UpdateStudent(student, userID, groupID); err != nil {
		return fmt.Errorf("error in DB while updating student: %w", err)
	}

	return nil
}

func (s *StudentService) DeleteStudent(id int, userID int) error {
	if err := s.db.DeleteStudent(id, userID); err != nil {
		return fmt.Errorf("error in DB while deleting student: %w", err)
	}
	return nil
}
