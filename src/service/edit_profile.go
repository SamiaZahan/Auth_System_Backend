package service

import (
	"context"
	"errors"
	"github.com/emamulandalib/airbringr-auth/dto"
	"github.com/emamulandalib/airbringr-auth/repository"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type EditProfile struct{}

func (ep EditProfile) EditUserProfile(input *dto.EditProfileInput, email string) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	//genericEditFailureMsg := errors.New("Profile Edit failed for some technical reason.")

	aRepo := repository.Auth{Ctx: ctx}
	user, err := aRepo.GetUserByEmail(email)
	if err != nil {
		return errors.New("user not found")
	}
	err = aRepo.SetUserProfileByID(user.ID.Hex(), input)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	if err != nil && mongo.ErrNoDocuments != err {
		log.Error(err)
		return errors.New("user not found")
	}
	return nil

}
