package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Alumni struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	NIM         string             `bson:"nim" json:"nim"`
	Nama        string             `bson:"nama" json:"nama"`
	Jurusan     string             `bson:"jurusan" json:"jurusan"`
	Angkatan    int                `bson:"angkatan" json:"angkatan"`
	Tahun_lulus int                `bson:"tahun_lulus" json:"tahun_lulus"`
	Email       string             `bson:"email" json:"email"`
	No_telp     string             `bson:"no_telepon" json:"no_telepon"`
	Alamat      string             `bson:"alamat" json:"alamat"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}

type AlumniResponse struct {
	Data []Alumni `json:"data"`
	Meta MetaInfo `json:"meta"`
}
