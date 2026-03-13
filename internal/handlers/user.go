package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ObjoradDdd/FeedbackTeachersHelper/internal/dto"
	"github.com/ObjoradDdd/FeedbackTeachersHelper/internal/services"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
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

	userID, err := GetUserID(w, r)
	if err != nil {
		return
	}

	var req dto.AddAPIKeyRequest
	if err := DecodeRequest(w, r, &req); err != nil {
		return
	}

	err = h.userService.AddApiKey(userID, req.APIKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dto.ErrorResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dto.AddApiKeyResponse{
		Message: "API key added successfully",
	})
}

// DeleteUser godoc
// @Summary Delete user
// @Description Deletes current user account
// @Tags users
// @Produce json
// @Security UserID
// @Success 200 {object} dto.DeleteUserResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /delete_user [delete]
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, err := GetUserID(w, r)
	if err != nil {
		return
	}

	err = h.userService.DeleteUser(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dto.ErrorResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dto.DeleteUserResponse{
		Message: "User deleted successfully",
	})
}
