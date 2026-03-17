package services

import (
	"context"
	"fmt"
)

type UserStorage interface {
	DeleteUserById(ctx context.Context, id int) error
	GetApiKey(ctx context.Context, userID int) (string, error)
	AddApiKey(ctx context.Context, userID int, apiKey string) error
	DeleteApiKey(ctx context.Context, userID int) error
}

type UserService struct {
	db        UserStorage
	masterKey string
}

func NewUserService(db UserStorage, masterKey string) *UserService {
	return &UserService{db: db, masterKey: masterKey}
}

func (s *UserService) DeleteUser(ctx context.Context, userID int) error {
	if err := s.db.DeleteUserById(ctx, userID); err != nil {
		return fmt.Errorf("error in DB while deleting user: %w", err)
	}
	return nil
}

func (s *UserService) AddApiKey(ctx context.Context, userID int, apiKey string) error {
	apiKeyHash, err := Encrypt(apiKey, s.masterKey)
	if err != nil {
		return fmt.Errorf("error while encrypting API key: %w", err)
	}

	if err := s.db.AddApiKey(ctx, userID, apiKeyHash); err != nil {
		return fmt.Errorf("error in DB while adding API key: %w", err)
	}
	return nil
}

func (s *UserService) DeleteApiKey(ctx context.Context, userID int) error {
	if err := s.db.DeleteApiKey(ctx, userID); err != nil {
		return fmt.Errorf("error in DB while deleting API key: %w", err)
	}
	return nil
}

func (s *UserService) GetApiKey(ctx context.Context, userID int) (string, error) {
	apiKeyHash, err := s.db.GetApiKey(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("error in DB while fetching API key: %w", err)
	}

	if apiKeyHash == "" {
		return "", nil
	}

	apiKey, err := Decrypt(apiKeyHash, s.masterKey)
	if err != nil {
		return "", fmt.Errorf("error while decrypting API key: %w", err)
	}

	return apiKey, nil
}
