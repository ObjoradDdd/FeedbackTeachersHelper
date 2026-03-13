package services

import (
	"fmt"

	encryption "github.com/ObjoradDdd/FeedbackTeachersHelper/internal/utils"
)

type UserStorage interface {
	DeleteUserById(id int) error
	GetApiKey(userID int) (string, error)
	AddApiKey(userID int, apiKey string) error
	DeleteApiKey(userID int) error
}

type UserService struct {
	db UserStorage
}

func NewUserService(db UserStorage) *UserService {
	return &UserService{db: db}
}

func (s *UserService) DeleteUser(userID int) error {
	if err := s.db.DeleteUserById(userID); err != nil {
		return fmt.Errorf("error in DB while deleting user: %w", err)
	}
	return nil
}

func (s *UserService) AddApiKey(userID int, apiKey string) error {
	apiKeyHash, err := encryption.Encrypt(apiKey)
	if err != nil {
		return fmt.Errorf("error while encrypting API key: %w", err)
	}

	if err := s.db.AddApiKey(userID, apiKeyHash); err != nil {
		return fmt.Errorf("error in DB while adding API key: %w", err)
	}
	return nil
}
