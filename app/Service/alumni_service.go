package service

import (
	"crud-app/app/models"
	"crud-app/app/repository"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type AlumniService struct {
	repo repository.AlumniRepository
}

func NewAlumniService(r repository.AlumniRepository) *AlumniService {
	return &AlumniService{repo: r}
}

// func (h *AlumniService) GetAll(w http.ResponseWriter, r *http.Request) {
// 	alumni, err := h.repo.FindAll()
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	json.NewEncoder(w).Encode(alumni)
// }

func (h *AlumniService) GetByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	alumni, err := h.repo.FindByID(id)
	if err != nil {
		http.Error(w, "Alumni not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(alumni)
}

func (h *AlumniService) Create(w http.ResponseWriter, r *http.Request) {
	var a models.Alumni
	if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	if err := h.repo.Create(&a); err != nil {
		http.Error(w, "Failed to create alumni", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(a)
}

func (h *AlumniService) Update(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var a models.Alumni
	if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	if err := h.repo.Update(id, &a); err != nil {
		http.Error(w, "Failed to update alumni", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "Updated successfully"})
}

func (h *AlumniService) Delete(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if err := h.repo.Delete(id); err != nil {
		http.Error(w, "Failed to delete alumni", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "Deleted successfully"})
}

func (h *AlumniService) GetAlumni(w http.ResponseWriter, r *http.Request) {
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

	alumni, total, err := h.repo.GetAlumni(search, sortBy, order, page, limit)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := models.AlumniResponse{
		Data: alumni,
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
