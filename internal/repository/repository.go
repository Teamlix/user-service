package repository

import "github.com/teamlix/user-service/internal/pkg/mongo"

type Repository struct {
	Db *mongo.Mongo
}

func NewRepository(db *mongo.Mongo) *Repository {
	return &Repository{
		Db: db,
	}
}
