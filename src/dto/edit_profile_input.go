package dto

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type EditProfileInput struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Gender    string `json:"gender"`
	Address   string `json:"address"`
}

func (input EditProfileInput) Validate() error {
	return validation.ValidateStruct(&input,
		validation.Field(&input.FirstName, validation.Length(2, 10), is.Alpha),
		validation.Field(&input.LastName, validation.Length(2, 10), is.Alpha),
		validation.Field(&input.Gender, is.Alpha),
		validation.Field(&input.Address),
	)
}
