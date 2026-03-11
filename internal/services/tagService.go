package services

import (
	"fmt"

	"github.com/ObjoradDdd/FeedbackTeachersHelper/internal/models"
)

type TagStorage interface {
	CreateTag(tag *models.Tag, teacherId int) (int, error)
	DeleteTag(id int, teacherId int) error
	GetTeachersTags(teacherId int) ([]models.Tag, error)
	UpdateTag(tag *models.Tag, teacherId int) error
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

func (s *TagService) CreateTag(input CreateTagInput, teacherId int) (int, error) {
	tag := &models.Tag{
		Name:    input.Name,
		Meaning: input.Meaning,
	}

	tagId, err := s.db.CreateTag(tag, teacherId)
	if err != nil {
		return 0, fmt.Errorf("error in DB while registering: %w", err)
	}

	return tagId, nil
}

func (s *TagService) GetTeachersTags(teacherId int) ([]models.Tag, error) {
	tags, err := s.db.GetTeachersTags(teacherId)
	if err != nil {
		return nil, fmt.Errorf("error in DB while fetching tags: %w", err)
	}

	return tags, nil
}

func (s *TagService) DeleteTag(id int, teacherId int) error {
	if err := s.db.DeleteTag(id, teacherId); err != nil {
		return fmt.Errorf("error in DB while deleting tag: %w", err)
	}
	return nil
}

func (s *TagService) UpdateTag(input UpdateTagInput, teacherId int) error {
	tag := &models.Tag{
		Id:      input.Id,
		Name:    input.Name,
		Meaning: input.Meaning,
	}

	if err := s.db.UpdateTag(tag, teacherId); err != nil {
		return fmt.Errorf("error in DB while updating tag: %w", err)
	}
	return nil
}
