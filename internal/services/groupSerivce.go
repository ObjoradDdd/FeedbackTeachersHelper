package services

import (
	"fmt"

	"github.com/ObjoradDdd/FeedbackTeachersHelper/internal/models"
)

type GroupStorage interface {
	CreateGroup(group *models.Group, teacherId int) (int, error)
	GetTeachersGroups(teacherId int) ([]models.Group, error)
	UpdateGroup(group *models.Group, teacherID int) error
	DeleteGroup(id int, teacherID int) error
}

type GroupService struct {
	db GroupStorage
}

func NewGroupService(db GroupStorage) *GroupService {
	return &GroupService{db: db}
}

func (s *GroupService) CreateGroup(name string, teacherId int) (int, error) {
	group := &models.Group{
		Name: name,
	}

	groupID, err := s.db.CreateGroup(group, teacherId)
	if err != nil {
		return 0, fmt.Errorf("error in DB while creating group: %w", err)
	}
	return groupID, nil
}

func (s *GroupService) GetTeachersGroups(teacherId int) ([]models.Group, error) {
	groups, err := s.db.GetTeachersGroups(teacherId)
	if err != nil {
		return nil, fmt.Errorf("error in DB while fetching groups: %w", err)
	}
	return groups, nil
}

func (s *GroupService) UpdateGroup(id int, name string, teacherID int) error {
	group := &models.Group{
		ID:   id,
		Name: name,
	}

	if err := s.db.UpdateGroup(group, teacherID); err != nil {
		return fmt.Errorf("error in DB while updating group: %w", err)
	}
	return nil
}

func (s *GroupService) DeleteGroup(id int, teacherId int) error {
	if err := s.db.DeleteGroup(id, teacherId); err != nil {
		return fmt.Errorf("error in DB while deleting group: %w", err)
	}
	return nil
}
