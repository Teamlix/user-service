package bcrypt

import (
	bLib "golang.org/x/crypto/bcrypt"
)

type Bcrypt struct {
	cost int
}

func NewBcrypt(cost int) *Bcrypt {
	return &Bcrypt{cost: cost}
}

func (b *Bcrypt) HashPassword(password string) (string, error) {
	bytes, err := bLib.GenerateFromPassword([]byte(password), b.cost)
	if err != nil {
		return "", err
	}

	hash := string(bytes)

	return hash, err
}

func (b *Bcrypt) CompareHashAndPassword(hash, password string) (bool, error) {
	err := bLib.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return false, err
	}

	return true, nil
}
