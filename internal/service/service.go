package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/teamlix/user-service/internal/domain"
	"golang.org/x/sync/errgroup"
)

type Repository interface {
	GetUserByName(ctx context.Context, name string) (*domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
}

type Cache interface {
}

type Bcrypt interface {
}

type Validator interface {
	ValidateSignUp(email, name, password, repeatedPassword string) error
	ValidateSignIn(email, password, repeatedPassword string) error
}

type Service struct {
	repository Repository
	cache      Cache
	bcrypt     Bcrypt
	validator  Validator
}

func NewService(
	repo Repository,
	cache Cache,
	b Bcrypt,
	v Validator,
) *Service {
	return &Service{
		repository: repo,
		cache:      cache,
		bcrypt:     b,
	}
}

func (s *Service) SignUp(ctx context.Context, name, email, password, repeatedPassword string) (domain.Tokens, error) {
	t := domain.Tokens{}
	err := s.validator.ValidateSignUp(email, name, password, repeatedPassword)
	if err != nil {
		return t, err
	}

	var nameExists, emailExists bool
	eg := errgroup.Group{}
	eg.Go(func() error {
		u, err := s.repository.GetUserByEmail(ctx, email)
		if err != nil {
			return err
		}
		if u != nil {
			emailExists = true
		}
		return nil
	})
	eg.Go(func() error {
		u, err := s.repository.GetUserByName(ctx, name)
		if err != nil {
			return err
		}
		if u != nil {
			nameExists = true
		}
		return nil
	})

	err = eg.Wait()
	if err != nil {
		return t, fmt.Errorf("error get user: %w", err)
	}

	if emailExists {
		return t, errors.New("provided email already exists")
	}
	if nameExists {
		return t, errors.New("provided user name already exists")
	}

	return t, nil
}
