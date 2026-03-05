package models

type GeneratedStudentFeedback struct {
	StudentID int
	Name      string
	Feedback  string
}

type GeneratedGroupFeedback struct {
	TeacherID         int
	GroupID           int
	LessonDescription string
	Students          []GeneratedStudentFeedback
}
