package dto

type StudentDto struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	GroupID int    `json:"group_id"`
}

type CreateStudentRequest struct {
	Name    string `json:"name"`
	GroupID int    `json:"group_id"`
}

type CreateStudentResponse struct {
	ID int `json:"id"`
}
