package dto

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type EmailOtpInput struct {
	Email string `json:"email"`
}

func (input EmailOtpInput) Validate() error {
	return validation.ValidateStruct(&input,
		validation.Field(&input.Email, validation.Required, is.Email),
	)
}
