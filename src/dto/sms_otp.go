package dto

import validation "github.com/go-ozzo/ozzo-validation/v4"

type SendSmsOtpInput struct {
	Mobile string `json:"mobile"`
}

func (input SendSmsOtpInput) Validate() error {
	return validation.ValidateStruct(&input,
		validation.Field(&input.Mobile, validation.Required),
	)
}
