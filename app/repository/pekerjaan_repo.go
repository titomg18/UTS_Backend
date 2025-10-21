package repository

import (
	"context"
	"crud-app/app/models"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PekerjaanRepository interface {
	FindByAlumni(alumniID string) ([]models.Pekerjaan, error)
	Create(p *models.Pekerjaan) error
	Update(id string, p *models.Pekerjaan) error
	Delete(id string) error
	GetPekerjaan(search, sortBy, order string, page, limit int) ([]models.Pekerjaan, int, error)
	SoftDeleteByAdmin(alumni_ID string) error
	SoftDeleteByUser(Id string, alumni_id string) error
	FindByPekerjaanID(id string) (*models.Pekerjaan, error)
	GetTrash(search, sortBy, order string, page, limit int) ([]models.Pekerjaan, int, error)
	Restore(pekerjaanID, alumniID string) error
	RestoreByAdmin(alumniID string) error
	HardDelete(pekerjaanID, alumniID string) error
	HardDeleteByAdmin(alumniID string) error
}

type pekerjaanMongo struct {
	collection *mongo.Collection
}

func NewPekerjaanRepository(db *mongo.Database) PekerjaanRepository {
	return &pekerjaanMongo{
		collection: db.Collection("pekerjaan_alumni"),
	}
}

func (r *pekerjaanMongo) FindByAlumni(alumniID string) ([]models.Pekerjaan, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(alumniID)
	if err != nil {
		return nil, err
	}

	filter := bson.M{
		"alumni_id":  objID,
		"is_deleted": nil,
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var list []models.Pekerjaan
	if err = cursor.All(ctx, &list); err != nil {
		return nil, err
	}

	return list, nil
}

func (r *pekerjaanMongo) FindByPekerjaanID(id string) (*models.Pekerjaan, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var p models.Pekerjaan
	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&p)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *pekerjaanMongo) Create(p *models.Pekerjaan) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	p.ID = primitive.NewObjectID()
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, p)
	return err
}

func (r *pekerjaanMongo) Update(id string, p *models.Pekerjaan) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	p.UpdatedAt = time.Now()

	update := bson.M{
		"$set": bson.M{
			"nama_perusahaan": p.Nama_Perusahaan,
			"posisi_jabatan":  p.Posisi_jabatan,
			"gaji_range":      p.Gaji_range,
			"updated_at":      p.UpdatedAt,
		},
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	return err
}

func (r *pekerjaanMongo) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objID})
	return err
}

func (r *pekerjaanMongo) GetPekerjaan(search, sortBy, order string, page, limit int) ([]models.Pekerjaan, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var pekerjaan []models.Pekerjaan

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	skip := int64((page - 1) * limit)

	allowedSort := map[string]bool{
		"_id": true, "alumni_id": true, "nama_perusahaan": true,
		"posisi_jabatan": true, "bidang_industri": true,
		"lokasi_kerja": true, "gaji_range": true,
	}
	if !allowedSort[sortBy] {
		sortBy = "_id"
	}
	if order != "asc" && order != "desc" {
		order = "asc"
	}

	// Build filter
	filter := bson.M{"is_deleted": nil}
	if search != "" {
		filter = bson.M{
			"$and": []bson.M{
				{"is_deleted": nil},
				{
					"$or": []bson.M{
						{"alumni_id": bson.M{"$regex": search, "$options": "i"}},
						{"nama_perusahaan": bson.M{"$regex": search, "$options": "i"}},
						{"posisi_jabatan": bson.M{"$regex": search, "$options": "i"}},
						{"bidang_industri": bson.M{"$regex": search, "$options": "i"}},
						{"lokasi_kerja": bson.M{"$regex": search, "$options": "i"}},
						{"gaji_range": bson.M{"$regex": search, "$options": "i"}},
					},
				},
			},
		}
	}

	// Count total
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	// Sort order
	sortOrder := int32(1)
	if order == "desc" {
		sortOrder = -1
	}

	// Query with pagination and sorting
	opts := options.Find().
		SetSkip(skip).
		SetLimit(int64(limit)).
		SetSort(bson.M{sortBy: sortOrder})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &pekerjaan); err != nil {
		return nil, 0, err
	}

	return pekerjaan, int(total), nil
}

