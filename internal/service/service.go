package service

import (
	"github.com/teamlix/user-service/internal/cache"
	"github.com/teamlix/user-service/internal/repository"
)

type Service struct {
	Repository *repository.Repository
	Cache      *cache.Cache
}

func NewService(repo *repository.Repository, cache *cache.Cache) *Service {
	return &Service{
		Repository: repo,
		Cache:      cache,
	}
}
