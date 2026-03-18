package dto

type GroupDto struct {
	Id       int          `json:"id"`
	Name     string       `json:"name"`
	Students []StudentDto `json:"students,omitempty"`
}

type GetGroupsResponse struct {
	Groups []GroupDto `json:"groups"`
}

type CreateGroupRequest struct {
	Name string `json:"name" validate:"required,min=2,max=100"`
}

type CreateGroupResponse struct {
	Id int `json:"id"`
}

type UpdateGroupRequest struct {
	Name string `json:"name" validate:"required,min=2,max=100"`
}

type UpdateGroupResponse struct {
	Message string `json:"message"`
}

type DeleteGroupRequest struct {
	Id int `json:"id" validate:"required,gt=0"`
}

type DeleteGroupResponse struct {
	Message string `json:"message"`
}
