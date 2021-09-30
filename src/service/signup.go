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

func (a *Auth) Signup(input dto.SignupInput) (err error) {
	genericSignupFailureMsg := errors.New("Signup failed for some technical reason.")
	ctx := context.Background()
	authRepo := repository.Auth{Ctx: ctx}

	// try to get existing user
	existingUser, err := authRepo.GetUserByEmail(input.Email)
	if err != nil && err != mongo.ErrNoDocuments {
		log.Error(err.Error())
		return genericSignupFailureMsg
	}
	if existingUser != nil {
		return errors.New("An user with this email is already exist.")
	}

	if a.ExistingEmail(input.Email) {
		return errors.New("An user with this email is already exist.")
	}

	// try to create user, user profile and verification link
	otp := GenerateRandNum()
	createVerificationLink := func(sessCtx mongo.SessionContext) (i interface{}, err error) {
		var userID string
		var verificationID string
		AuthRpo := repository.Auth{Ctx: sessCtx}
		VerificationRepo := repository.Verification{Ctx: sessCtx}

		if userID, err = AuthRpo.CreateUser(input.Email); err != nil {
			return
		}
		if err = AuthRpo.CreateUserProfile(userID, input.FirstName, input.LastName); err != nil {
			return
		}
		if verificationID, err = VerificationRepo.Create(input.Email, otp, userID); err != nil {
			return
		}

		err = a.SendEmail(input.Email, otp, verificationID)
		return
	}

	var sess mongo.Session
	if sess, err = repository.MongoClient.StartSession(); err != nil {
		log.Error(err.Error())
		return genericSignupFailureMsg
	}
	defer sess.EndSession(ctx)

	if _, err = sess.WithTransaction(ctx, createVerificationLink); err != nil {
		log.Error(err.Error())
		return genericSignupFailureMsg
	}

	return
}

func (a *Auth) ExistingEmail(email string) (exists bool) {
	if code, _, errs := fiber.
		Post(fmt.Sprintf("%s/helper/exist-email", config.Params.AirBringrDomain)).
		JSON(fiber.Map{
			"email": email,
		}).
		String(); code != fiber.StatusOK {
		log.Error(errs)
		exists = true
		return
	}

	exists = false
	return
}

func (a *Auth) SendEmail(email string, otp int, verificationID string) error {
	emailSvcURI := fmt.Sprintf("%s/v1/send-email", config.Params.NotificationSvcDomain)
	verificationLink := fmt.Sprintf("%s/verification/?otp=%d&auth=%s", config.Params.ServiceFrontend, otp, verificationID)

	code, _, errs := fiber.
		Post(emailSvcURI).
		JSON(fiber.Map{
			"data": fiber.Map{
				"link": verificationLink,
			},
			"to":            email,
			"from":          "contact@airbringr.com",
			"message":       "Please click the link to verify your account.",
			"subject":       "AirBringr Signup Verification",
			"template_code": "signup_verification",
		}).
		String()

	if code != fiber.StatusOK {
		log.Error(errs)
		return errors.New("Email send failed.")
	}
	return nil
}
