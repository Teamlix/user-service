package mongo

import (
	"context"
	"fmt"

	drv "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
)

const usersCollection = "users"

type Mongo struct {
	client *drv.Client
	Users  *drv.Collection
}

func NewMongo(ctx context.Context, connString string) (*Mongo, error) {
	dburl, err := connstring.Parse(connString)
	if err != nil {
		return nil, fmt.Errorf("error parsing mongo url: %w", err)
	}

	client, err := drv.Connect(ctx, options.Client().ApplyURI(connString))
	if err != nil {
		return nil, err
	}

	db := client.Database(dburl.Database)
	users := db.Collection(usersCollection)

	return &Mongo{
		client: client,
		Users:  users,
	}, nil
}

func (m *Mongo) Disconnect(ctx context.Context) error {
	return m.client.Disconnect(ctx)
}
