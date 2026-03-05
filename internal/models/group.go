package models

import "github.com/ObjoradDdd/FeedbackTeachersHelper/internal/dto"

type Group struct {
	Name string
	ID   int
}

func (g *Group) ToDto() *dto.GroupDto {
	return &dto.GroupDto{
		ID:   g.ID,
		Name: g.Name,
	}
}
