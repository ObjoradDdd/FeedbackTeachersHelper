package handlers

import (
	"context"
	"net/http"

	"github.com/ObjoradDdd/FeedbackTeachersHelper/internal/dto"
	"github.com/ObjoradDdd/FeedbackTeachersHelper/internal/models"
	"github.com/ObjoradDdd/FeedbackTeachersHelper/internal/services"
	"github.com/go-playground/validator/v10"
)

type feedbackService interface {
	GenerateFeedback(ctx context.Context, input *services.GenerateFeedbackInput, userID int) (*models.GeneratedGroupFeedback, error)
}

type FeedbackHandler struct {
	feedbackService feedbackService
	validator       *validator.Validate
}

func NewFeedbackHandler(feedbackService feedbackService, validator *validator.Validate) *FeedbackHandler {
	return &FeedbackHandler{
		feedbackService: feedbackService,
		validator:       validator,
	}
}

// GetFeedback godoc
// @Summary Generate feedback
// @Description Generates lesson feedback for students in a group
// @Tags feedback
// @Accept json
// @Produce json
// @Security UserID
// @Param input body dto.GetFeedbackRequest true "Feedback payload"
// @Success 200 {object} dto.GetFeedbackResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /feedback [post]
func (h *FeedbackHandler) GetFeedback(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, err := getUserID(w, r)
	if err != nil {
		return
	}

	var req dto.GetFeedbackRequest
	if err := decodeAndValidateRequest(w, r, &req, h.validator); err != nil {
		return
	}

	feedback, err := h.feedbackService.GenerateFeedback(r.Context(), &services.GenerateFeedbackInput{
		GroupID:           req.GroupId,
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
	}, userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to generate feedback")
		return
	}

	respondWithJSON(w, http.StatusOK, feedback.ToDto())
}
