package dto

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type EditProfileInput struct {
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Gender    string  `json:"gender"`
	Age       string  `json:"age"`
	Address   Address `json:"address"`
}
type Address struct {
	Division string `json:"division"`
	District string `json:"district"`
	Area     string `json:"area"`
	Text     string `json:"text"`
	Zone     string `json:"zone"`
}

func (input EditProfileInput) Validate() error {
	return validation.ValidateStruct(&input,
		validation.Field(&input.FirstName, validation.Length(2, 10), is.Alpha),
		validation.Field(&input.LastName, validation.Length(2, 10), is.Alpha),
		validation.Field(&input.Gender, is.Alpha),
		validation.Field(&input.Age),
		validation.Field(&input.Address),
	)
}
