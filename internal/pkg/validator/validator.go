package validator

import (
	"errors"
	"net/mail"
	"unicode/utf8"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	maxUserNameLen = 50
	minPasswordLen = 5
	maxPasswordLen = 30
)

type Validator struct {
}

func NewValidator() *Validator {
	return &Validator{}
}

func (v *Validator) validateEmail(email string) error {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return errors.New("email has wrong format")
	}

	return nil
}

func (v *Validator) validateName(name string) error {
	l := utf8.RuneCountInString(name)
	if l == 0 || l > maxUserNameLen {
		return errors.New("username has wrong format")
	}

	return nil
}

func (v *Validator) validatePassword(password string) error {
	l := utf8.RuneCountInString(password)
	if l == 0 || l > maxPasswordLen || l < minPasswordLen {
		return errors.New("password has wrong length")
	}

	return nil
}

func (v *Validator) ValidateSignUp(email, name, password, repeatedPassword string) error {
	err := v.validateEmail(email)
	if err != nil {
		return err
	}

	err = v.validateName(name)
	if err != nil {
		return err
	}

	err = v.validatePassword(password)
	if err != nil {
		return err
	}

	if password != repeatedPassword {
		return errors.New("passwords are not equal")
	}

	return nil
}

func (v *Validator) ValidateSignIn(email, password string) error {
	err := v.validateEmail(email)
	if err != nil {
		return err
	}

	err = v.validatePassword(password)
	if err != nil {
		return err
	}

	return nil
}

func (v *Validator) ValidateGetUserByID(userID string) error {
	ok := primitive.IsValidObjectID(userID)
	if !ok {
		return errors.New("wrong userID")
	}
	return nil
}

func (v *Validator) ValidateGetUsersList(skip, limit int) error {
	if skip < 0 || limit <= 0 {
		return errors.New("invalid payload")
	}

	return nil
}
