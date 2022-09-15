package redis

import (
	"context"
	"fmt"

	drv "github.com/go-redis/redis/v9"
)

type Redis struct {
	Db *drv.Client
}

func NewRedis(ctx context.Context, host, port, password string, db int) (*Redis, error) {
	rdb := drv.NewClient(&drv.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: password,
		DB:       db,
	})

	err := rdb.Ping(context.Background()).Err()
	if err != nil {
		return nil, fmt.Errorf("error connect to redis: %w", err)
	}

	return &Redis{
		Db: rdb,
	}, nil
}

func (r *Redis) Disconnect() error {
	return r.Db.Close()
}
