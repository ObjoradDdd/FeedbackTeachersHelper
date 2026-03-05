package models

import "github.com/ObjoradDdd/FeedbackTeachersHelper/internal/dto"

type Student struct {
	Name    string
	ID      int
	GroupID int
}

func (s *Student) ToDto() *dto.StudentDto {
	return &dto.StudentDto{
		ID:      s.ID,
		Name:    s.Name,
		GroupID: s.GroupID,
	}
}
