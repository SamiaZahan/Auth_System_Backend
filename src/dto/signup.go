package dto

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type SignupInput struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Mobile    string `json:"mobile"`
	Password  string `json:"password"`
}

func (input SignupInput) Validate() error {
	return validation.ValidateStruct(&input,
		validation.Field(&input.FirstName, validation.Required),
		validation.Field(&input.LastName, validation.Required),
		validation.Field(&input.Email, validation.Required, is.Email),
		validation.Field(&input.Mobile, validation.Required, is.Digit),
		validation.Field(&input.Password, validation.Length(8, 20)),
		//Match(regexp.
		//MustCompile("^(?=.*[0-9])(?=.*[A-Za-z]).{8,20}$")).
		//Error("Password must have minimum eight characters, at least one letter and one number")),
	)
}
