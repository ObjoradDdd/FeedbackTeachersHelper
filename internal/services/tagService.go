package services

import (
	"fmt"

	"github.com/ObjoradDdd/FeedbackTeachersHelper/internal/models"
)

type TagStorage interface {
	CreateTag(tag *models.Tag) (int, error)
	DeleteTag(id int) error
	GetTeachersTags(teacherId int) ([]models.Tag, error)
	UpdateTag(tag *models.Tag) error
}

type TagService struct {
	db TagStorage
}

func NewTagService(db TagStorage) *TagService {
	return &TagService{db: db}
}

type CreateTagInput struct {
	Name      string
	Meaning   string
	IsBad     bool
	TeacherID int
}

func (s *TagService) CreateTag(input CreateTagInput) (int, error) {
	tag := &models.Tag{
		Name:      input.Name,
		Meaning:   input.Meaning,
		TeacherID: input.TeacherID,
	}

	tagID, err := s.db.CreateTag(tag)
	if err != nil {
		return 0, fmt.Errorf("error in DB while registering: %w", err)
	}

	return tagID, nil
}

func (s *TagService) GetTeachersTags(teacherId int) ([]models.Tag, error) {
	tags, err := s.db.GetTeachersTags(teacherId)
	if err != nil {
		return nil, fmt.Errorf("error in DB while fetching tags: %w", err)
	}

	return tags, nil
}

func (s *TagService) DeleteTag(id int) error {
	if err := s.db.DeleteTag(id); err != nil {
		return fmt.Errorf("error in DB while deleting tag: %w", err)
	}
	return nil
}

func (s *TagService) UpdateTag(tag *models.Tag) error {
	if err := s.db.UpdateTag(tag); err != nil {
		return fmt.Errorf("error in DB while updating tag: %w", err)
	}
	return nil
}
