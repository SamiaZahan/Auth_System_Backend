package service

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"time"

	"github.com/emamulandalib/airbringr-auth/config"
	"github.com/emamulandalib/airbringr-auth/dto"
	"github.com/emamulandalib/airbringr-auth/repository"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
)

func (a *Auth) Signup(input dto.SignupInput) (err error) {
	genericSignupFailureMsg := errors.New("Signup failed for some technical reason.")
	ctx := context.Background()
	authRepo := repository.Auth{Ctx: ctx}
	passwordService := PasswordService{}

	// try to get existing user
	existingUser, err := authRepo.GetUserByEmail(input.Email)
	if err != nil && err != mongo.ErrNoDocuments {
		log.Error(err.Error())
		return genericSignupFailureMsg
	}
	if existingUser != nil {
		return errors.New("An user with this email already exists")
	}

	userExistResponse := ExistingEmail(input.Email)
	if userExistResponse.Status {
		return errors.New("An user with this email already exists")
	}

	var otp string
	otpSvc := OtpSvc{MicroAPIToken: config.Params.MicroAPIToken}
	if otp, err = otpSvc.Generate(OtpGenerateRequest{
		Expiry: int64(time.Hour * 24),
		Id:     input.Email,
	}); err != nil {
		//fmt.Print(otp)
		return genericSignupFailureMsg
	}
	createVerificationLink := func(sessCtx mongo.SessionContext) (i interface{}, err error) {
		var userID string
		AuthRpo := repository.Auth{Ctx: sessCtx}
		hashedPassword := passwordService.HashPassword(input.Password)
		if userID, err = AuthRpo.CreateUser(input.Email, hashedPassword); err != nil {
			return
		}
		if err = AuthRpo.CreateUserProfile(userID, input.FirstName, input.LastName); err != nil {
			return
		}
		if err = a.SendEmail(input.Email, otp); err != nil {
			return
		}
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
		return errors.New("Email send failed")
	}
	return nil
}
