package services

import (
	"errors"

	"fmt"

	"github.com/ObjoradDdd/FeedbackTeachersHelper/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type TeacherStorage interface {
	CreateTeacher(teacher *models.Teacher, hash string) (int, error)
	GetTeacherByLogin(login string) (*models.Teacher, error)
	DeleteTeacherById(id int) error
}

type TeacherService struct {
	db TeacherStorage
}

func NewTeacherService(db TeacherStorage) *TeacherService {
	return &TeacherService{db: db}
}

func (s *TeacherService) Register(login, password string) (int, error) {

	if login == "" || password == "" {
		return 0, errors.New("login and password cannot be empty")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, fmt.Errorf("error while hashing password: %w", err)
	}

	teacher := &models.Teacher{
		Login: login,
	}

	teacherID, err := s.db.CreateTeacher(teacher, string(hash))
	if err != nil {
		return 0, fmt.Errorf("error in DB while registering: %w", err)
	}

	return teacherID, nil
}

func (s *TeacherService) Login(login, password string) (int, error) {
	if login == "" || password == "" {
		return 0, errors.New("login and password cannot be empty")
	}

	teacher, err := s.db.GetTeacherByLogin(login)
	if err != nil {
		return 0, fmt.Errorf("error in DB while fetching teacher: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(teacher.Hash), []byte(password)); err != nil {
		return 0, errors.New("invalid login or password")
	}

	return teacher.ID, nil
}

func (s *TeacherService) DeleteTeacher(teacherId int) error {
	if err := s.db.DeleteTeacherById(teacherId); err != nil {
		return fmt.Errorf("error in DB while deleting teacher: %w", err)
	}
	return nil
}
