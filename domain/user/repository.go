package user

import (
	"context"
	"errors"
	"github.com/viktor8881/service-utilities/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository struct {
	*db.MongoDB
	collName string
}

func NewRepository(mongoDB *db.MongoDB) *Repository {
	return &Repository{MongoDB: mongoDB, collName: "items"}
}

func (r *Repository) FetchAll(ctx context.Context) ([]Model, error) {
	users := new([]Model)
	return *users, r.MongoDB.FetchAll(ctx, "sql: fetchAll", r.collName, users, map[string]interface{}{})
}

func (r *Repository) FetchAllByEmail(ctx context.Context, email string) ([]Model, error) {
	users := new([]Model)
	args := map[string]interface{}{"value": email}
	return *users, r.MongoDB.FetchAll(ctx, "sql- fetchAllByEmail", r.collName, users, args)
}

func (r *Repository) Get(ctx context.Context, ID string) (*Model, error) {
	user := new(Model)
	objectID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objectID}
	err = r.MongoDB.Get(ctx, "sql- get", r.collName, user, filter)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrNotFound
		}
	}
	return user, err
}

func (r *Repository) Create(ctx context.Context, user Model) (string, error) {
	return r.MongoDB.Create(ctx, "sql- create", r.collName, user)
}

func (r *Repository) Update(ctx context.Context, user Model) (int64, error) {
	objectID, err := primitive.ObjectIDFromHex(user.ID)
	if err != nil {
		return 0, err
	}

	filter := bson.M{"_id": objectID}

	// Создаем обновляемый документ без поля _id
	update := bson.M{
		"$set": bson.M{
			"name":  user.Name,
			"value": user.Email,
		},
	}

	return r.MongoDB.Update(ctx, "sql- update", r.collName, filter, update)
}

func (r *Repository) Delete(ctx context.Context, ID string) (int64, error) {
	objectID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return 0, err
	}

	filter := bson.M{"_id": objectID}
	return r.MongoDB.Delete(ctx, "sql- delete", r.collName, filter)
}
