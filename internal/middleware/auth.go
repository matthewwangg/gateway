package middleware

import (
	"context"
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type ContextKey string

const key ContextKey = "claims"

func CheckJWTToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authorizationHeader, "Bearer ") {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authorizationHeader, "Bearer ")
		claims := &jwt.RegisteredClaims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, GetJWTSecret)
		if err != nil || !token.Valid {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), key, claims)))
	})
}

func GetJWTSecret(token *jwt.Token) (interface{}, error) {
	if token.Method != jwt.SigningMethodHS256 {
		return nil, errors.New("invalid token signing method")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, errors.New("no jwt secret set")
	}

	return []byte(jwtSecret), nil
}

func GenerateJWTToken(userID string) (string, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return "", errors.New("environment variable not set")
	}

	claims := jwt.RegisteredClaims{
		Subject:   userID,
		Issuer:    "gateway",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
