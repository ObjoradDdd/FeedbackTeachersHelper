package models

import "github.com/ObjoradDdd/FeedbackTeachersHelper/internal/dto"

type Group struct {
	Name string
	Id   int
}

func (g *Group) ToDto() *dto.GroupDto {
	return &dto.GroupDto{
		Id:   g.Id,
		Name: g.Name,
	}
}
