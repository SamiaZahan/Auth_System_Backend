package service

import (
	"errors"
	"fmt"
	"github.com/emamulandalib/airbringr-auth/config"
	"github.com/emamulandalib/airbringr-auth/dto"
	"github.com/emamulandalib/airbringr-auth/repository"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
	"go.m3o.com/otp"
	"go.mongodb.org/mongo-driver/mongo"
)

func (a *Auth) VerifyAndEditEmail(input dto.EmailVerificationInput, email string, c *fiber.Ctx) (err error) {
	genericErrMsg := errors.New("something went wrong with the verification. Please try again later")
	ctx := c.Context()
	otpSvc := otp.NewOtpService(config.Params.MicroAPIToken)
	resp, err := otpSvc.Validate(&otp.ValidateRequest{
		Code: fmt.Sprintf("%d", input.OTP),
		Id:   input.Auth,
	})
	if err != nil {
		log.Error(err.Error())
		return genericErrMsg
	}
	if !resp.Success {
		log.Error(errors.New("OTP verification not success from M30"))
		return genericErrMsg
	}

	cb := func(sessCtx mongo.SessionContext) (i interface{}, err error) {
		aRepo := repository.Auth{Ctx: sessCtx}
		if err = aRepo.SetUserEmailByEmail(input.Auth, email); err != nil {
			log.Error(err.Error())
			return genericErrMsg, nil
		}
		if code, _, _ := fiber.
			Post(fmt.Sprintf("%s/helper/edit-email", config.Params.AirBringrDomain)).
			JSON(fiber.Map{
				"old_email": email,
				"new_email": input.Auth,
			}).
			String(); code != fiber.StatusOK {
			return nil, genericErrMsg
		}
		return
	}
	var sess mongo.Session
	if sess, err = repository.MongoClient.StartSession(); err != nil {
		return genericErrMsg
	}
	defer sess.EndSession(ctx)
	if _, err = sess.WithTransaction(ctx, cb); err != nil {
		log.Error(err.Error())
		return genericErrMsg
	}
	return
}
