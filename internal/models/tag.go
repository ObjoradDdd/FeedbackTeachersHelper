package models

import "github.com/ObjoradDdd/FeedbackTeachersHelper/internal/dto"

type Tag struct {
	Name    string
	Meaning string
	Id      int
}

func (t *Tag) ToDto() dto.TagDto {
	return dto.TagDto{
		Id:      t.Id,
		Meaning: t.Meaning,
		Name:    t.Name,
	}
}
