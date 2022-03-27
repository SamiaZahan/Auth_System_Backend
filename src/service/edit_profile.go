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

func (ep EditProfile) EditUserProfile(input *dto.EditProfileInput, email string, c *fiber.Ctx) (err error) {
	ctx, cancel := context.WithTimeout(c.Context(), time.Second*2)
	defer cancel()
	genericEditFailureMsg := errors.New("profile Edit failed for some technical reason")
	aRepo := repository.Auth{Ctx: ctx}
	user, err := aRepo.GetUserByEmail(email)
	if err != nil {
		return errors.New("user not found")
	}
	cb := func(sessCtx mongo.SessionContext) (i interface{}, err error) {
		aRepo := repository.Auth{Ctx: sessCtx}
		if err = aRepo.SetUserProfileByID(user.ID.Hex(), input); err != nil {
			log.Error(err)
			return err, nil
		}

		if code, body, errs := fiber.
			Post(fmt.Sprintf("%s/helper/edit-profile", config.Params.AirBringrDomain)).
			JSON(fiber.Map{
				"email": email,
				"data":  input,
			}).
			String(); code != fiber.StatusOK {
			log.Error(body)
			log.Error(errs)
			return nil, genericEditFailureMsg
		}
		return
	}

	var sess mongo.Session
	if sess, err = repository.MongoClient.StartSession(); err != nil {
		return genericEditFailureMsg
	}
	defer sess.EndSession(ctx)
	if _, err = sess.WithTransaction(ctx, cb); err != nil {
		log.Error(err.Error())
		return genericEditFailureMsg
	}
	return nil
}
func (ep EditProfile) EditEmailOtp(input *dto.EditEmailInput, c *fiber.Ctx) (err error) {
	ctx, cancel := context.WithTimeout(c.Context(), time.Second*2)
	defer cancel()
	genericEmailEditFailureMsg := errors.New("email update failed for some technical reason")
	aRepo := repository.Auth{Ctx: ctx}
	user, _ := aRepo.GetUserByEmail(input.Email)
	if user != nil {
		return errors.New("an user with this email already exists")
	}
	userExistResponse := ExistingEmail(input.Email)
	if userExistResponse.Status {
		return errors.New("an user with this email already exists")
	}
	var otp string
	otpSvc := OtpSvc{MicroAPIToken: config.Params.MicroAPIToken}
	if otp, err = otpSvc.Generate(OtpGenerateRequest{
		Expiry: int64(time.Hour * 6),
		Id:     input.Email,
	}); err != nil {
		return genericEmailEditFailureMsg
	}
	err = ep.EditSendEmail(input.Email, otp)
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
		//return errors.New("email send failed")
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
