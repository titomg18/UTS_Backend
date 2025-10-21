package service

import (
	"encoding/json"
	"net/http"
	"strconv"

	"crud-app/app/models"
	"crud-app/app/repository"

	"github.com/gorilla/mux"
)

type UserService struct {
	Repo repository.UserRepository
}

func NewUserHandler(repo repository.UserRepository) *UserService {
	return &UserService{Repo: repo}
}


func (h *UserService) GetUsers(w http.ResponseWriter, r *http.Request) {
	// Ambil query params
	query := r.URL.Query()
	page, _ := strconv.Atoi(query.Get("page"))
	if page == 0 {
		page = 1
	}
	limit, _ := strconv.Atoi(query.Get("limit"))
	if limit == 0 {
		limit = 10
	}
	search := query.Get("search")
	sortBy := query.Get("sortBy")
	if sortBy == "" {
		sortBy = "id"
	}
	order := query.Get("order")
	if order == "" {
		order = "asc"
	}

	// Ambil data dari repository
	users, total, err := h.Repo.GetUser(search, sortBy, order, page, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := models.UserResponse{
		Data: users,
		Meta: models.MetaInfo{
			Page:   page,
			Limit:  limit,
			Total:  total,
			Pages:  (total + limit - 1) / limit,
			SortBy: sortBy,
			Order:  order,
			Search: search,
		},
	}

	// Set header dan kirim response JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *UserService) SoftDeleteUser(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	if err := h.Repo.SoftDelete(id); err != nil {
		http.Error(w, "Failed to soft delete user", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "User soft deleted"})
}
