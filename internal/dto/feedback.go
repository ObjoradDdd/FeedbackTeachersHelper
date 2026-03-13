package dto

type GetFeedbackRequest struct {
	GroupId           int                      `json:"group_id"`
	LessonDescription string                   `json:"lesson_description"`
	Activities        string                   `json:"activities"`
	Students          []StudentFeedbackRequest `json:"students"`
}

type StudentFeedbackRequest struct {
	StudentId int    `json:"student_id"`
	Comment   string `json:"comment"`
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
