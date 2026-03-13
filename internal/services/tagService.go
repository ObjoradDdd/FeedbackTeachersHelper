package services

import (
	"fmt"

	"github.com/ObjoradDdd/FeedbackTeachersHelper/internal/models"
)

type TagStorage interface {
	CreateTag(tag *models.Tag, userID int) (int, error)
	DeleteTag(id int, userID int) error
	GetUserTags(userID int) ([]models.Tag, error)
	UpdateTag(tag *models.Tag, userID int) error
}

type TagService struct {
	db TagStorage
}

func NewTagService(db TagStorage) *TagService {
	return &TagService{db: db}
}

type CreateTagInput struct {
	Name    string
	Meaning string
}

type UpdateTagInput struct {
	Id      int
	Name    string
	Meaning string
}

func (s *TagService) CreateTag(input CreateTagInput, userID int) (int, error) {
	tag := &models.Tag{
		Name:    input.Name,
		Meaning: input.Meaning,
	}

	tagId, err := s.db.CreateTag(tag, userID)
	if err != nil {
		return 0, fmt.Errorf("error in DB while registering: %w", err)
	}

	return tagId, nil
}

func (s *TagService) GetUserTags(userID int) ([]models.Tag, error) {
	tags, err := s.db.GetUserTags(userID)
	if err != nil {
		return nil, fmt.Errorf("error in DB while fetching tags: %w", err)
	}

	return tags, nil
}

func (s *TagService) DeleteTag(id int, userID int) error {
	if err := s.db.DeleteTag(id, userID); err != nil {
		return fmt.Errorf("error in DB while deleting tag: %w", err)
	}
	return nil
}

func (s *TagService) UpdateTag(input UpdateTagInput, userID int) error {
	tag := &models.Tag{
		Id:      input.Id,
		Name:    input.Name,
		Meaning: input.Meaning,
	}

	if err := s.db.UpdateTag(tag, userID); err != nil {
		return fmt.Errorf("error in DB while updating tag: %w", err)
	}
	return nil
}
