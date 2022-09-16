package repository

import (
	"github.com/teamlix/user-service/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type user struct {
	ID   primitive.ObjectID `bson:"_id"`
	Name string             `bson:"name"`
}

func (u *user) ToDomain() domain.User {
	user := domain.User{
		ID:   u.ID.Hex(),
		Name: u.Name,
	}

	return user
}

func (u *user) FromDomain(input domain.User) error {
	var user user
	oId, err := primitive.ObjectIDFromHex(input.ID)
	if err != nil {
		return err
	}
	user.ID = oId
	user.Name = input.Name

	return nil
}
