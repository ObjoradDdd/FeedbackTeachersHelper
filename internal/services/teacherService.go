package services

import (
	"errors"

	"fmt"

	"github.com/ObjoradDdd/FeedbackTeachersHelper/internal/models"
	"github.com/ObjoradDdd/FeedbackTeachersHelper/internal/utils"
	encryption "github.com/ObjoradDdd/FeedbackTeachersHelper/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type TeacherStorage interface {
	CreateTeacher(teacher *models.Teacher, hash string) (int, error)
	GetTeacherByLogin(login string) (*models.Teacher, error)
	DeleteTeacherById(id int) error
	GetApiKey(teacherId int) (string, error)
	AddApiKey(teacherId int, apiKey string) error
	DeleteApiKey(teacherId int) error
	GetTeacherHash(teacherId int) (string, error)
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

	teacherId, err := s.db.CreateTeacher(teacher, string(hash))
	if err != nil {
		return 0, fmt.Errorf("error in DB while registering: %w", err)
	}

	return teacherId, nil
}

func (s *TeacherService) Login(login, password string) (string, error) {
	if login == "" || password == "" {
		return "", errors.New("login and password cannot be empty")
	}

	teacher, err := s.db.GetTeacherByLogin(login)
	if err != nil {
		return "", fmt.Errorf("error in DB while fetching teacher: %w", err)
	}

	hash, err := s.db.GetTeacherHash(teacher.Id)
	if err != nil {
		return "", fmt.Errorf("error in DB while fetching teacher hash: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return "", errors.New("invalid login or password")
	}

	token, err := utils.GenerateToken(teacher.Id)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return token, nil
}

func (s *TeacherService) DeleteTeacher(teacherId int) error {
	if err := s.db.DeleteTeacherById(teacherId); err != nil {
		return fmt.Errorf("error in DB while deleting teacher: %w", err)
	}
	return nil
}

func (s *TeacherService) AddApiKey(teacherId int, apiKey string) error {
	apiKeyHash, err := encryption.Encrypt(apiKey)
	if err != nil {
		return fmt.Errorf("error while encrypting API key: %w", err)
	}

	if err := s.db.AddApiKey(teacherId, apiKeyHash); err != nil {
		return fmt.Errorf("error in DB while adding API key: %w", err)
	}
	return nil
}
