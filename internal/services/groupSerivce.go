package services

import (
	"fmt"

	"github.com/ObjoradDdd/FeedbackTeachersHelper/internal/models"
)

type GroupStorage interface {
	CreateGroup(group *models.Group) (int, error)
	GetTeacherGroups(teacherId int) ([]models.Group, error)
	UpdateGroup(group *models.Group) error
	DeleteGroup(id int) error
}

type GroupService struct {
	db GroupStorage
}

func NewGroupService(db GroupStorage) *GroupService {
	return &GroupService{db: db}
}

func (s *GroupService) CreateGroup(name string, teacherId int) (int, error) {
	group := &models.Group{
		Name:      name,
		TeacherID: teacherId,
	}

	groupID, err := s.db.CreateGroup(group)
	if err != nil {
		return 0, fmt.Errorf("error in DB while creating group: %w", err)
	}
	return groupID, nil
}

func (s *GroupService) GetTeacherGroups(teacherId int) ([]models.Group, error) {
	groups, err := s.db.GetTeacherGroups(teacherId)
	if err != nil {
		return nil, fmt.Errorf("error in DB while fetching groups: %w", err)
	}
	return groups, nil
}

func (s *GroupService) UpdateGroup(id int, name string) error {
	group := &models.Group{
		ID:   id,
		Name: name,
	}

	if err := s.db.UpdateGroup(group); err != nil {
		return fmt.Errorf("error in DB while updating group: %w", err)
	}
	return nil
}

func (s *GroupService) DeleteGroup(id int) error {
	if err := s.db.DeleteGroup(id); err != nil {
		return fmt.Errorf("error in DB while deleting group: %w", err)
	}
	return nil
}
