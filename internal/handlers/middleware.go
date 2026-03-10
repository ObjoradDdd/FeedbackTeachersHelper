package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/ObjoradDdd/FeedbackTeachersHelper/internal/dto"
	"github.com/ObjoradDdd/FeedbackTeachersHelper/internal/utils"
)

type contextKey string

const TeacherIdKey contextKey = "teacherId"

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(dto.ErrorResponse{Error: "Authorization header missing"})
			return
		}

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(dto.ErrorResponse{Error: "Invalid Authorization header format"})
			return
		}

		tokenString := headerParts[1]

		teacherId, err := utils.ParseToken(tokenString)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(dto.ErrorResponse{Error: err.Error()})
			return
		}

		ctx := context.WithValue(r.Context(), TeacherIdKey, teacherId)

		reqWithContext := r.WithContext(ctx)

		next(w, reqWithContext)
	}
}
