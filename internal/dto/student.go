package dto

type StudentDto struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type CreateStudentRequest struct {
	Name    string `json:"name" validate:"required,min=2,max=100"`
	GroupId int    `json:"group_id" validate:"required,gt=0"`
}

type CreateStudentResponse struct {
	Id int `json:"id"`
}

type GetStudentsGroupRequest struct {
	Id int `json:"group_id" validate:"required,gt=0"`
}

type GetStudentsGroupResponse struct {
	Students []StudentDto `json:"students"`
}

type UpdateStudentRequest struct {
	Name    string `json:"name" validate:"required,min=2,max=100"`
	GroupId int    `json:"group_id" validate:"required,gt=0"`
}

type UpdateStudentResponse struct {
	Message string `json:"message"`
}

type DeleteStudentRequest struct {
	Id int `json:"id" validate:"required,gt=0"`
}

type DeleteStudentResponse struct {
	Message string `json:"message"`
}
