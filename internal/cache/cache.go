package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/teamlix/user-service/internal/pkg/redis"
)

const (
	whiteListKey = "white_list"
	refreshKey   = "refresh"
)

type Cache struct {
	redis           *redis.Redis
	accessDuration  time.Duration
	refreshDuration time.Duration
}

func NewCache(redis *redis.Redis, accessDuration, refreshDuration time.Duration) *Cache {
	return &Cache{
		redis:           redis,
		accessDuration:  accessDuration,
		refreshDuration: refreshDuration,
	}
}

func (c *Cache) SetAccessToken(ctx context.Context, userID, token string) error {
	key := fmt.Sprintf("%s:%s:%s", whiteListKey, userID, token)
	return c.redis.SetKey(ctx, key, nil, c.accessDuration)
}

func (c *Cache) SetRefreshToken(ctx context.Context, userID, token string) error {
	key := fmt.Sprintf("%s:%s:%s", refreshKey, userID, token)
	return c.redis.SetKey(ctx, key, nil, c.refreshDuration)
}

func (c *Cache) CheckAccessToken(ctx context.Context, userID, token string) (bool, error) {
	key := fmt.Sprintf("%s:%s:%s", whiteListKey, userID, token)
	return c.redis.CheckKey(ctx, key)
}

func (c *Cache) CheckRefreshToken(ctx context.Context, userID, token string) (bool, error) {
	key := fmt.Sprintf("%s:%s:%s", refreshKey, userID, token)
	return c.redis.CheckKey(ctx, key)
}

func (c *Cache) RemoveRefreshToken(ctx context.Context, userID, token string) error {
	key := fmt.Sprintf("%s:%s:%s", refreshKey, userID, token)
	return c.redis.DeleteKey(ctx, key)
}
