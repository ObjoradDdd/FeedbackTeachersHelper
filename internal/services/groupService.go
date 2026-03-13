package services

import (
	"fmt"

	"github.com/ObjoradDdd/FeedbackTeachersHelper/internal/models"
)

type GroupStorage interface {
	CreateGroup(group *models.Group, userID int) (int, error)
	GetUserGroups(userID int) ([]models.Group, error)
	UpdateGroup(group *models.Group, userID int) error
	DeleteGroup(id int, userID int) error
}

type GroupService struct {
	db GroupStorage
}

func NewGroupService(db GroupStorage) *GroupService {
	return &GroupService{db: db}
}

func (s *GroupService) CreateGroup(name string, userID int) (int, error) {
	group := &models.Group{
		Name: name,
	}

	groupId, err := s.db.CreateGroup(group, userID)
	if err != nil {
		return 0, fmt.Errorf("error in DB while creating group: %w", err)
	}
	return groupId, nil
}

func (s *GroupService) GetUserGroups(userID int) ([]models.Group, error) {
	groups, err := s.db.GetUserGroups(userID)
	if err != nil {
		return nil, fmt.Errorf("error in DB while fetching groups: %w", err)
	}
	return groups, nil
}

func (s *GroupService) UpdateGroup(id int, name string, userID int) error {
	group := &models.Group{
		Id:   id,
		Name: name,
	}

	if err := s.db.UpdateGroup(group, userID); err != nil {
		return fmt.Errorf("error in DB while updating group: %w", err)
	}
	return nil
}

func (s *GroupService) DeleteGroup(id int, userID int) error {
	if err := s.db.DeleteGroup(id, userID); err != nil {
		return fmt.Errorf("error in DB while deleting group: %w", err)
	}
	return nil
}
