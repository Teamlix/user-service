package service

import (
	"github.com/teamlix/user-service/internal/cache"
	"github.com/teamlix/user-service/internal/pkg/bcrypt"
	"github.com/teamlix/user-service/internal/pkg/validator"
	"github.com/teamlix/user-service/internal/repository"
)

type Service struct {
	repository *repository.Repository
	cache      *cache.Cache
	bcrypt     *bcrypt.Bcrypt
	validator  *validator.Validator
}

func NewService(
	repo *repository.Repository,
	cache *cache.Cache,
	b *bcrypt.Bcrypt,
	v *validator.Validator,
) *Service {
	return &Service{
		repository: repo,
		cache:      cache,
		bcrypt:     b,
	}
}
