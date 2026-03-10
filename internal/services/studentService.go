package services

import (
	"fmt"

	"github.com/ObjoradDdd/FeedbackTeachersHelper/internal/models"
)

type StudentStorage interface {
	CreateStudent(student *models.Student, teacherId int, groupId int) (int, error)
	GetGroupStudents(id int, teacherId int) ([]models.Student, error)
	UpdateStudent(student *models.Student, teacherId int, groupId int) error
	DeleteStudent(id int, teacherId int) error
}

type StudentService struct {
	db StudentStorage
}

func NewStudentService(db StudentStorage) *StudentService {
	return &StudentService{db: db}
}

func (s *StudentService) CreateStudent(name string, groupId int, teacherId int) (int, error) {
	student := &models.Student{
		Name: name,
	}

	studentId, err := s.db.CreateStudent(student, teacherId, groupId)
	if err != nil {
		return 0, fmt.Errorf("error in DB while creating student: %w", err)
	}

	return studentId, nil
}

func (s *StudentService) GetGroupStudents(groupId int, teacherId int) ([]models.Student, error) {
	students, err := s.db.GetGroupStudents(groupId, teacherId)

	if err != nil {
		return nil, fmt.Errorf("error in DB while fetching students: %w", err)
	}

	return students, nil
}

func (s *StudentService) UpdateStudent(id int, name string, groupId int, teacherId int) error {
	student := &models.Student{
		Id:   id,
		Name: name,
	}

	if err := s.db.UpdateStudent(student, teacherId, groupId); err != nil {
		return fmt.Errorf("error in DB while updating student: %w", err)
	}

	return nil
}

func (s *StudentService) DeleteStudent(id int, teacherId int) error {
	if err := s.db.DeleteStudent(id, teacherId); err != nil {
		return fmt.Errorf("error in DB while deleting student: %w", err)
	}
	return nil
}