func (r *pekerjaanMongo) SoftDeleteByAdmin(alumniID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(alumniID)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"is_deleted": time.Now(),
		},
	}

	_, err = r.collection.UpdateMany(ctx, bson.M{"alumni_id": objID, "is_deleted": nil}, update)
	return err
}

func (r *pekerjaanMongo) SoftDeleteByUser(pekerjaanID, alumniID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pekerjaanObjID, err := primitive.ObjectIDFromHex(pekerjaanID)
	if err != nil {
		return err
	}

	alumniObjID, err := primitive.ObjectIDFromHex(alumniID)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"is_deleted": time.Now(),
		},
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": pekerjaanObjID, "alumni_id": alumniObjID, "is_deleted": nil}, update)
	return err
}

func (r *pekerjaanMongo) GetTrash(search, sortBy, order string, page, limit int) ([]models.Pekerjaan, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var pekerjaan []models.Pekerjaan

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	skip := int64((page - 1) * limit)

	allowedSort := map[string]bool{
		"_id": true, "alumni_id": true, "nama_perusahaan": true,
		"posisi_jabatan": true, "is_deleted": true,
	}
	if !allowedSort[sortBy] {
		sortBy = "is_deleted"
	}
	if order != "asc" && order != "desc" {
		order = "desc"
	}

	// Build filter for deleted items
	filter := bson.M{"is_deleted": bson.M{"$ne": nil}}
	if search != "" {
		filter = bson.M{
			"$and": []bson.M{
				{"is_deleted": bson.M{"$ne": nil}},
				{
					"$or": []bson.M{
						{"alumni_id": bson.M{"$regex": search, "$options": "i"}},
						{"nama_perusahaan": bson.M{"$regex": search, "$options": "i"}},
						{"posisi_jabatan": bson.M{"$regex": search, "$options": "i"}},
					},
				},
			},
		}
	}

	// Count total
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	// Sort order
	sortOrder := int32(1)
	if order == "desc" {
		sortOrder = -1
	}

	// Query with pagination and sorting
	opts := options.Find().
		SetSkip(skip).
		SetLimit(int64(limit)).
		SetSort(bson.M{sortBy: sortOrder})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &pekerjaan); err != nil {
		return nil, 0, err
	}

	return pekerjaan, int(total), nil
}

func (r *pekerjaanMongo) Restore(pekerjaanID, alumniID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pekerjaanObjID, err := primitive.ObjectIDFromHex(pekerjaanID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": pekerjaanObjID, "is_deleted": bson.M{"$ne": nil}}
	if alumniID != "" {
		alumniObjID, err := primitive.ObjectIDFromHex(alumniID)
		if err != nil {
			return err
		}
		filter["alumni_id"] = alumniObjID
	}

	update := bson.M{
		"$set": bson.M{
			"is_deleted": nil,
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("data not found or already restored")
	}

	return nil
}

func (r *pekerjaanMongo) RestoreByAdmin(alumniID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(alumniID)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"is_deleted": nil,
		},
	}

	result, err := r.collection.UpdateMany(ctx, bson.M{"alumni_id": objID, "is_deleted": bson.M{"$ne": nil}}, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("no data found to restore")
	}

	return nil
}

func (r *pekerjaanMongo) HardDelete(pekerjaanID, alumniID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pekerjaanObjID, err := primitive.ObjectIDFromHex(pekerjaanID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": pekerjaanObjID, "is_deleted": bson.M{"$ne": nil}}
	if alumniID != "" {
		alumniObjID, err := primitive.ObjectIDFromHex(alumniID)
		if err != nil {
			return err
		}
		filter["alumni_id"] = alumniObjID
	}

	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("data not found or not in trash")
	}

	return nil
}

func (r *pekerjaanMongo) HardDeleteByAdmin(alumniID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(alumniID)
	if err != nil {
		return err
	}

	result, err := r.collection.DeleteMany(ctx, bson.M{"alumni_id": objID, "is_deleted": bson.M{"$ne": nil}})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("no data found in trash")
	}

	return nil
}
