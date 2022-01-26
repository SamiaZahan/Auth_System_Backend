package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/emamulandalib/airbringr-auth/config"
	"github.com/emamulandalib/airbringr-auth/dto"
	"github.com/emamulandalib/airbringr-auth/repository"
	log "github.com/sirupsen/logrus"
	"go.m3o.com/otp"
)

func (a *Auth) EmailVerification(input dto.EmailVerificationInput) (err error) {
	genericErrMsg := errors.New("something went wrong with the verification. Please try again later")
	ctx := context.Background()
	aRepo := repository.Auth{Ctx: ctx}

	otpSvc := otp.NewOtpService(config.Params.MicroAPIToken)
	resp, err := otpSvc.Validate(&otp.ValidateRequest{
		Code: fmt.Sprintf("%d", input.OTP),
		Id:   input.Auth,
	})

	if err != nil {
		return genericErrMsg
	}

	if !resp.Success {
		log.Error(errors.New("OTP verification not success from M30"))
		return genericErrMsg
	}

	if err = aRepo.ActivateUserByEmail(input.Auth); err != nil {
		return genericErrMsg
	}

	return
}
