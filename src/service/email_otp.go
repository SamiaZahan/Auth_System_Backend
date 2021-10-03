package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/emamulandalib/airbringr-auth/config"
	"github.com/emamulandalib/airbringr-auth/dto"
	"github.com/emamulandalib/airbringr-auth/repository"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

type EmailOtp struct{}

func (e *EmailOtp) Send(input dto.EmailOtpInput) (err error) {
	genericFailureMsg := errors.New("OTP send failed")
	ctx := context.Background()
	aRepo := repository.Auth{Ctx: ctx}

	if _, err = aRepo.GetUserByEmail(input.Email); err != nil && err != mongo.ErrNoDocuments {
		return genericFailureMsg
	}

	if !ExistingEmail(input.Email) {
		return errors.New("no user found with this email. Please signup")
	}

	otp := GenerateRandNum()
	emailSvcURI := fmt.Sprintf("%s/v1/send-email", config.Params.NotificationSvcDomain)

	cb := func(sessCtx mongo.SessionContext) (d interface{}, err error) {
		vRepo := repository.Verification{Ctx: sessCtx}
		if _, err = vRepo.Create(input.Email, otp, ""); err != nil {
			return
		}

		if code, _, errs := fiber.
			Post(emailSvcURI).
			JSON(fiber.Map{
				"data": fiber.Map{
					"otp": otp,
				},
				"to":            input.Email,
				"from":          "contact@airbringr.com",
				"message":       "Please use this OTP for login into account.",
				"subject":       "AirBringr OTP",
				"template_code": "otp",
			}).
			String(); code != fiber.StatusOK {
			log.Error(errs)
			return nil, errors.New("falied to send SMS")
		}
		return
	}

	var sess mongo.Session
	if sess, err = repository.MongoClient.StartSession(); err != nil {
		log.Error(err.Error())
		return genericFailureMsg
	}
	defer sess.EndSession(ctx)

	if _, err = sess.WithTransaction(ctx, cb); err != nil {
		log.Error(err.Error())
		return genericFailureMsg
	}

	return
}
