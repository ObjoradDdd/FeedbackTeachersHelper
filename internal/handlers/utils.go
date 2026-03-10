package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ObjoradDdd/FeedbackTeachersHelper/internal/dto"
)

func RespondWithError(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(dto.ErrorResponse{Error: message})
}

func GetTeacherIdFromToken(w http.ResponseWriter, r *http.Request) (int, error) {
	teacherId, ok := r.Context().Value(TeacherIdKey).(int)
	if !ok {
		RespondWithError(w, http.StatusInternalServerError, "internal server error: failed to get teacher id from context")
		return 0, errors.New("teacher id not found in context")
	}
	return teacherId, nil
}

func DecodeRequest(w http.ResponseWriter, r *http.Request, req any) error {
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		RespondWithError(w, http.StatusBadRequest, "bad request: invalid json")
		return err
	}
	return nil
}
