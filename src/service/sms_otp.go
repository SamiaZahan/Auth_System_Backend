package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/emamulandalib/airbringr-auth/config"
	"github.com/emamulandalib/airbringr-auth/dto"
	"github.com/emamulandalib/airbringr-auth/repository"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

type SmsOtp struct{}

func (s *SmsOtp) Send(mobile string) (err error) {
	genericFailureMsg := errors.New("OTP send failed")
	otpSvc := OtpSvc{MicroAPIToken: config.Params.MicroAPIToken}
	otp, err := otpSvc.Generate(OtpGenerateRequest{
		Expiry: int64(time.Minute * 5),
		Id:     mobile,
	})
	if err != nil {
		return genericFailureMsg
	}

	smsSvcURI := fmt.Sprintf("%s/v1/send-sms", config.Params.NotificationSvcDomain)
	if code, _, errs := fiber.
		Post(smsSvcURI).
		JSON(fiber.Map{
			"message": fmt.Sprintf("AirBringr: %s", otp),
			"number":  mobile,
		}).
		String(); code != fiber.StatusOK {
		log.Error(errs)
		//return errors.New("failed to send SMS")
	}
	return
}

func (s *SmsOtp) SendSmsOtp(input dto.SendSmsOtpInput) (err error) {
	genericErrMsg := errors.New("OTP send failed")
	ctx := context.Background()
	authRepo := repository.Auth{Ctx: ctx}
	_, err = authRepo.GetUserByMobile(input.Mobile)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New("this number is not registered")
		}
		return genericErrMsg
	}
	err = s.Send(input.Mobile)
	if err != nil {
		return genericErrMsg
	}
	return
}

//func (s *SmsOtp) EditMobileOtpSend(input dto.SendSmsOtpInput, email string) (err error) {
//	genericFailureMsg := errors.New("OTP send failed")
//	mblNmbrExistMsg := errors.New("mobile number already taken")
//	var u *repository.UserDoc
//	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
//	defer cancel()
//	aRepo := repository.Auth{Ctx: ctx}
//	if u, err = aRepo.GetUserByMobile(input.Mobile); err != nil && err != mongo.ErrNoDocuments {
//		return genericFailureMsg
//	}
//	if u != nil {
//		return mblNmbrExistMsg
//	}
//	err = s.Send(input.Mobile)
//	return
//}

func (s *SmsOtp) MobileVerificationOtp(input dto.SendSmsOtpInput) (err error) {
	genericFailureMsg := errors.New("OTP send failed")
	mblNmbrExistMsg := errors.New("mobile number already taken")
	var u *repository.UserDoc

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	aRepo := repository.Auth{Ctx: ctx}
	if u, err = aRepo.GetUserByMobile(input.Mobile); err != nil && err != mongo.ErrNoDocuments {
		return genericFailureMsg
	}
	if u != nil {
		return mblNmbrExistMsg
	}
	err = s.Send(input.Mobile)
	return
}

func (s *SmsOtp) Verify(input dto.VerifyOtpInput) (err error) {
	genericFailureMsg := errors.New("OTP verification failed")
	otpSvc := OtpSvc{MicroAPIToken: config.Params.MicroAPIToken}
	isValid := otpSvc.Validate(OtpValidateRequest{
		Code: fmt.Sprintf("%d", input.OTP),
		Id:   input.EmailOrMobile,
	})
	if !isValid {
		return genericFailureMsg
	}
	return
}

func (s *SmsOtp) VerifyAndRegisterMobileNumber(input dto.VerifyMobileInput) (err error) {
	genericFailureMsg := errors.New("mobile verification failed")
	otpSvc := OtpSvc{MicroAPIToken: config.Params.MicroAPIToken}
	isValid := otpSvc.Validate(OtpValidateRequest{
		Code: fmt.Sprintf("%d", input.OTP),
		Id:   input.Mobile,
	})
	if !isValid {
		log.Error(errors.New("OTP verification not success from M30"))
		return genericFailureMsg
	}

	cb := func(sessCtx mongo.SessionContext) (i interface{}, err error) {
		aRepo := repository.Auth{Ctx: sessCtx}
		if err = aRepo.SetUserMobileByEmail(input.Auth, input.Mobile); err != nil {
			return
		}

		//register user into legacy system
		var userDoc *repository.UserDoc
		var userProfileDoc *repository.UserProfileDoc
		if userDoc, err = aRepo.GetUserByEmail(input.Auth); err != nil {
			return
		}
		if userProfileDoc, err = aRepo.GetUserProfileByID(userDoc.ID.Hex()); err != nil {
			return
		}

		if userDoc.ExistingUser {
			code, _, errs := fiber.
				Post(fmt.Sprintf("%s/helper/update-phone-number", config.Params.AirBringrDomain)).
				JSON(fiber.Map{
					"name":     fmt.Sprintf("%s %s", userProfileDoc.FirstName, userProfileDoc.LastName),
					"email":    userDoc.Email,
					"phone":    userDoc.Mobile,
					"password": "*qSdn<<rha7eFb6<rPFF.!4=Nk%=6R",
				}).
				String()

			if code != fiber.StatusOK || errs != nil {
				return nil, errors.New("phone number update failed")
			}
			return
		}

		code, _, errs := fiber.
			Post(fmt.Sprintf("%s/helper/register-v2", config.Params.AirBringrDomain)).
			JSON(fiber.Map{
				"name":     fmt.Sprintf("%s %s", userProfileDoc.FirstName, userProfileDoc.LastName),
				"email":    userDoc.Email,
				"phone":    userDoc.Mobile,
				"password": "Vi$FV/kBi<VuZCW2Y9JT_G(NbUj~rV",
			}).
			String()

		if code != fiber.StatusOK || errs != nil {
			return nil, errors.New("user registration failed")
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
		log.Error(err.Error())
		return err
	}
	return
}
