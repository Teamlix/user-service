package redis

import (
	"context"
	"fmt"
	"time"

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

func (c *Redis) SetKey(ctx context.Context, key string, value interface{}, dur time.Duration) error {
	return c.Db.Set(
		ctx,
		key,
		value,
		dur,
	).Err()
}

func (c *Redis) CheckKey(ctx context.Context, key string) (bool, error) {
	_, err := c.Db.Get(ctx, key).Result()
	if err != nil {
		if err == drv.Nil {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (c *Redis) DeleteKey(ctx context.Context, key string) error {
	err := c.Db.Del(
		ctx,
		key,
	).Err()

	if err != nil {
		if err == drv.Nil {
			return nil
		}
		return err
	}

	return nil
}

func (r *Redis) Disconnect() error {
	return r.Db.Close()
}
