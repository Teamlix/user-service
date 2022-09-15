package service

import (
	"github.com/teamlix/user-service/internal/cache"
	"github.com/teamlix/user-service/internal/pkg/bcrypt"
	"github.com/teamlix/user-service/internal/repository"
)

type Service struct {
	repository *repository.Repository
	cache      *cache.Cache
	bcrypt     *bcrypt.Bcrypt
}

func NewService(repo *repository.Repository, cache *cache.Cache, b *bcrypt.Bcrypt) *Service {
	return &Service{
		repository: repo,
		cache:      cache,
		bcrypt:     b,
	}
}
