package dto

type GetFeedbackRequest struct {
	GroupId           int                      `json:"group_id" validate:"required,gt=0"`
	LessonDescription string                   `json:"lesson_description" validate:"required,min=2,max=512"`
	Activities        string                   `json:"activities" validate:"required,min=2,max=512"`
	Students          []StudentFeedbackRequest `json:"students"`
}

type StudentFeedbackRequest struct {
	StudentId int    `json:"student_id" validate:"required,gt=0"`
	Comment   string `json:"comment" validate:"max=512"`
	TagIds    []int  `json:"tag_ids"`
}

type GetFeedbackResponse struct {
	UserID            int                       `json:"user_id"`
	GroupId           int                       `json:"group_id"`
	LessonDescription string                    `json:"lesson_description"`
	Students          []StudentFeedbackResponse `json:"students"`
}

type StudentFeedbackResponse struct {
	StudentId int    `json:"student_id"`
	Name      string `json:"name"`
	Feedback  string `json:"feedback"`
}
