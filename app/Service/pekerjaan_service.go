package service

import (
	"crud-app/app/models"
	"crud-app/app/repository"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type PekerjaanService struct {
	repo repository.PekerjaanRepository
}

func NewPekerjaanService(r repository.PekerjaanRepository) *PekerjaanService {
	return &PekerjaanService{repo: r}
}

func (h *PekerjaanService) GetByAlumni(w http.ResponseWriter, r *http.Request) {
	alumniID := mux.Vars(r)["alumni_id"]
	data, err := h.repo.FindByAlumni(alumniID)
	if err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(data)
}

func (h *PekerjaanService) Create(w http.ResponseWriter, r *http.Request) {
	var p models.Pekerjaan
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	if err := h.repo.Create(&p); err != nil {
		http.Error(w, "Failed to create pekerjaan", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(p)
}

func (h *PekerjaanService) GetByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	pekerjaan, err := h.repo.FindByPekerjaanID(id)
	if err != nil {
		http.Error(w, "Data not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pekerjaan)
}

func (h *PekerjaanService) Update(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var p models.Pekerjaan
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	if err := h.repo.Update(id, &p); err != nil {
		http.Error(w, "Failed to update pekerjaan", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "Updated successfully"})
}

// func (h *PekerjaanService) Delete(w http.ResponseWriter, r *http.Request) {
// 	id, _ := strconv.Atoi(mux.Vars(r)["id"])
// 	if err := h.repo.Delete(id); err != nil {
// 		http.Error(w, "Failed to delete pekerjaan", http.StatusInternalServerError)
// 		return
// 	}
// 	json.NewEncoder(w).Encode(map[string]string{"message": "Deleted successfully"})
// }

// func (h *PekerjaanService) GetAll(w http.ResponseWriter, r *http.Request) {
// 	pekerjaan, err := h.repo.FindAllPekerjaan()
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	json.NewEncoder(w).Encode(pekerjaan)
// }
 

func (h *PekerjaanService) GetPekerjaan(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	page, _ := strconv.Atoi(q.Get("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(q.Get("limit"))
	if limit < 1 {
		limit = 10
	}
	search := q.Get("search")
	sortBy := q.Get("sortBy")
	if sortBy == "" {
		sortBy = "id"
	}
	order := q.Get("order")
	if order == "" {
		order = "asc"
	}

	data, total, err := h.repo.GetPekerjaan(search, sortBy, order, page, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := map[string]interface{}{
		"data": data,
		"meta": map[string]interface{}{
			"page":   page,
			"limit":  limit,
			"total":  total,
			"pages":  (total + limit - 1) / limit,
			"sortBy": sortBy,
			"order":  order,
			"search": search,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}


func (s *PekerjaanService) SoftDeletePekerjaan(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pekerjaanID := vars["id"]

	userVal := r.Context().Value("user")
	if userVal == nil {
		http.Error(w, "Unauthorized: user not found in context", http.StatusUnauthorized)
		return
	}

	user := userVal.(models.User)

	if user.Role == "admin" {
		alumniID := r.URL.Query().Get("alumni_id")

		if err := s.repo.SoftDeleteByAdmin(alumniID); err != nil {
			http.Error(w, "Failed to soft delete pekerjaan (admin)", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"message": "Semua riwayat pekerjaan alumni berhasil dihapus"})
		return
	}

	userIDStr := user.ID.Hex()
	if err := s.repo.SoftDeleteByUser(pekerjaanID, userIDStr); err != nil {
		http.Error(w, "Failed to soft delete pekerjaan (user)", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Pekerjaan berhasil dihapus"})
}

// GetTrash - Get semua data yang sudah di-soft delete
func (h *PekerjaanService) GetTrash(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	page, _ := strconv.Atoi(q.Get("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(q.Get("limit"))
	if limit < 1 {
		limit = 10
	}
	search := q.Get("search")
	sortBy := q.Get("sortBy")
	if sortBy == "" {
		sortBy = "is_delete"
	}
	order := q.Get("order")
	if order == "" {
		order = "desc"
	}

	data, total, err := h.repo.GetTrash(search, sortBy, order, page, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := map[string]interface{}{
		"data": data,
		"meta": map[string]interface{}{
			"page":   page,
			"limit":  limit,
			"total":  total,
			"pages":  (total + limit - 1) / limit,
			"sortBy": sortBy,
			"order":  order,
			"search": search,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// RestorePekerjaan - Restore data dari trash
func (h *PekerjaanService) RestorePekerjaan(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pekerjaanID := vars["id"]

	userVal := r.Context().Value("user")
	if userVal == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	user := userVal.(models.User)

	var err error

	if user.Role == "admin" {
		alumniIDStr := r.URL.Query().Get("alumni_id")
		if alumniIDStr != "" {
			err = h.repo.RestoreByAdmin(alumniIDStr)
		} else {
			err = h.repo.Restore(pekerjaanID, "")
		}
	} else {
		userIDStr := user.ID.Hex()
		err = h.repo.Restore(pekerjaanID, userIDStr)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Data restored successfully"})
}

// HardDeletePekerjaan - Hapus permanen data dari trash
func (h *PekerjaanService) HardDeletePekerjaan(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pekerjaanID := vars["id"]

	userVal := r.Context().Value("user")
	if userVal == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	user := userVal.(models.User)

	var err error

	if user.Role == "admin" {
		alumniIDStr := r.URL.Query().Get("alumni_id")
		if alumniIDStr != "" {
			err = h.repo.HardDeleteByAdmin(alumniIDStr)
		} else {
			err = h.repo.HardDelete(pekerjaanID, "")
		}
	} else {
		userIDStr := user.ID.Hex()
		err = h.repo.HardDelete(pekerjaanID, userIDStr)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Data permanently deleted"})
}
