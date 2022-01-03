package service

import (
	"context"
	"fmt"
	"time"

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

	if ExistingEmail(input.Email) {
		return errors.New("An user with this email is already exist.")
	}

	var otp string
	otpSvc := OtpSvc{MicroAPIToken: config.Params.MicroAPIToken}
	if otp, err = otpSvc.Generate(OtpGenerateRequest{
		Expiry: int64(time.Hour * 24),
		Id:     input.Email,
	}); err != nil {
		return genericSignupFailureMsg
	}

	createVerificationLink := func(sessCtx mongo.SessionContext) (i interface{}, err error) {
		var userID string
		AuthRpo := repository.Auth{Ctx: sessCtx}
		hashedPassword, passwordHasingError := authRepo.HashPassword(input.Password)
		if err != nil {
			log.Error(passwordHasingError.Error())
			return
		}
		if userID, err = AuthRpo.CreateUser(input.Email, hashedPassword); err != nil {
			return
		}
		if err = AuthRpo.CreateUserProfile(userID, input.FirstName, input.LastName); err != nil {
			return
		}

		err = a.SendEmail(input.Email, otp)
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

func (a *Auth) SendEmail(email string, otp string) error {
	emailSvcURI := fmt.Sprintf("%s/v1/send-email", config.Params.NotificationSvcDomain)
	verificationLink := fmt.Sprintf("%s/verification/?otp=%s&auth=%s", config.Params.ServiceFrontend, otp, email)

	if code, _, errs := fiber.
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
		String(); code != fiber.StatusOK {
		log.Error(errs)
		return errors.New("Email send failed.")
	}
	return nil
}
