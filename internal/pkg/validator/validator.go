package validator

type Validator struct {
}

func NewValidator() *Validator {
	return &Validator{}
}

func (v *Validator) ValidateEmail(email string) error {
	return nil
}

func (v *Validator) ValidateName(name string) error {
	return nil
}
