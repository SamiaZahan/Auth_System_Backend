package dto

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type LoginInput struct {
	EmailOrMobile string `json:"email_or_mobile"`
	Password      string `json:"password"`
}

func (input LoginInput) Validate() error {
	return validation.ValidateStruct(&input,
		validation.Field(&input.Password, validation.Required),
		validation.Field(&input.EmailOrMobile, validation.Required),
	)
}
