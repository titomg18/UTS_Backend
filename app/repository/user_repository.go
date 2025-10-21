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

type UserRepository interface {
	GetByUsername(username string) (*models.User, error)
	GetByID(id string) (*models.User, error)
	Create(user *models.User) error
	GetUser(search, sortBy, order string, page, limit int) ([]models.User, int, error)
	SoftDelete(id string) error
}

type userMongo struct {
	collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) UserRepository {
	return &userMongo{
		collection: db.Collection("users"),
	}
}

func (r *userMongo) GetByUsername(username string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var u models.User
	err := r.collection.FindOne(ctx, bson.M{"username": username}).Decode(&u)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userMongo) GetByID(id string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var u models.User
	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&u)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userMongo) Create(u *models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	u.ID = primitive.NewObjectID()

	_, err := r.collection.InsertOne(ctx, u)
	return err
}

func (r *userMongo) GetUser(search, sortBy, order string, page, limit int) ([]models.User, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var users []models.User

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	skip := int64((page - 1) * limit)

	allowedSort := map[string]bool{"_id": true, "username": true, "email": true}
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
				{"username": bson.M{"$regex": search, "$options": "i"}},
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

	if err = cursor.All(ctx, &users); err != nil {
		return nil, 0, err
	}

	return users, int(total), nil
}

func (r *userMongo) SoftDelete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"is_delete": time.Now(),
		},
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}
