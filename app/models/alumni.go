package models

import "time"

type Alumni struct {
	ID          int       `json:"id"`
	NIM         string    `json:"nim"`
	Nama        string    `json:"nama"`
	Jurusan     string    `json:"jurusan"`
	Angkatan    int       `json:"angkatan"`
	Tahun_lulus int       `json:"tahun_lulus"`
	Email       string    `json:"email"`
	No_telp     string    `json:"no_telepon"`
	Alamat      string    `json:"alamat"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type AlumniResponse struct {
	Data []Alumni `json: "data"`
	Meta MetaInfo `json: "meta"`
}
