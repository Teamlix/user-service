package validator

type Validator struct {
}

func NewValidator() *Validator {
	return &Validator{}
}

func (v *Validator) validateEmail(email string) error {
	return nil
}

func (v *Validator) validateName(name string) error {
	return nil
}

func (v *Validator) ValidateSignUp(email, name, password, repeatedPassword string) error {
	return nil
}

func (v *Validator) ValidateSignIn(email, password, repeatedPassword string) error {
	return nil
}
