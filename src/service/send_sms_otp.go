package service

import (
	"errors"
	"fmt"
	"github.com/emamulandalib/airbringr-auth/config"
	"github.com/emamulandalib/airbringr-auth/dto"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
)

type SmsOtp struct{}

func (s *SmsOtp) Send(input dto.SendSmsOtpInput) (err error) {
	otp := GenerateRandNum()
	smsSvcURI := fmt.Sprintf("%s/v1/send-sms", config.Params.NotificationSvcDomain)
	code, _, errs := fiber.
		Post(smsSvcURI).
		JSON(fiber.Map{
			"message": fmt.Sprintf("AirBringr OTP: %d", otp),
			"number":  input.Mobile,
		}).
		String()

	if code != fiber.StatusOK {
		log.Error(errs)
		return errors.New("OTP send failed.")
	}
	return
}
