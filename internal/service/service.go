package service

const userRole = 1

type Repository interface {
}

type Cache interface {
}

type Validator interface {
}

type Service struct {
	repository Repository
	cache      Cache
	validator  Validator
}

func NewService(
	repo Repository,
	cache Cache,
	v Validator,
) *Service {
	return &Service{
		repository: repo,
		cache:      cache,
		validator:  v,
	}
}
