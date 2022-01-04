package dto

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"regexp"
)

type SignupInput struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func (input SignupInput) Validate() error {
	return validation.ValidateStruct(&input,
		validation.Field(&input.FirstName, validation.Required),
		validation.Field(&input.LastName, validation.Required),
		validation.Field(&input.Email, validation.Required, is.Email),
		validation.Field(&input.Password, validation.Match(regexp.
			MustCompile("^(?=.*[0-9])(?=.*[A-Za-z]).{8,}$")).
			Error("Password must have minimum eight characters, at least one letter and one number")),
	)
}
