package bcrypt

import "testing"

func TestBcryptModule(t *testing.T) {
	password := "!curved_cucumber_0!"
	b := NewBcrypt(10)

	hash, err := b.HashPassword(password)
	if err != nil {
		panic(err)
	}

	ok, err := b.CompareHashAndPassword(hash, password)
	if err != nil {
		panic(err)
	}

	if !ok {
		t.Error("passwords are not equal")
	}
}
