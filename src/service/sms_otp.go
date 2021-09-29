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

type SmsOtp struct{}

func (s *SmsOtp) Send(input dto.SendSmsOtpInput) (err error) {
	genericSignupFailureMsg := errors.New("OTP send failed")
	otp := GenerateRandNum()
	smsSvcURI := fmt.Sprintf("%s/v1/send-sms", config.Params.NotificationSvcDomain)

	ctx := context.Background()
	var sess mongo.Session
	if sess, err = repository.MongoClient.StartSession(); err != nil {
		log.Error(err.Error())
		return genericSignupFailureMsg
	}
	defer sess.EndSession(ctx)

	tryToSend := func(sessCtx mongo.SessionContext) (i interface{}, err error) {
		vRepo := repository.Verification{Ctx: sessCtx}
		if _, err = vRepo.Create(input.Mobile, otp, ""); err != nil {
			return
		}

		if code, _, errs := fiber.
			Post(smsSvcURI).
			JSON(fiber.Map{
				"message": fmt.Sprintf("AirBringr %d", otp),
				"number":  input.Mobile,
			}).
			String(); code != fiber.StatusOK {
			log.Error(errs)
			return nil, errors.New("falied to send SMS")
		}

		return
	}

	if _, err = sess.WithTransaction(ctx, tryToSend); err != nil {
		log.Error(err.Error())
		return genericSignupFailureMsg
	}

	return
}
