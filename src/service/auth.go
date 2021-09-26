package service

import (
	"github.com/emamulandalib/airbringr-auth/dto"
	"github.com/emamulandalib/airbringr-auth/repository"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type Auth struct{}

func (a *Auth) Signup(input dto.SignupInput) (err error) {
	authRepo := repository.Auth{}
	ID, err := authRepo.CreateUser(input.Email)

	if err != nil {
		log.Error(err.Error())
		return errors.New("Signup failed for some technical reason.")
	}

	err = authRepo.CreateUserProfile(ID, input.FirstName, input.LastName)

	if err != nil {
		log.Error(err.Error())
		return errors.New("Signup failed for some technical reason.")
	}

	// delete user

	// send email

	return
}
