package repository

import (
	"context"

	"github.com/teamlix/user-service/internal/domain"
	"github.com/teamlix/user-service/internal/pkg/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	drv "go.mongodb.org/mongo-driver/mongo"
)

const (
	nameField  = "name"
	emailField = "email"
)

type Repository struct {
	Db *mongo.Mongo
}

func NewRepository(db *mongo.Mongo) *Repository {
	return &Repository{
		Db: db,
	}
}

func (r *Repository) getUserByField(ctx context.Context, field, value string) (*domain.User, error) {
	var user user
	err := r.Db.Users.FindOne(
		ctx,
		bson.D{{Key: field, Value: value}},
	).Decode(&user)
	if err != nil {
		if err == drv.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	du := user.ToDomain()

	return &du, nil
}

func (r *Repository) GetUserByID(ctx context.Context, userID string) (*domain.User, error) {
	var user user
	objUserID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	err = r.Db.Users.FindOne(
		ctx,
		bson.D{{Key: "_id", Value: objUserID}},
	).Decode(&user)
	if err != nil {
		if err == drv.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	du := user.ToDomain()

	return &du, nil
}

func (r *Repository) GetUserByName(ctx context.Context, name string) (*domain.User, error) {
	return r.getUserByField(ctx, nameField, name)
}

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	return r.getUserByField(ctx, emailField, email)
}

func (r *Repository) AddUser(ctx context.Context, name, email, password string) (string, error) {
	user := user{
		ID:       primitive.NewObjectID(),
		Name:     name,
		Email:    email,
		Password: password,
	}
	res, err := r.Db.Users.InsertOne(
		ctx,
		user,
	)
	if err != nil {
		return "", err
	}

	id := res.InsertedID.(primitive.ObjectID).Hex()

	return id, nil
}
