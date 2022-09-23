package repository

import (
	"context"

	"github.com/teamlix/user-service/internal/domain"
	"github.com/teamlix/user-service/internal/pkg/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	drv "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func (r *Repository) GetUsersTotalCount(ctx context.Context) (int, error) {
	cnt, err := r.Db.Users.CountDocuments(ctx, bson.D{})
	if err != nil {
		return 0, err
	}

	return int(cnt), nil
}

func (r *Repository) GetUsers(ctx context.Context, skip, limit int) ([]domain.User, error) {
	res := make([]domain.User, 0)
	opts := options.Find().
		SetSort(bson.D{{Key: "_id", Value: -1}}).
		SetSkip(int64(skip)).
		SetLimit(int64(limit))

	cur, err := r.Db.Users.Find(
		ctx,
		bson.D{},
		opts,
	)
	if err != nil {
		return res, err
	}

	for cur.Next(ctx) {
		var result user
		err = cur.Decode(&result)
		if err != nil {
			return res, err
		}
		res = append(res, result.ToDomain())
	}
	if err = cur.Err(); err != nil {
		return res, err
	}

	err = cur.Close(ctx)
	if err != nil {
		return res, err
	}

	return res, nil
}
