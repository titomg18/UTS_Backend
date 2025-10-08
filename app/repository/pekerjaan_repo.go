package repository

import (
	"crud-app/app/models"
	"database/sql"
	"fmt"
)

type PekerjaanRepository interface {
	FindByAlumni(alumniID int) ([]models.Pekerjaan, error)
	Create(p *models.Pekerjaan) error
	Update(id int, p *models.Pekerjaan) error
	Delete(id int) error
	GetPekerjaan(search, sortBy, order string, page, limit int) ([]models.Pekerjaan, int, error)
	SoftDeleteByAdmin(alumni_ID int) error
	SoftDeleteByUser(Id int, alumni_id int) error
	FindByPekerjaanID(id int) (*models.Pekerjaan, error)
	
	// New methods for trash management
	GetTrash(search, sortBy, order string, page, limit int) ([]models.Pekerjaan, int, error)
	Restore(pekerjaanID, alumniID int) error
	RestoreByAdmin(alumniID int) error
	HardDelete(pekerjaanID, alumniID int) error
	HardDeleteByAdmin(alumniID int) error
}

type pekerjaanPostgres struct{ db *sql.DB }

func NewPekerjaanRepository(db *sql.DB) PekerjaanRepository {
	return &pekerjaanPostgres{db}
}

func (r *pekerjaanPostgres) FindByAlumni(alumniID int) ([]models.Pekerjaan, error) {
	rows, err := r.db.Query("SELECT id, alumni_id, nama_perusahaan, posisi_jabatan FROM pekerjaan_alumni WHERE alumni_id=$1 AND is_delete IS NULL", alumniID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Pekerjaan
	for rows.Next() {
		var p models.Pekerjaan
		rows.Scan(&p.ID, &p.Alumni_ID, &p.Nama_Perusahaan, &p.Posisi_jabatan)
		list = append(list, p)
	}
	return list, nil
}

func (r *pekerjaanPostgres) FindByPekerjaanID(id int) (*models.Pekerjaan, error) {
	var p models.Pekerjaan
	err := r.db.QueryRow("SELECT id, nama_perusahaan, posisi_jabatan FROM pekerjaan_alumni WHERE id=$1", id).
		Scan(&p.ID, &p.Nama_Perusahaan, &p.Posisi_jabatan)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *pekerjaanPostgres) Create(p *models.Pekerjaan) error {
	_, err := r.db.Exec(
		"INSERT INTO pekerjaan_alumni(alumni_id,nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja, gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan, created_at, updated_at) VALUES($1,$2,$3, $4, $5, $6, $7, $8, $9, $10, $11)",
		p.Alumni_ID, p.Nama_Perusahaan, p.Posisi_jabatan, p.Bidang_industri, p.Lokasi_kerja, p.Gaji_range, p.Tanggal_kerja, p.Tanggal_selesai, p.Status, p.CreatedAt, p.UpdatedAt,
	)
	return err
}

func (r *pekerjaanPostgres) Update(id int, p *models.Pekerjaan) error {
	_, err := r.db.Exec("UPDATE pekerjaan_alumni SET nama_perusahaan=$1, posisi_jabatan=$2 , gaji_range=$3 WHERE id=$4", p.Nama_Perusahaan, p.Posisi_jabatan, p.Gaji_range, id)
	return err
}

func (r *pekerjaanPostgres) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM pekerjaan_alumni WHERE id=$1", id)
	return err
}

