package dto

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type SignupInput struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

func (input SignupInput) Validate() error {
	return validation.ValidateStruct(&input,
		validation.Field(&input.FirstName, validation.Required),
		validation.Field(&input.LastName, validation.Required),
		validation.Field(&input.Email, validation.Required, is.Email),
	)
}
