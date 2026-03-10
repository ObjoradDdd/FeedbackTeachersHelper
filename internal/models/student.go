package models

import "github.com/ObjoradDdd/FeedbackTeachersHelper/internal/dto"

type Student struct {
	Name string
	Id   int
}

func (s *Student) ToDto() *dto.StudentDto {
	return &dto.StudentDto{
		Id:   s.Id,
		Name: s.Name,
	}
}
