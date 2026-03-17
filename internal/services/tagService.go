package services

import (
	"context"
	"fmt"

	"github.com/ObjoradDdd/FeedbackTeachersHelper/internal/models"
)

type TagStorage interface {
	CreateTag(ctx context.Context, tag *models.Tag, userID int) (int, error)
	DeleteTag(ctx context.Context, id int, userID int) error
	GetUserTags(ctx context.Context, userID int) ([]models.Tag, error)
	UpdateTag(ctx context.Context, tag *models.Tag, userID int) error
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

func (s *TagService) CreateTag(ctx context.Context, input CreateTagInput, userID int) (int, error) {
	tag := &models.Tag{
		Name:    input.Name,
		Meaning: input.Meaning,
	}

	tagId, err := s.db.CreateTag(ctx, tag, userID)
	if err != nil {
		return 0, fmt.Errorf("error in DB while registering: %w", err)
	}

	return tagId, nil
}

func (s *TagService) GetUserTags(ctx context.Context, userID int) ([]models.Tag, error) {
	tags, err := s.db.GetUserTags(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("error in DB while fetching tags: %w", err)
	}

	return tags, nil
}

func (s *TagService) DeleteTag(ctx context.Context, id int, userID int) error {
	if err := s.db.DeleteTag(ctx, id, userID); err != nil {
		return fmt.Errorf("error in DB while deleting tag: %w", err)
	}
	return nil
}

func (s *TagService) UpdateTag(ctx context.Context, input UpdateTagInput, userID int) error {
	tag := &models.Tag{
		Id:      input.Id,
		Name:    input.Name,
		Meaning: input.Meaning,
	}

	if err := s.db.UpdateTag(ctx, tag, userID); err != nil {
		return fmt.Errorf("error in DB while updating tag: %w", err)
	}
	return nil
}
