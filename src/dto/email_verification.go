package dto

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type EmailVerificationInput struct {
	Auth string `json:"auth"`
	OTP  int    `json:"otp"`
}

func (input EmailVerificationInput) Validate() error {
	return validation.ValidateStruct(&input,
		validation.Field(&input.Auth, validation.Required),
		validation.Field(&input.OTP, validation.Required),
	)
}
