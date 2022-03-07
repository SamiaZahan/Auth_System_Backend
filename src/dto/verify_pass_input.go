package dto

import validation "github.com/go-ozzo/ozzo-validation/v4"

type VerifyPassword struct {
	Password string `json:"password"`
}

func (input VerifyPassword) Validate() error {
	return validation.ValidateStruct(&input,
		validation.Field(&input.Password, validation.Required),
	)
}