func (r *pekerjaanPostgres) GetPekerjaan(search, sortBy, order string, page, limit int) ([]models.Pekerjaan, int, error) {
	var Pekerjaan []models.Pekerjaan
	var total int

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	allowedSort := map[string]bool{
		"id": true, "alumni_id": true, "nama_perusahaan": true,
		"posisi_jabatan": true, "bidang_industri": true,
		"lokasi_kerja": true, "gaji_range": true,
	}
	if !allowedSort[sortBy] {
		sortBy = "id"
	}
	if order != "asc" && order != "desc" {
		order = "asc"
	}

	base := "FROM pekerjaan_alumni WHERE 1=1"
	args := []interface{}{}
	idx := 1

	if search != "" {
		base += fmt.Sprintf(` AND (
			CAST(alumni_id AS TEXT) ILIKE $%d OR
			nama_perusahaan ILIKE $%d OR 
			posisi_jabatan ILIKE $%d OR 
			bidang_industri ILIKE $%d OR 
			lokasi_kerja ILIKE $%d OR 
			gaji_range ILIKE $%d
		)`, idx, idx+1, idx+2, idx+3, idx+4, idx+5)

		args = append(args,
			"%"+search+"%",
			"%"+search+"%",
			"%"+search+"%",
			"%"+search+"%",
			"%"+search+"%",
			"%"+search+"%",
		)
		idx += 6
	}

	// hitung total data
	countQuery := "SELECT COUNT(*) " + base
	if err := r.db.QueryRow(countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// ambil data dengan pagination & sorting
	dataQuery := fmt.Sprintf(`
		SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja, gaji_range 
		%s ORDER BY %s %s LIMIT $%d OFFSET $%d`,
		base, sortBy, order, idx, idx+1)

	args = append(args, limit, offset)

	rows, err := r.db.Query(dataQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var p models.Pekerjaan
		if err := rows.Scan(
			&p.ID,
			&p.Alumni_ID,
			&p.Nama_Perusahaan,
			&p.Posisi_jabatan,
			&p.Bidang_industri,
			&p.Lokasi_kerja,
			&p.Gaji_range,
		); err != nil {
			return nil, 0, err
		}
		Pekerjaan = append(Pekerjaan, p)
	}

	return Pekerjaan, total, nil
}

// Soft delete by Admin
func (r *pekerjaanPostgres) SoftDeleteByAdmin(alumniID int) error {
	_, err := r.db.Exec(`
		UPDATE pekerjaan_alumni 
		SET is_delete = NOW() 
		WHERE alumni_id = $1 AND is_delete IS NULL
	`, alumniID)
	return err
}

// Soft delete by User
func (r *pekerjaanPostgres) SoftDeleteByUser(pekerjaanID, userID int) error {
	_, err := r.db.Exec(`
		UPDATE pekerjaan_alumni 
		SET is_delete = NOW() 
		WHERE id = $1 AND alumni_id = $2 AND is_delete IS NULL
	`, pekerjaanID, userID)
	return err
}

// Method untuk trash management
func (r *pekerjaanPostgres) GetTrash(search, sortBy, order string, page, limit int) ([]models.Pekerjaan, int, error) {
	var pekerjaan []models.Pekerjaan
	var total int

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	// Validasi sorting
	allowedSort := map[string]bool{
		"id": true, "alumni_id": true, "nama_perusahaan": true,
		"posisi_jabatan": true, "is_deleted": true,
	}
	if !allowedSort[sortBy] {
		sortBy = "is_deleted"
	}
	if order != "asc" && order != "desc" {
		order = "desc"
	}

	// Base query untuk data yang sudah di-soft delete
	base := "FROM pekerjaan_alumni WHERE is_delete IS NOT NULL"
	args := []interface{}{}
	idx := 1

	if search != "" {
		base += fmt.Sprintf(` AND (
			CAST(alumni_id AS TEXT) ILIKE $%d OR
			nama_perusahaan ILIKE $%d OR 
			posisi_jabatan ILIKE $%d
		)`, idx, idx+1, idx+2)

		args = append(args,
			"%"+search+"%",
			"%"+search+"%",
			"%"+search+"%",
		)
		idx += 3
	}

	// Hitung total data trash
	countQuery := "SELECT COUNT(*) " + base
	if err := r.db.QueryRow(countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Ambil data trash dengan pagination
	dataQuery := fmt.Sprintf(`
		SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, 
		       lokasi_kerja, gaji_range, is_delete
		%s ORDER BY %s %s LIMIT $%d OFFSET $%d`,
		base, sortBy, order, idx, idx+1)

	args = append(args, limit, offset)

	rows, err := r.db.Query(dataQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var p models.Pekerjaan
		if err := rows.Scan(
			&p.ID,
			&p.Alumni_ID,
			&p.Nama_Perusahaan,
			&p.Posisi_jabatan,
			&p.Bidang_industri,
			&p.Lokasi_kerja,
			&p.Gaji_range,
			&p.IsDelete,
		); err != nil {
			return nil, 0, err
		}
		pekerjaan = append(pekerjaan, p)
	}

	return pekerjaan, total, nil
}

// Restore soft deleted data
func (r *pekerjaanPostgres) Restore(pekerjaanID, alumniID int) error {
	query := `
		UPDATE pekerjaan_alumni 
		SET is_delete = NULL 
		WHERE id = $1 AND alumni_id = $2 AND is_delete IS NOT NULL
	`
	result, err := r.db.Exec(query, pekerjaanID, alumniID)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("data not found or already restored")
	}

	return nil
}

// Restore all by admin
func (r *pekerjaanPostgres) RestoreByAdmin(alumniID int) error {
	query := `
		UPDATE pekerjaan_alumni 
		SET is_delete = NULL 
		WHERE alumni_id = $1 AND is_delete IS NOT NULL
	`
	result, err := r.db.Exec(query, alumniID)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("no data found to restore")
	}

	return nil
}

// Hard delete (permanen hapus)
func (r *pekerjaanPostgres) HardDelete(pekerjaanID, alumniID int) error {
	query := `
		DELETE FROM pekerjaan_alumni 
		WHERE id = $1 AND alumni_id = $2 AND is_delete IS NOT NULL
	`
	result, err := r.db.Exec(query, pekerjaanID, alumniID)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("data not found or not in trash")
	}

	return nil
}

// Hard delete all by admin
func (r *pekerjaanPostgres) HardDeleteByAdmin(alumniID int) error {
	query := `
		DELETE FROM pekerjaan_alumni 
		WHERE alumni_id = $1 AND is_delete IS NOT NULL
	`
	result, err := r.db.Exec(query, alumniID)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("no data found in trash")
	}

	return nil
}