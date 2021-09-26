package dto

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"strings"
)

type LoginInput struct {
	UsernameOrEmail string `json:"username_or_email"`
}

func (input LoginInput) Validate() error {
	return validation.ValidateStruct(&input,
		validation.Field(&input.UsernameOrEmail,
			validation.Required,
			validation.When(
				strings.Contains(input.UsernameOrEmail, "@"), is.Email,
			),
		),
	)
}
