package dto

type StudentDto struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type CreateStudentRequest struct {
	Name    string `json:"name"`
	GroupId int    `json:"group_id"`
}

type CreateStudentResponse struct {
	Id int `json:"id"`
}

type GetStudentsGroupRequest struct {
	Id int `json:"group_id"`
}

type GetStudentsGroupResponse struct {
	Students []StudentDto `json:"students"`
}

type UpdateStudentRequest struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	GroupId int    `json:"group_id"`
}

type UpdateStudentResponse struct {
	Message string `json:"message"`
}

type DeleteStudentRequest struct {
	Id int `json:"id"`
}

type DeleteStudentResponse struct {
	Message string `json:"message"`
}
