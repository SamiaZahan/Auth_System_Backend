package dto

import validation "github.com/go-ozzo/ozzo-validation/v4"

type PasswordReset struct {
	Auth     string `json:"auth"`
	Password string `json:"password"`
}

func (input PasswordReset) Validate() error {
	return validation.ValidateStruct(&input,
		validation.Field(&input.Auth, validation.Required),
		validation.Field(&input.Password, validation.Required),
	)
}
