package service

import (
	"context"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"time"

	"github.com/emamulandalib/airbringr-auth/config"
	"github.com/emamulandalib/airbringr-auth/dto"
	"github.com/emamulandalib/airbringr-auth/repository"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type EmailOtp struct{}

func (e *EmailOtp) Send(input dto.EmailOtpInput) (err error) {
	genericFailureMsg := errors.New("OTP send failed")
	userNotExistErrMsg := errors.New("no user found with this email. Please signup")
	ctx := context.Background()
	aRepo := repository.Auth{Ctx: ctx}

	if _, err = aRepo.GetUserByEmail(input.Email); err != nil {
		if err == mongo.ErrNoDocuments {
			return userNotExistErrMsg
		}
		return genericFailureMsg
	}

	if !ExistingEmail(input.Email) {
		return userNotExistErrMsg
	}

	otpSvc := OtpSvc{MicroAPIToken: config.Params.MicroAPIToken}
	otp, err := otpSvc.Generate(OtpGenerateRequest{
		Expiry: int64(time.Minute * 5),
		Id:     input.Email,
	})

	if err != nil {
		log.Error(err.Error())
		return genericFailureMsg
	}

	emailSvcURI := fmt.Sprintf("%s/v1/send-email", config.Params.NotificationSvcDomain)

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
		return errors.New("failed to send SMS")
	}

	return
}
