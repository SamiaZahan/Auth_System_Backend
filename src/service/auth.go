package service

import (
	"context"
	"fmt"
	"github.com/emamulandalib/airbringr-auth/config"
	"github.com/emamulandalib/airbringr-auth/dto"
	"github.com/emamulandalib/airbringr-auth/repository"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

type Auth struct{}

func (a *Auth) Signup(input dto.SignupInput) (err error) {
	genericSignupFailureMsg := errors.New("Signup failed for some technical reason.")
	ctx := context.Background()
	authRepo := repository.Auth{ctx}

	// try to get existing user
	existingUser, err := authRepo.GetUserByEmail(input.Email)
	if err != nil && err != mongo.ErrNoDocuments {
		log.Error(err.Error())
		return genericSignupFailureMsg
	}
	if existingUser != nil {
		return errors.New("An user with this email is already exist.")
	}

	// try to create user, user profile and verification link
	otp := GenerateRandNum()
	createVerificationLink := func(sessCtx mongo.SessionContext) (i interface{}, err error) {
		var userID string
		AuthRpo := repository.Auth{sessCtx}
		VerificationRepo := repository.Verification{sessCtx}

		if userID, err = AuthRpo.CreateUser(input.Email); err != nil {
			return
		}
		if err = AuthRpo.CreateUserProfile(userID, input.FirstName, input.LastName); err != nil {
			return
		}
		if err = VerificationRepo.Create(input.Email, otp, userID); err != nil {
			return
		}

		err = a.SendEmail(input.Email, otp)
		return
	}

	sess, err := repository.MongoClient.StartSession()
	if err != nil {
		log.Error(err.Error())
		return genericSignupFailureMsg
	}
	defer sess.EndSession(ctx)

	_, err = sess.WithTransaction(ctx, createVerificationLink)
	if err != nil {
		log.Error(err.Error())
		return genericSignupFailureMsg
	}

	return
}

func (a *Auth) SendEmail(email string, otp int) error {
	emailSvcURI := fmt.Sprintf("%s/v1/send-email", config.Params.NotificationSvcDomain)
	verificationLink := fmt.Sprintf("%s/verification/?otp=%s", config.Params.MainDomain, otp)

	code, _, errs := fiber.
		Post(emailSvcURI).
		JSON(fiber.Map{
			"data": fiber.Map{
				"link": verificationLink,
			},
			"to":            email,
			"from":          "contact@airbringr.com",
			"message":       "Please click the link to verify your account.",
			"subject":       "AirBringr Signup VerificationDoc",
			"template_code": "signup_verification",
		}).
		String()

	if code != fiber.StatusOK {
		log.Error(errs)
		return errors.New("Email send failed.")
	}
	return nil
}
