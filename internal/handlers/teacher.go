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
