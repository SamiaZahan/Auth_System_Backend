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
	genericFailureMsg := errors.New("OTP send failed")
	otp := GenerateRandNum()
	smsSvcURI := fmt.Sprintf("%s/v1/send-sms", config.Params.NotificationSvcDomain)

	ctx := context.Background()
	var sess mongo.Session
	if sess, err = repository.MongoClient.StartSession(); err != nil {
		log.Error(err.Error())
		return genericFailureMsg
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
				"message": fmt.Sprintf("AirBringr: %d", otp),
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
		return genericFailureMsg
	}

	return
}

func (s *SmsOtp) Verify(input dto.VerifySmsOtpInput) (err error) {
	ctx := context.Background()
	genericFailureMsg := errors.New("OTP verification failed")
	vRepo := repository.Verification{Ctx: ctx}
	var vDoc repository.VerificationDoc

	if vDoc, err = vRepo.GetByEmailOrMobileAndCode(input.Mobile, input.OTP); err != nil {
		log.Error(err.Error())
		return genericFailureMsg
	}

	_ = vRepo.DeleteByID(vDoc.ID.Hex())
	return
}

func (s *SmsOtp) VerifyAndRegisterMobileNumber(input dto.VerifyMobileInput) (err error) {
	ctx := context.Background()
	genericFailureMsg := errors.New("mobile verification failed")
	vRepo := repository.Verification{Ctx: ctx}
	var vDoc repository.VerificationDoc
	var authVerDoc repository.VerificationDoc

	if vDoc, err = vRepo.GetByEmailOrMobileAndCode(input.Mobile, input.OTP); err != nil {
		log.Error(err.Error())
		return genericFailureMsg
	}

	if authVerDoc, err = vRepo.GetByID(input.Auth); err != nil {
		log.Error(err.Error())
		return genericFailureMsg
	}

	userID := authVerDoc.UserID.Hex()

	cb := func(sessCtx mongo.SessionContext) (i interface{}, err error) {
		VRepo := repository.Verification{Ctx: sessCtx}
		aRepo := repository.Auth{Ctx: sessCtx}

		if err = aRepo.SetUserMobileByID(userID, input.Mobile); err != nil {
			return
		}

		if err = VRepo.DeleteByID(input.Auth); err != nil {
			return
		}

		if err = VRepo.DeleteByID(vDoc.ID.Hex()); err != nil {
			return
		}

		// register user into legacy system
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
