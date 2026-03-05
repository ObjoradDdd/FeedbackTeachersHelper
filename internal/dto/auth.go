package dto

type RegisterRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	TeacherID int    `json:"teacher_id"`
	Message   string `json:"message"`
}

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type DeleteTeacherResponse struct {
	Message string `json:"message"`
}
