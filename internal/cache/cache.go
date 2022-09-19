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

func (c *Cache) set(ctx context.Context, key string, value interface{}, dur time.Duration) error {
	return c.redis.Db.Set(
		ctx,
		key,
		value,
		dur,
	).Err()
}

func (c *Cache) SetAccessToken(ctx context.Context, userID, token string) error {
	key := fmt.Sprintf("%s:%s:%s", whiteListKey, userID, token)
	return c.set(ctx, key, nil, c.accessDuration)
}

func (c *Cache) SetRefreshToken(ctx context.Context, userID, token string) error {
	key := fmt.Sprintf("%s:%s:%s", refreshKey, userID, token)
	return c.set(ctx, key, nil, c.refreshDuration)
}
