package repository

import (
	"crud-app/app/models"
	"database/sql"
	"fmt"

)

type AlumniRepository interface {
	FindByID(id int) (*models.Alumni, error)
	GetAlumni(search, sortBy, order string, page, limit int) ([]models.Alumni, int, error)
	Create(a *models.Alumni) error
	Update(id int, a *models.Alumni) error
	Delete(id int) error
}

type alumniPostgres struct{ db *sql.DB }

func NewAlumniRepository(db *sql.DB) AlumniRepository {
	return &alumniPostgres{db}
}

func (r *alumniPostgres) GetAlumni(search, sortBy, order string, page, limit int) ([]models.Alumni, int, error) {
	var alumni []models.Alumni
	var total int

	if page < 1 {
	   page = 1
	}

	if limit < 1 {
	   limit = 10 
	}

	offset := (page - 1) * limit

	allowedSort := map[string]bool{"id": true, "nama": true, "angkatan": true, "jurusan":true, "email": true}

	if !allowedSort[sortBy] {
		sortBy = "id"
	}
	if order != "asc" && order != "desc" {
		order = "asc"
	}

	query := `SELECT id, nama, angkatan, jurusan, email FROM alumni WHERE 1=1`
	args := []interface{}{}

	if search != "" {
		query += " AND (nama ILIKE $1 OR jurusan ILIKE $2 OR CAST(angkatan AS TEXT) ILIKE $3 OR email ILIKE $4)"
		args = append(args, "%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	countQuery := "SELECT COUNT(*) FROM (" + query + ") AS sub"
	if len(args) > 0 {
		if err := r.db.QueryRow(countQuery, args...).Scan(&total); err != nil {
			return nil, 0, err
		}
	} else {
		if err := r.db.QueryRow(countQuery).Scan(&total); err != nil {
			return nil, 0, err
		}
	}

	query = fmt.Sprintf("%s ORDER BY %s %s LIMIT %d OFFSET %d", query, sortBy, order, limit, offset)
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}

	defer rows.Close()

	for rows.Next() {
		var a models.Alumni
		if err := rows.Scan(&a.ID, &a.Nama, &a.Angkatan, &a.Jurusan, &a.Email); err != nil {
				return nil, 0, err
		}
		alumni = append(alumni,a)
	}
		return alumni, total, nil
}

func (r *alumniPostgres) FindByID(id int) (*models.Alumni, error) {
	var a models.Alumni
	err := r.db.QueryRow("SELECT id, nama, email FROM alumni WHERE id=$1", id).
		Scan(&a.ID, &a.Nama, &a.Email)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *alumniPostgres) Create(a *models.Alumni) error {
	_, err := r.db.Exec("INSERT INTO alumni(nim, nama, jurusan, angkatan, tahun_lulus, email, no_telepon) VALUES($1,$2,$3,$4,$5,$6,$7)", a.NIM, a.Nama, a.Jurusan, a.Angkatan, a.Tahun_lulus, a.Email, a.No_telp)
	return err
}

func (r *alumniPostgres) Update(id int, a *models.Alumni) error {
	_, err := r.db.Exec("UPDATE alumni SET nama=$1, jurusan=$2, no_telepon=$3 WHERE id=$4", a.Nama, a.Jurusan, a.No_telp, id)
	return err
}

func (r *alumniPostgres) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM alumni WHERE id=$1", id)
	return err
}
