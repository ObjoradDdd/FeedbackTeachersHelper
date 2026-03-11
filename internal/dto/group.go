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
	Name string `json:"name"`
}

type CreateGroupResponse struct {
	Id int `json:"id"`
}

type UpdateGroupRequest struct {
	Name string `json:"name"`
}

type UpdateGroupResponse struct {
	Message string `json:"message"`
}

type DeleteGroupRequest struct {
	Id int `json:"id"`
}

type DeleteGroupResponse struct {
	Message string `json:"message"`
}
