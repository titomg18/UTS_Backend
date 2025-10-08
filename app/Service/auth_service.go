package service

import (
	"crud-app/app/models"
	"crud-app/app/repository"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo repository.UserRepository
}

func NewAuthService(r repository.UserRepository) *AuthService {
	return &AuthService{repo: r}
}

func (h *AuthService) Register(w http.ResponseWriter, r *http.Request) {
    var u models.User
    if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
        http.Error(w, "Invalid input", http.StatusBadRequest)
        return
    }

    // Validasi minimal
    if u.Email == "" || u.Username == "" || u.Password == "" {
        http.Error(w, "Email, Username, and Password are required", http.StatusBadRequest)
        return
    }

    if len(u.Password) < 6 {
        http.Error(w, "Password must be at least 6 characters", http.StatusBadRequest)
        return
    }

    // Hash password
    hash, _ := bcrypt.GenerateFromPassword([]byte(u.Password), 10)
    u.Password = string(hash)

    // Default role
    if u.Role == "" {
        u.Role = "user"
    }

    // Simpan user ke DB
    if err := h.repo.Create(&u); err != nil {
        http.Error(w, "User already exists or DB error", http.StatusBadRequest)
        return
    }

    json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
}


// Login user
// func (h *AuthService) Login(w http.ResponseWriter, r *http.Request) {
// 	var req models.User
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		http.Error(w, "Invalid input", http.StatusBadRequest)
// 		return
// 	}

// 	user, err := h.repo.GetByUsername(req.Username)
// 	if err != nil {
// 		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
// 		return
// 	}

// 	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) != nil {
// 		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
// 		return
// 	}

// 	claims := jwt.MapClaims{
// 		"sub":  user.ID,
// 		"role": user.Role,
// 		"exp":  time.Now().Add(time.Hour * 24).Unix(),
// 	}

// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 	secret := os.Getenv("JWT_SECRET")
// 	t, _ := token.SignedString([]byte(secret))

// 	json.NewEncoder(w).Encode(map[string]string{"token": t})
// }

// Login user
func (h *AuthService) Login(w http.ResponseWriter, r *http.Request) {
	var req models.User
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	user, err := h.repo.GetByUsername(req.Username)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	claims := jwt.MapClaims{
		"sub":  user.ID,
		"role": user.Role,
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")
	t, _ := token.SignedString([]byte(secret))

	json.NewEncoder(w).Encode(map[string]string{"token": t})
}