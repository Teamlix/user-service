package cache

import "github.com/teamlix/user-service/internal/pkg/redis"

type Cache struct {
	Redis *redis.Redis
}

func NewCache(redis *redis.Redis) *Cache {
	return &Cache{
		Redis: redis,
	}
}
