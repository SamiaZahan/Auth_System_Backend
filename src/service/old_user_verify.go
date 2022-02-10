package service

import (
	"fmt"
	"github.com/emamulandalib/airbringr-auth/config"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"time"
)

func (a *Auth) OldUserVerify(email string) (err error) {
	genericSignupFailureMsg := errors.New("Login failed for some technical reason.")
	var otp string
	otpSvc := OtpSvc{MicroAPIToken: config.Params.MicroAPIToken}
	if otp, err = otpSvc.Generate(OtpGenerateRequest{
		Expiry: int64(time.Hour * 24),
		Id:     email,
	}); err != nil {
		err = genericSignupFailureMsg
		return
	}
	err = a.OldUserVerifySendEmail(email, otp)
	return
}
func (a *Auth) OldUserVerifySendEmail(email string, otp string) error {
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
			"subject":       "AirBringr Old Shopper Verification",
			"template_code": "signup_verification",
		}).
		String(); code != fiber.StatusOK {
		log.Error(errs)
		return errors.New("Email send failed")
	}
	return nil
}
