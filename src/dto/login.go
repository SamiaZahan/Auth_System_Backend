package dto

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"regexp"
	"strings"
)

type LoginInput struct {
	EmailOrMobile string `json:"email_or_mobile"`
	Password      string `json:"password"`
}

func (input LoginInput) Validate() error {
	return validation.ValidateStruct(&input,
		validation.Field(&input.Password, validation.Required),
		validation.Field(&input.EmailOrMobile, validation.
			When(strings.Contains(input.EmailOrMobile, "@"), is.Email).
			Else(validation.Match(regexp.MustCompile("^(\\+[1-9]{1})?[0-9]{4,14}$")))),
	)
}
