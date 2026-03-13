package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ObjoradDdd/FeedbackTeachersHelper/internal/dto"
)

type contextKey string

const UserIDKey contextKey = "userId"

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userIDHeader := r.Header.Get("X-User-ID")
		if userIDHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(dto.ErrorResponse{Error: "X-User-ID header missing"})
			return
		}

		userID, err := strconv.Atoi(userIDHeader)
		if err != nil || userID <= 0 {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(dto.ErrorResponse{Error: "Invalid X-User-ID header format"})
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, userID)

		reqWithContext := r.WithContext(ctx)

		next(w, reqWithContext)
	}
}

func EnableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "*")

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-User-ID, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
