package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ObjoradDdd/FeedbackTeachersHelper/internal/dto"
	"github.com/go-playground/validator/v10"
)

func getUserID(w http.ResponseWriter, r *http.Request) (int, error) {
	userID, ok := r.Context().Value(UserIDKey).(int)
	if !ok {
		respondWithError(w, http.StatusInternalServerError, "internal server error: failed to get user id from context")
		return 0, errors.New("user id not found in context")
	}
	return userID, nil
}

func decodeRequest(w http.ResponseWriter, r *http.Request, req any) error {
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		respondWithError(w, http.StatusBadRequest, "bad request: invalid json")
		return err
	}
	return nil
}

func respondWithJSON(w http.ResponseWriter, code int, payload any) {
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(payload); err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to encode response")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(buf.Bytes())
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(dto.ErrorResponse{Error: message})
}

func decodeAndValidateRequest(w http.ResponseWriter, r *http.Request, req any, validator *validator.Validate) error {
	if err := decodeRequest(w, r, req); err != nil {
		respondWithError(w, http.StatusBadRequest, "error decoding request body")
		return err
	}

	if err := validator.Struct(req); err != nil {
		respondWithError(w, http.StatusBadRequest, "validation error")
		return err
	}

	return nil
}
