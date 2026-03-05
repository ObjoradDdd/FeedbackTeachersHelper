package services

import (
	"fmt"

	"github.com/ObjoradDdd/FeedbackTeachersHelper/internal/models"
)

type StudentStorage interface {
	CreateStudent(student *models.Student) (int, error)
	GetGroupStudents(id int) ([]models.Student, error)
	UpdateStudent(student *models.Student) error
	DeleteStudent(id int) error
}

type StudentService struct {
	db StudentStorage
}

func NewStudentService(db StudentStorage) *StudentService {
	return &StudentService{db: db}
}

func (s *StudentService) CreateStudent(name string, groupId int) (int, error) {
	student := &models.Student{
		Name:    name,
		GroupID: groupId,
	}

	studentID, err := s.db.CreateStudent(student)
	if err != nil {
		return 0, fmt.Errorf("error in DB while creating student: %w", err)
	}

	return studentID, nil
}

func (s *StudentService) GetGroupStudents(groupId int) ([]models.Student, error) {
	students, err := s.db.GetGroupStudents(groupId)
	if err != nil {
		return nil, fmt.Errorf("error in DB while fetching students: %w", err)
	}

	return students, nil
}

func (s *StudentService) UpdateStudent(id int, name string, groupId int) error {
	student := &models.Student{
		ID:      id,
		Name:    name,
		GroupID: groupId,
	}

	if err := s.db.UpdateStudent(student); err != nil {
		return fmt.Errorf("error in DB while updating student: %w", err)
	}

	return nil
}

func (s *StudentService) DeleteStudent(id int) error {
	if err := s.db.DeleteStudent(id); err != nil {
		return fmt.Errorf("error in DB while deleting student: %w", err)
	}
	return nil
}
