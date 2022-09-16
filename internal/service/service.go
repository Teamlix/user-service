package service

type Repository interface {
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
