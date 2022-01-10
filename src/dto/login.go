package dto

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"strings"
)

type LoginInput struct {
	EmailOrMobile string `json:"email_or_mobile"`
	Password      string `json:"password"`
	CountryPrefix string `json:"country_prefix"`
}

func (input LoginInput) Validate() error {
	return validation.ValidateStruct(&input,
		validation.Field(&input.Password, validation.Required),
		validation.Field(&input.EmailOrMobile, validation.
			When(strings.Contains(input.EmailOrMobile, "@"), is.Email).
			Else(is.Digit)),
		validation.Field(&input.CountryPrefix, validation.
			When(!strings.Contains(input.EmailOrMobile, "@"), validation.Required), is.Digit),
	)
}
