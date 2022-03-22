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
	"time"
)

type EditProfile struct{}

func (ep EditProfile) EditUserProfile(input *dto.EditProfileInput, email string) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	//genericEditFailureMsg := errors.New("Profile Edit failed for some technical reason.")

	aRepo := repository.Auth{Ctx: ctx}
	user, err := aRepo.GetUserByEmail(email)
	if err != nil {
		return errors.New("user not found")
	}
	err = aRepo.SetUserProfileByID(user.ID.Hex(), input)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	if err != nil && mongo.ErrNoDocuments != err {
		log.Error(err)
		return errors.New("user not found")
	}
	return nil

}
func (ep EditProfile) EditEmailOtp(input *dto.EditEmailInput, email string) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	genericEmailEditFailureMsg := errors.New("email update failed for some technical reason")
	aRepo := repository.Auth{Ctx: ctx}
	user, err := aRepo.GetUserByEmail(email)
	fmt.Printf(user.Email)
	if err != nil {
		return errors.New("user not found")
	}
	if err != nil && mongo.ErrNoDocuments != err {
		log.Error(err)
		return errors.New("error while finding user")
	}
	var otp string
	otpSvc := OtpSvc{MicroAPIToken: config.Params.MicroAPIToken}
	if otp, err = otpSvc.Generate(OtpGenerateRequest{
		Expiry: int64(time.Hour * 24),
		Id:     input.Email,
	}); err != nil {
		return genericEmailEditFailureMsg
	}
	err = ep.EditSendEmail(input.Email, otp)
	return
	return nil
}

func (ep EditProfile) EditSendEmail(email string, otp string) error {
	emailSvcURI := fmt.Sprintf("%s/v1/send-email", config.Params.NotificationSvcDomain)
	emailUpdateLink := fmt.Sprintf("%s/verify-and-update-email/?otp=%s&auth=%s", config.Params.ServiceFrontend, otp, email)
	if code, _, errs := fiber.
		Post(emailSvcURI).
		JSON(fiber.Map{
			"data": fiber.Map{
				"link": emailUpdateLink,
			},
			"to":            email,
			"from":          "contact@airbringr.com",
			"message":       "Please click the link to updated your email.",
			"subject":       "AirBringr Email Update Verification",
			"template_code": "signup_verification",
		}).
		String(); code != fiber.StatusOK {
		log.Error(errs)
		return errors.New("email send failed")
	}
	return nil
}

func (ep *EditProfile) EditUserEmail(input *dto.EditEmailInput, email string) (err error) {
	genericFailureMsg := errors.New("email edit failed")
	cb := func(sessCtx mongo.SessionContext) (i interface{}, err error) {
		aRepo := repository.Auth{Ctx: sessCtx}
		if err = aRepo.SetUserEmailByEmail(input.Email, email); err != nil {
			log.Error(err)
		}

		if code, body, errs := fiber.
			Post(fmt.Sprintf("%s/helper/edit-email", config.Params.AirBringrDomain)).
			JSON(fiber.Map{
				"oldEmail": email,
				"newEmail": input.Email,
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
		log.Error(err.Error())
		return genericFailureMsg
	}
	return
}

func (ep EditProfile) EditUserMobile(input *dto.EditMobileInput, email string) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	aRepo := repository.Auth{Ctx: ctx}
	user, err := aRepo.GetUserByEmail(email)
	fmt.Printf(user.Email)
	if err != nil {
		return errors.New("user not found")
	}
	if err != nil && mongo.ErrNoDocuments != err {
		log.Error(err)
		return errors.New("user not found")
	}
	return nil
}
