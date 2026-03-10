package models

type GeneratedStudentFeedback struct {
	StudentId int
	Name      string
	Feedback  string
}

type GeneratedGroupFeedback struct {
	TeacherId         int
	GroupId           int
	LessonDescription string
	Students          []GeneratedStudentFeedback
}
