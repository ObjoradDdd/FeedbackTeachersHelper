package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ObjoradDdd/FeedbackTeachersHelper/internal/dto"
	"github.com/ObjoradDdd/FeedbackTeachersHelper/internal/services"
)

type FeedbackHandler struct {
	feedbackService *services.FeedbackService
}

func NewFeedbackHandler(feedbackService *services.FeedbackService) *FeedbackHandler {
	return &FeedbackHandler{
		feedbackService: feedbackService,
	}
}

func (h *FeedbackHandler) GetFeedback(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	teacherId, err := GetTeacherIdFromToken(w, r)
	if err != nil {
		return
	}

	var req dto.GetFeedbackRequest
	if err := DecodeRequest(w, r, &req); err != nil {
		return
	}

	feedback, err := h.feedbackService.GenerateFeedback(&services.GenerateFeedbackInput{
		TeacherId:         teacherId,
		GroupId:           req.GroupId,
		LessonDescription: req.LessonDescription,
		Activities:        req.Activities,
		Students: func(students []dto.StudentFeedbackRequest) []services.StudentFeedbackInput {
			studentsInput := make([]services.StudentFeedbackInput, len(students))
			for i, student := range students {
				studentsInput[i] = services.StudentFeedbackInput{
					StudentId: student.StudentId,
					Comment:   student.Comment,
					TagIds:    student.TagIds,
				}
			}
			return studentsInput
		}(req.Students),
	}, teacherId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dto.ErrorResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(feedback.ToDto())
}
