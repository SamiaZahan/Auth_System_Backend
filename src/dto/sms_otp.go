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

type VerifyOtpInput struct {
	EmailOrMobile string `json:"email_or_mobile"`
	OTP           int    `json:"otp"`
}

func (input VerifyOtpInput) Validate() error {
	return validation.ValidateStruct(&input,
		validation.Field(&input.EmailOrMobile, validation.Required),
		validation.Field(&input.OTP, validation.Required),
	)
}

type VerifyMobileInput struct {
	Auth   string `json:"auth"`
	Mobile string `json:"mobile"`
	OTP    int    `json:"otp"`
}

func (input VerifyMobileInput) Validate() error {
	return validation.ValidateStruct(&input,
		validation.Field(&input.Auth, validation.Required),
		validation.Field(&input.Mobile, validation.Required),
		validation.Field(&input.OTP, validation.Required),
	)
}
