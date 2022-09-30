package cache

import (
	"github.com/teamlix/user-service/internal/pkg/redis"
)

type Cache struct {
	redis *redis.Redis
}

func NewCache(redis *redis.Redis) *Cache {
	return &Cache{
		redis: redis,
	}
}
