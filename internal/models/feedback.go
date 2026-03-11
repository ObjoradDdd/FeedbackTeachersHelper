package models

import "github.com/ObjoradDdd/FeedbackTeachersHelper/internal/dto"

type GeneratedStudentFeedback struct {
	StudentId int
	Name      string
	Feedback  string
}

func (s GeneratedStudentFeedback) ToDto() dto.StudentFeedbackResponse {
	return dto.StudentFeedbackResponse{
		StudentId: s.StudentId,
		Name:      s.Name,
		Feedback:  s.Feedback,
	}
}

type GeneratedGroupFeedback struct {
	TeacherId         int
	GroupId           int
	LessonDescription string
	Students          []GeneratedStudentFeedback
}

func (g GeneratedGroupFeedback) ToDto() dto.GetFeedbackResponse {
	studentsDto := make([]dto.StudentFeedbackResponse, len(g.Students))
	for i, student := range g.Students {
		studentsDto[i] = student.ToDto()
	}

	return dto.GetFeedbackResponse{
		TeacherId:         g.TeacherId,
		GroupId:           g.GroupId,
		LessonDescription: g.LessonDescription,
		Students:          studentsDto,
	}
}
