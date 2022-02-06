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
	genericPassResetFailureMsg := errors.New("Password Reset failed for some technical reason.")
	var otp string
	otpSvc := OtpSvc{MicroAPIToken: config.Params.MicroAPIToken}
	if otp, err = otpSvc.Generate(OtpGenerateRequest{
		Expiry: int64(time.Hour * 24),
		Id:     input.Email,
	}); err != nil {
		fmt.Print(otp)
		return genericPassResetFailureMsg
	}
	//if err = p.PassResetSendEmail(input.Email, otp); err != nil {
	//	return
	//}
	return
}
func (p *PassReset) PassResetSendEmail(email string, otp string) error {
	emailSvcURI := fmt.Sprintf("%s/v1/send-email", config.Params.NotificationSvcDomain)
	passwordResetLink := fmt.Sprintf("%s/password-reset/?otp=%s&auth=%s", config.Params.ServiceFrontend, otp, email)
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
			return
		}

		// register user into legacy system
		var userDoc *repository.UserDoc
		var userProfileDoc *repository.UserProfileDoc
		if userDoc, err = aRepo.GetUserByEmail(input.Auth); err != nil {
			return
		}
		if userProfileDoc, err = aRepo.GetUserProfileByID(userDoc.ID.Hex()); err != nil {
			return
		}
		if code, body, errs := fiber.
			Post(fmt.Sprintf("%s/helper/register", config.Params.AirBringrDomain)).
			JSON(fiber.Map{
				"name":     fmt.Sprintf("%s %s", userProfileDoc.FirstName, userProfileDoc.LastName),
				"email":    userDoc.Email,
				"phone":    userDoc.Mobile,
				"password": "Vi$FV/kBi<VuZCW2Y9JT_G(NbUj~rV",
			}).
			String(); code != fiber.StatusOK {
			log.Error(body)
			log.Error(errs)
			return nil, genericFailureMsg
		}
		return
	}

	ctx := context.Background()
	var sess mongo.Session
	if sess, err = repository.MongoClient.StartSession(); err != nil {
		return genericFailureMsg
	}
	defer sess.EndSession(ctx)
	if _, err = sess.WithTransaction(ctx, cb); err != nil {
		return genericFailureMsg
	}
	return
}
