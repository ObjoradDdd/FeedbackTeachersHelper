package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func getJWTSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		panic("fatal error: JWT_SECRET environment variable is not set")
	}
	return []byte(secret)
}

func GenerateToken(teacherID int) (string, error) {
	claims := jwt.MapClaims{
		"teacher_id": teacherID,
		"exp":        time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(getJWTSecret())
}

func ParseToken(tokenString string) (int, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Invalid signing method")
		}
		return getJWTSecret(), nil
	})

	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		teacherIDFloat, ok := claims["teacher_id"].(float64)
		if !ok {
			return 0, errors.New("Invalid teacher_id format in token")
		}
		return int(teacherIDFloat), nil
	}

	return 0, errors.New("Invalid token")
}
