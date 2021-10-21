package service

import (
	"github.com/micro/services/clients/go/otp"
	log "github.com/sirupsen/logrus"
)

type OtpSvc struct {
	MicroAPIToken string
}

type OtpGenerateRequest otp.GenerateRequest

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
