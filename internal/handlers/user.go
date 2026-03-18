package handlers

import (
	"context"
	"net/http"

	"github.com/ObjoradDdd/FeedbackTeachersHelper/internal/dto"
	"github.com/go-playground/validator/v10"
)

type userService interface {
	AddApiKey(ctx context.Context, userID int, apiKey string) error
	DeleteUser(ctx context.Context, userID int) error
}

type UserHandler struct {
	userService userService
	validator   *validator.Validate
}

func NewUserHandler(userService userService, validator *validator.Validate) *UserHandler {
	return &UserHandler{
		userService: userService,
		validator:   validator,
	}
}

// AddAPIKey godoc
// @Summary Add API key
// @Description Saves encrypted external API key for current user
// @Tags users
// @Accept json
// @Produce json
// @Security UserID
// @Param input body dto.AddAPIKeyRequest true "API key payload"
// @Success 200 {object} dto.AddApiKeyResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /add_api_key [post]
func (h *UserHandler) AddAPIKey(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, err := getUserID(w, r)
	if err != nil {
		return
	}

	var req dto.AddAPIKeyRequest
	if err := decodeAndValidateRequest(w, r, &req, h.validator); err != nil {
		return
	}

	err = h.userService.AddApiKey(r.Context(), userID, req.APIKey)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to add api key")
		return
	}

	respondWithJSON(w, http.StatusOK, dto.AddApiKeyResponse{
		Message: "API key added successfully",
	})
}
