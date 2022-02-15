package service

import (
	"context"
	"fmt"
	"github.com/emamulandalib/airbringr-auth/repository"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"

	log "github.com/sirupsen/logrus"

	"github.com/emamulandalib/airbringr-auth/config"
	"github.com/emamulandalib/airbringr-auth/dto"
	"github.com/pkg/errors"

	"time"
)

type PassReset struct{}

func (p *PassReset) PasswordResetOtp(input *dto.EmailOtpInput) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	genericPassResetFailureMsg := errors.New("Password Reset failed for some technical reason.")

	aRepo := repository.Auth{Ctx: ctx}
	_, err = aRepo.GetUserByEmail(input.Email)

	if err != nil && mongo.ErrNoDocuments != err {
		log.Error(err)
		return
	}

	if err == mongo.ErrNoDocuments {
		emailExists := ExistingEmail(input.Email)

		if emailExists.Error {
			return errors.New("something went wrong. please try again later")
		}

		if !emailExists.Status {
			return errors.New("not an user. please signup")
		}
	}

	var otp string
	otpSvc := OtpSvc{MicroAPIToken: config.Params.MicroAPIToken}
	if otp, err = otpSvc.Generate(OtpGenerateRequest{
		Expiry: int64(time.Hour * 24),
		Id:     input.Email,
	}); err != nil {
		return genericPassResetFailureMsg
	}
	err = p.PassResetSendEmail(input.Email, otp)
	return
}
func (p *PassReset) PassResetSendEmail(email string, otp string) error {
	emailSvcURI := fmt.Sprintf("%s/v1/send-email", config.Params.NotificationSvcDomain)
	passwordResetLink := fmt.Sprintf("%s/reset/?otp=%s&auth=%s", config.Params.ServiceFrontend, otp, email)
	if code, _, errs := fiber.
		Post(emailSvcURI).
		JSON(fiber.Map{
			"data": fiber.Map{
				"link": passwordResetLink,
			},
			"to":            email,
			"from":          "contact@airbringr.com",
			"message":       "Please click the link to reset password.",
			"subject":       "AirBringr Password Reset",
			"template_code": "password_reset",
		}).
		String(); code != fiber.StatusOK {
		log.Error(errs)
		return errors.New("Email send failed")
	}
	return nil
}

func (p *PassReset) UpdatePassword(input dto.PasswordReset) (err error) {
	genericFailureMsg := errors.New("password reset failed")
	passwordService := PasswordService{}
	hashedPassword := passwordService.HashPassword(input.Password)
	cb := func(sessCtx mongo.SessionContext) (i interface{}, err error) {
		aRepo := repository.Auth{Ctx: sessCtx}
		if err = aRepo.SetUserPasswordByEmail(input.Auth, hashedPassword); err != nil {
			log.Error(err)
		}

		if code, body, errs := fiber.
			Post(fmt.Sprintf("%s/helper/update-password", config.Params.AirBringrDomain)).
			JSON(fiber.Map{
				"email":    input.Auth,
				"password": input.Password,
			}).
			String(); code != fiber.StatusOK {
			log.Error(body)
			log.Error(errs)
			return nil, genericFailureMsg
		}
		err = p.UpdatePasswordConfirmEmail(input.Auth)
		return
	}

	ctx := context.Background()
	var sess mongo.Session
	if sess, err = repository.MongoClient.StartSession(); err != nil {
		return genericFailureMsg
	}
	defer sess.EndSession(ctx)
	if _, err = sess.WithTransaction(ctx, cb); err != nil {
		log.Error(err.Error())
		return genericFailureMsg
	}
	return
}

func (p *PassReset) UpdatePasswordConfirmEmail(email string) error {
	emailSvcURI := fmt.Sprintf("%s/v1/send-email", config.Params.NotificationSvcDomain)
	if code, _, errs := fiber.
		Post(emailSvcURI).
		JSON(fiber.Map{
			"data":          fiber.Map{},
			"to":            email,
			"from":          "contact@airbringr.com",
			"message":       "Your password has been updated successfully.",
			"subject":       "AirBringr Password Reset Confirmation",
			"template_code": "password_reset_confirmation",
		}).
		String(); code != fiber.StatusOK {
		log.Error(errs)
		return errors.New("Email send failed")
	}
	return nil
}
