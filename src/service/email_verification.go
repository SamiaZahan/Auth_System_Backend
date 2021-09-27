package service

import (
	"context"
	"errors"

	"github.com/emamulandalib/airbringr-auth/dto"
	"github.com/emamulandalib/airbringr-auth/repository"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

func (a *Auth) EmailVerification(input dto.EmailVerificationInput) (err error) {
	genericErrMsg := errors.New("Something went wrong with the verification. Please try again later.")
	ctx := context.Background()
	vRepo := repository.Verification{Ctx: ctx}
	aRepo := repository.Auth{Ctx: ctx}

	var vDoc repository.VerificationDoc
	if vDoc, err = vRepo.GetByIDAndCode(input.Auth, input.OTP); err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New("Verification code does not exist.")
		}
		log.Error(err.Error())
		return genericErrMsg
	}

	userID := vDoc.UserID.Hex()
	if err = aRepo.ActivateUserByID(userID); err != nil {
		log.Error(err.Error())
		return genericErrMsg
	}

	_ = vRepo.DeleteByID(input.Auth)
	return
}
