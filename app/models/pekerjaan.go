package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Pekerjaan struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Alumni_ID       primitive.ObjectID `bson:"alumni_id" json:"alumni_id"`
	Nama_Perusahaan string             `bson:"nama_perusahaan" json:"nama_perusahaan"`
	Posisi_jabatan  string             `bson:"posisi_jabatan" json:"posisi_jabatan"`
	Bidang_industri string             `bson:"bidang_industri" json:"bidang_industri"`
	Lokasi_kerja    string             `bson:"lokasi_kerja" json:"lokasi_kerja"`
	Gaji_range      string             `bson:"gaji_range" json:"gaji_range"`
	Tanggal_kerja   time.Time          `bson:"tanggal_mulai_kerja" json:"tanggal_mulai_kerja"`
	Tanggal_selesai time.Time          `bson:"tanggal_selesai_kerja" json:"tanggal_selesai_kerja"`
	Status          string             `bson:"status_pekerjaan" json:"status_pekerjaan"`
	Deskripsi       string             `bson:"deskripsi_pekerjaan" json:"deskripsi_pekerjaan"`
	CreatedAt       time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt       time.Time          `bson:"updated_at" json:"updated_at"`
	IsDelete        *time.Time         `bson:"is_deleted" json:"is_deleted"`
}
