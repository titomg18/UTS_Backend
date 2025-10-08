package service

import (
	"crud-app/app/models"
	"crud-app/app/repository"
	"encoding/json"
	"net/http"
	"strconv"

	// "github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
)

type PekerjaanService struct {
	repo repository.PekerjaanRepository
}

func NewPekerjaanService(r repository.PekerjaanRepository) *PekerjaanService {
	return &PekerjaanService{repo: r}
}

func (h *PekerjaanService) GetByAlumni(w http.ResponseWriter, r *http.Request) {
	// vars := mux.Vars(r)
    // alumniID := vars["alumni_id"]
	alumniID, _ := strconv.Atoi(mux.Vars(r)["alumni_id"])
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
    id, err := strconv.Atoi(mux.Vars(r)["id"])
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    pekerjaan, err := h.repo.FindByPekerjaanID(id)
    if err != nil {
        http.Error(w, "Data not found", http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(pekerjaan)
}

func (h *PekerjaanService) Update(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
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
    pekerjaanID, _ := strconv.Atoi(vars["id"])

    // Ambil user dari context (dimasukkan di middleware)
    userVal := r.Context().Value("user")
    if userVal == nil {
        http.Error(w, "Unauthorized: user not found in context", http.StatusUnauthorized)
        return
    }

    user := userVal.(models.User)

    if user.Role == "admin" {
        // Admin hapus semua pekerjaan berdasarkan alumni_id (query param)
        alumniIDStr := r.URL.Query().Get("alumni_id")
        alumniID, _ := strconv.Atoi(alumniIDStr)

        if err := s.repo.SoftDeleteByAdmin(alumniID); err != nil {
            http.Error(w, "Failed to soft delete pekerjaan (admin)", http.StatusInternalServerError)
            return
        }
        json.NewEncoder(w).Encode(map[string]string{"message": "Semua riwayat pekerjaan alumni berhasil dihapus"})
        return
    }

    // User biasa hanya boleh hapus pekerjaannya sendiri
    if err := s.repo.SoftDeleteByUser(pekerjaanID, user.ID); err != nil {
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
	pekerjaanID, _ := strconv.Atoi(vars["id"])

	// Ambil user dari context
	userVal := r.Context().Value("user")
	if userVal == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	user := userVal.(models.User)

	var err error

	if user.Role == "admin" {
		// Admin bisa restore semua data atau data alumni tertentu
		alumniIDStr := r.URL.Query().Get("alumni_id")
		if alumniIDStr != "" {
			alumniID, _ := strconv.Atoi(alumniIDStr)
			err = h.repo.RestoreByAdmin(alumniID)
		} else {
			// Restore data tertentu
			err = h.repo.Restore(pekerjaanID, 0) // 0 berarti tidak check alumni_id
		}
	} else {
		// User hanya bisa restore data milik sendiri
		err = h.repo.Restore(pekerjaanID, user.ID)
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
	pekerjaanID, _ := strconv.Atoi(vars["id"])

	// Ambil user dari context
	userVal := r.Context().Value("user")
	if userVal == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	user := userVal.(models.User)

	var err error

	if user.Role == "admin" {
		// Admin bisa hard delete semua data atau data alumni tertentu
		alumniIDStr := r.URL.Query().Get("alumni_id")
		if alumniIDStr != "" {
			alumniID, _ := strconv.Atoi(alumniIDStr)
			err = h.repo.HardDeleteByAdmin(alumniID)
		} else {
			// Hard delete data tertentu
			err = h.repo.HardDelete(pekerjaanID, 0) // 0 berarti tidak check alumni_id
		}
	} else {
		// User hanya bisa hard delete data milik sendiri
		err = h.repo.HardDelete(pekerjaanID, user.ID)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Data permanently deleted"})
}
