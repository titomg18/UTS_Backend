package middleware

import (
	"context"
	// "crud-app/app/models"
	"crud-app/app/repository"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(userRepo repository.UserRepository, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		tokenStr := strings.TrimPrefix(auth, "Bearer ")

		parsed, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil || !parsed.Valid {
			http.Error(w, "Token tidak sesuai", http.StatusUnauthorized)
			return
		}

		claims := parsed.Claims.(jwt.MapClaims)
		id := int(claims["sub"].(float64))

		user, err := userRepo.GetByID(id)
		if err != nil {
			http.Error(w, "User tidak ditemukan", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "user", *user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

