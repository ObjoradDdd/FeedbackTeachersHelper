package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ObjoradDdd/FeedbackTeachersHelper/internal/dto"
	"github.com/ObjoradDdd/FeedbackTeachersHelper/internal/services"
)

type TeacherHandler struct {
	teacherService *services.TeacherService
}

func NewTeacherHandler(teacherService *services.TeacherService) *TeacherHandler {
	return &TeacherHandler{
		teacherService: teacherService,
	}
}

// Register godoc
// @Summary Register teacher
// @Description Creates a new teacher account
// @Tags teacher
// @Accept json
// @Produce json
// @Param input body dto.RegisterRequest true "Register payload"
// @Success 201 {object} dto.RegisterResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /register [post]
func (h *TeacherHandler) Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req dto.RegisterRequest
	if err := DecodeRequest(w, r, &req); err != nil {
		return
	}

	teacherId, err := h.teacherService.Register(req.Login, req.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dto.ErrorResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(dto.RegisterResponse{
		TeacherId: teacherId,
		Message:   "Teacher registered successfully",
	})
}

// Login godoc
// @Summary Login teacher
// @Description Authenticates teacher and returns JWT token
// @Tags teacher
// @Accept json
// @Produce json
// @Param input body dto.LoginRequest true "Login payload"
// @Success 200 {object} dto.LoginResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Router /login [post]
func (h *TeacherHandler) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req dto.LoginRequest
	if err := DecodeRequest(w, r, &req); err != nil {
		return
	}

	token, err := h.teacherService.Login(req.Login, req.Password)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(dto.ErrorResponse{Error: err.Error()})
		return
	}

	json.NewEncoder(w).Encode(dto.LoginResponse{
		Token: token,
	})
}

// AddAPIKey godoc
// @Summary Add API key
// @Description Saves encrypted external API key for current teacher
// @Tags teacher
// @Accept json
// @Produce json
// @Security Bearer
// @Param input body dto.AddAPIKeyRequest true "API key payload"
// @Success 200 {object} dto.AddApiKeyResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /add_api_key [post]
func (h *TeacherHandler) AddAPIKey(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	teacherId, err := GetTeacherIdFromToken(w, r)
	if err != nil {
		return
	}

	var req dto.AddAPIKeyRequest
	if err := DecodeRequest(w, r, &req); err != nil {
		return
	}

	err = h.teacherService.AddApiKey(teacherId, req.APIKey)
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

// DeleteTeacher godoc
// @Summary Delete teacher
// @Description Deletes current teacher account
// @Tags teacher
// @Produce json
// @Security Bearer
// @Success 200 {object} dto.DeleteTeacherResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /delete_teacher [delete]
func (h *TeacherHandler) DeleteTeacher(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	teacherId, err := GetTeacherIdFromToken(w, r)
	if err != nil {
		return
	}

	err = h.teacherService.DeleteTeacher(teacherId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dto.ErrorResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dto.DeleteTeacherResponse{
		Message: "Teacher deleted successfully",
	})
}
