package auth

import (
	"bank/internal/models"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	jwtSecret   = "MyBankSecret"
	jwtDuration = 24 * time.Hour
)

func AuthenticatedAdmin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := authenticated(w, r)
		if user == nil {
			return
		}

		if !user.Admin {
			http.Error(w, "Insufficient privileges", http.StatusForbidden)
		}

		ctx := user.ClientToContext(r.Context())
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func Authenticated(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := authenticated(w, r)
		if user == nil {
			return
		}

		ctx := user.ClientToContext(r.Context())
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func authenticated(w http.ResponseWriter, r *http.Request) *models.Client {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Authorization header required", http.StatusUnauthorized)
		return nil
	}

	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) != 2 {
		http.Error(w, "Invalid auth header", http.StatusBadRequest)
		return nil
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})
	if err != nil || !token.Valid {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return nil
	}
	user := models.Client{}
	if err = json.Unmarshal([]byte(claims.Subject), &user); err != nil {
		http.Error(w, "Invalid user info", http.StatusUnauthorized)
		return nil
	}

	return &user
}

func GenerateJWTToken(userJson string) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   userJson,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(jwtDuration)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(jwtSecret))
}
