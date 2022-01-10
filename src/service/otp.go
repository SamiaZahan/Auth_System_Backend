package service

import (
	"errors"
	"go.m3o.com/otp"
	//"github.com/micro/services/clients/go/otp"
	log "github.com/sirupsen/logrus"
)

type OtpSvc struct {
	MicroAPIToken string
}

type OtpGenerateRequest otp.GenerateRequest
type OtpValidateRequest otp.ValidateRequest

func (o *OtpSvc) Generate(genReq OtpGenerateRequest) (code string, err error) {
	svc := otp.NewOtpService(o.MicroAPIToken)
	req := (otp.GenerateRequest)(genReq)
	var resp *otp.GenerateResponse
	if resp, err = svc.Generate(&req); err != nil {
		log.Error(err.Error())
		return
	}
	code = resp.Code
	return
}

func (o *OtpSvc) Validate(validatedReq OtpValidateRequest) bool {
	svc := otp.NewOtpService(o.MicroAPIToken)
	req := (otp.ValidateRequest)(validatedReq)
	var resp *otp.ValidateResponse
	resp, err := svc.Validate(&req)
	if err != nil {
		log.Error(err.Error())
		return false
	}
	if !resp.Success {
		log.Error(errors.New("OTP verification not success from M30"))
		return false
	}
	return true
}
