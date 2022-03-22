package dto

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type EditEmailInput struct {
	Email string `json:"email"`
}

func (input EditEmailInput) Validate() error {
	return validation.ValidateStruct(&input,
		validation.Field(&input.Email, validation.Required, is.Email),
	)
}
