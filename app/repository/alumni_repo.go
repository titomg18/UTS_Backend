package repository

import (
	"context"
	"crud-app/app/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AlumniRepository interface {
	FindByID(id string) (*models.Alumni, error)
	GetAlumni(search, sortBy, order string, page, limit int) ([]models.Alumni, int, error)
	Create(a *models.Alumni) error
	Update(id string, a *models.Alumni) error
	Delete(id string) error
}

type alumniMongo struct {
	collection *mongo.Collection
}

func NewAlumniRepository(db *mongo.Database) AlumniRepository {
	return &alumniMongo{
		collection: db.Collection("alumni"),
	}
}

func (r *alumniMongo) GetAlumni(search, sortBy, order string, page, limit int) ([]models.Alumni, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var alumni []models.Alumni

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	skip := int64((page - 1) * limit)

	allowedSort := map[string]bool{"_id": true, "nama": true, "angkatan": true, "jurusan": true, "email": true}
	if !allowedSort[sortBy] {
		sortBy = "_id"
	}
	if order != "asc" && order != "desc" {
		order = "asc"
	}

	// Build filter
	filter := bson.M{}
	if search != "" {
		filter = bson.M{
			"$or": []bson.M{
				{"nama": bson.M{"$regex": search, "$options": "i"}},
				{"jurusan": bson.M{"$regex": search, "$options": "i"}},
				{"angkatan": search},
				{"email": bson.M{"$regex": search, "$options": "i"}},
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

	if err = cursor.All(ctx, &alumni); err != nil {
		return nil, 0, err
	}

	return alumni, int(total), nil
}

func (r *alumniMongo) FindByID(id string) (*models.Alumni, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var a models.Alumni
	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&a)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *alumniMongo) Create(a *models.Alumni) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	a.ID = primitive.NewObjectID()
	a.CreatedAt = time.Now()
	a.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, a)
	return err
}

func (r *alumniMongo) Update(id string, a *models.Alumni) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	a.UpdatedAt = time.Now()

	update := bson.M{
		"$set": bson.M{
			"nama":       a.Nama,
			"jurusan":    a.Jurusan,
			"no_telepon": a.No_telp,
			"updated_at": a.UpdatedAt,
		},
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	return err
}

func (r *alumniMongo) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objID})
	return err
}
