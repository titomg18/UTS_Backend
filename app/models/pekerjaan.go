package models

import "time"

type Pekerjaan struct {
	ID              int       `json:"id"`
	Alumni_ID       int       `json:"alumni_id"`
	Nama_Perusahaan string    `json:"nama_perusahaan"`
	Posisi_jabatan  string    `json:"posisi_jabatan"`
	Bidang_industri string    `json:"bidang_industri"`
	Lokasi_kerja    string    `json:"lokasi_kerja"`
	Gaji_range      string    `json:"gaji_range"`
	Tanggal_kerja   time.Time `json:"tanggal_mulai_kerja"`
	Tanggal_selesai time.Time `json:"tanggal_selesai_kerja"`
	Status          string    `json:"status_pekerjaan"`
	Deskripsi       string    `json:"deskripsi_pekerjaan"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	IsDelete        *time.Time `json:"is_deleted"`
}