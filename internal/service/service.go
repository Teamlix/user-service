package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/teamlix/user-service/internal/domain"
	"golang.org/x/sync/errgroup"
)

const userRole = 1

type Repository interface {
	GetUserByID(ctx context.Context, userID string) (*domain.User, error)
	GetUserByName(ctx context.Context, name string) (*domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	AddUser(ctx context.Context, name, email, password string) (string, error)
	GetUsersTotalCount(ctx context.Context) (int, error)
	GetUsers(ctx context.Context, skip, limit int) ([]domain.User, error)
}

type Cache interface {
	SetAccessToken(ctx context.Context, userID, token string) error
	SetRefreshToken(ctx context.Context, userID, token string) error
	CheckAccessToken(ctx context.Context, userID, token string) (bool, error)
	CheckRefreshToken(ctx context.Context, userID, token string) (bool, error)
	RemoveAccessToken(ctx context.Context, userID, token string) error
	RemoveRefreshToken(ctx context.Context, userID, token string) error
}

type Bcrypt interface {
	HashPassword(password string) (string, error)
	CompareHashAndPassword(hash, password string) (bool, error)
}

type Validator interface {
	ValidateSignUp(email, name, password, repeatedPassword string) error
	ValidateSignIn(email, password string) error
	ValidateGetUserByID(userID string) error
	ValidateGetUsersList(skip, limit int) error
}

type Tokener interface {
	SignAccessToken(userID string, role uint32) (string, error)
	SignRefreshToken(userID string) (string, error)
	ValidateAccessToken(token string) (string, uint32, error)
	ValidateRefreshToken(token string) (string, error)
}

type Service struct {
	repository Repository
	cache      Cache
	bcrypt     Bcrypt
	validator  Validator
	tokens     Tokener
}

func NewService(
	repo Repository,
	cache Cache,
	b Bcrypt,
	v Validator,
	t Tokener,
) *Service {
	return &Service{
		repository: repo,
		cache:      cache,
		bcrypt:     b,
		validator:  v,
		tokens:     t,
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

	hashedPassword, err := s.bcrypt.HashPassword(password)
	if err != nil {
		return t, err
	}

	id, err := s.repository.AddUser(ctx, name, email, hashedPassword)
	if err != nil {
		return t, err
	}

	eg = errgroup.Group{}
	var accessToken, refreshToken string
	eg.Go(func() error {
		accessToken, err = s.tokens.SignAccessToken(id, userRole)
		return err
	})
	eg.Go(func() error {
		refreshToken, err = s.tokens.SignRefreshToken(id)
		return err
	})
	err = eg.Wait()
	if err != nil {
		return t, err
	}

	eg = errgroup.Group{}
	eg.Go(func() error {
		return s.cache.SetAccessToken(ctx, id, accessToken)
	})
	eg.Go(func() error {
		return s.cache.SetRefreshToken(ctx, id, refreshToken)
	})
	err = eg.Wait()
	if err != nil {
		return t, err
	}

	t.AccessToken = accessToken
	t.RefreshToken = refreshToken

	return t, nil
}

func (s *Service) SignIn(ctx context.Context, email, password string) (domain.Tokens, error) {
	var t domain.Tokens

	err := s.validator.ValidateSignIn(email, password)
	if err != nil {
		return t, err
	}

	user, err := s.repository.GetUserByEmail(ctx, email)
	if err != nil {
		return t, err
	}
	if user == nil {
		return t, errors.New("user not found")
	}

	ok, err := s.bcrypt.CompareHashAndPassword(user.Password, password)
	if err != nil && !ok {
		return t, errors.New("user not found")
	}

	eg := errgroup.Group{}
	var accessToken, refreshToken string
	eg.Go(func() error {
		accessToken, err = s.tokens.SignAccessToken(user.ID, userRole)
		return err
	})
	eg.Go(func() error {
		refreshToken, err = s.tokens.SignRefreshToken(user.ID)
		return err
	})
	err = eg.Wait()
	if err != nil {
		return t, err
	}

	eg = errgroup.Group{}
	eg.Go(func() error {
		return s.cache.SetAccessToken(ctx, user.ID, accessToken)
	})
	eg.Go(func() error {
		return s.cache.SetRefreshToken(ctx, user.ID, refreshToken)
	})
	err = eg.Wait()
	if err != nil {
		return t, err
	}

	t.AccessToken = accessToken
	t.RefreshToken = refreshToken

	return t, nil
}

func (s *Service) Refresh(ctx context.Context, rt string) (domain.Tokens, error) {
	t := domain.Tokens{}

	userID, err := s.tokens.ValidateRefreshToken(rt)
	if err != nil {
		return t, err
	}

	if userID == "" {
		return t, errors.New("user not found")
	}

	ok, err := s.cache.CheckRefreshToken(ctx, userID, rt)
	if err != nil {
		return t, err
	}
	if !ok {
		return t, errors.New("unauthorized")
	}

	user, err := s.repository.GetUserByID(ctx, userID)
	if err != nil {
		return t, err
	}
	if user == nil {
		return t, errors.New("user not found")
	}

	eg := errgroup.Group{}
	var accessToken, refreshToken string
	eg.Go(func() error {
		accessToken, err = s.tokens.SignAccessToken(userID, userRole)
		return err
	})
	eg.Go(func() error {
		refreshToken, err = s.tokens.SignRefreshToken(userID)
		return err
	})
	err = eg.Wait()
	if err != nil {
		return t, err
	}

	eg = errgroup.Group{}
	eg.Go(func() error {
		return s.cache.RemoveRefreshToken(ctx, userID, rt)
	})
	eg.Go(func() error {
		return s.cache.SetAccessToken(ctx, userID, accessToken)
	})
	eg.Go(func() error {
		return s.cache.SetRefreshToken(ctx, userID, refreshToken)
	})
	err = eg.Wait()
	if err != nil {
		return t, err
	}

	t.AccessToken = accessToken
	t.RefreshToken = refreshToken

	return t, nil
}

func (s *Service) LogOut(ctx context.Context, accessToken, refreshToken string) error {
	userID, err := s.tokens.ValidateRefreshToken(refreshToken)
	if err != nil {
		return err
	}

	if userID == "" {
		return errors.New("user not found")
	}

	ok, err := s.cache.CheckRefreshToken(ctx, userID, refreshToken)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("unauthorized")
	}

	eg := errgroup.Group{}
	eg.Go(func() error {
		return s.cache.RemoveAccessToken(ctx, userID, accessToken)
	})
	eg.Go(func() error {
		return s.cache.RemoveRefreshToken(ctx, userID, refreshToken)
	})
	err = eg.Wait()
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) GetUserByID(ctx context.Context, userID string) (domain.User, error) {
	u := domain.User{}

	if err := s.validator.ValidateGetUserByID(userID); err != nil {
		return u, err
	}

	user, err := s.repository.GetUserByID(ctx, userID)
	if err != nil {
		return u, err
	}
	if user == nil {
		return u, errors.New("user not found")
	}

	return *user, nil
}

func (s *Service) GetUsersList(ctx context.Context, skip, limit int) ([]domain.User, int, error) {
	users := make([]domain.User, 0)
	cnt := 0

	if err := s.validator.ValidateGetUsersList(skip, limit); err != nil {
		return users, cnt, err
	}

	eg := errgroup.Group{}
	eg.Go(func() error {
		ul, err := s.repository.GetUsers(ctx, skip, limit)
		if err != nil {
			return err
		}
		users = ul
		return nil
	})
	eg.Go(func() error {
		c, err := s.repository.GetUsersTotalCount(ctx)
		if err != nil {
			return err
		}
		cnt = c
		return nil
	})
	if err := eg.Wait(); err != nil {
		return users, cnt, err
	}

	return users, cnt, nil
}

func (s *Service) CheckAccessToken(ctx context.Context, accessToken string) error {
	userID, _, err := s.tokens.ValidateAccessToken(accessToken)
	if err != nil {
		return err
	}

	if userID == "" {
		return errors.New("user not found")
	}

	ok, err := s.cache.CheckAccessToken(ctx, userID, accessToken)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("unauthorized")
	}
	return nil
}
