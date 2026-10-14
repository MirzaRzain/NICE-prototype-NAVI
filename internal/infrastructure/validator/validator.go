package validator

import "github.com/go-playground/validator/v10"

func NewValidator() *Validator {
	return &Validator{Validator: validator.New()}
}

type Validator struct {
	Validator *validator.Validate
}

func (v Validator) Validate(i interface{}) (err error) {
	err = v.Validator.Struct(i)

	return
}
