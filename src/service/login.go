package service

import (
	"context"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/emamulandalib/airbringr-auth/config"
	"github.com/emamulandalib/airbringr-auth/dto"
	"github.com/emamulandalib/airbringr-auth/repository"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"strings"
	"time"
)

type LoginResponse struct {
	Redirect bool
	Code     string // base64 encoded base64.encode({"emailOrPhone": "", "password": ""})
	Error    error
}

func (a *Auth) Login(input dto.LoginInput) (res LoginResponse) {
	genericLoginFailureMsg := errors.New("Login failed for some technical reason.")
	var ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	authRepo := repository.Auth{Ctx: ctx}
	passwordService := PasswordService{}
	// input modify for force login
	modifiedInput := map[string]string{
		"email_or_mobile": input.EmailOrMobile,
		"password":        input.Password,
	}
	code, inputMarshalError := json.Marshal(modifiedInput)

	// try to get existing user
	existingUser, err := authRepo.GetUserByEmailOrMobile(input.EmailOrMobile)
	if err == nil {
		passwordMatched := passwordService.ComparePasswords(existingUser.Password, []byte(input.Password))
		if !existingUser.EmailVerified {
			return LoginResponse{Error: errors.New("Email is not verified")}
		}
		if !existingUser.MobileVerified {
			return LoginResponse{Error: errors.New("Mobile number is not verified")}
		}
		//TODO: future scope: a scheduler will remove the  unverified users within certain time.
		if passwordMatched {
			return LoginResponse{
				Redirect: true,
				Code:     b64.StdEncoding.EncodeToString([]byte(code)),
				Error:    inputMarshalError,
			}
		}
		return LoginResponse{Error: errors.New("Wrong password")}
	}
	if err != nil {
		if err != mongo.ErrNoDocuments {
			log.Error(errors.New("User not found"))
		}
		log.Error(err.Error())
		//Lookup in Old DB
		doesUserExists := DoesUserExists{}
		response := doesUserExists.DoesUserExists(input.EmailOrMobile, input.Password)
		if !response.UserExists {
			return LoginResponse{Error: errors.New("Not an user. Please sign up")}
		}
		if !response.PassValid {
			return LoginResponse{Error: errors.New("Wrong password")}
		}

		//Old Shopper Email, Mobile  Verification

		var otp string
		otpSvc := OtpSvc{MicroAPIToken: config.Params.MicroAPIToken}
		if otp, err = otpSvc.Generate(OtpGenerateRequest{
			Expiry: int64(time.Hour * 24),
			Id:     response.User.Email,
		}); err != nil {
			fmt.Print(otp)
			return LoginResponse{Error: genericLoginFailureMsg}
		}
		createVerificationLink := func(sessCtx mongo.SessionContext) (i interface{}, err error) {
			var userID string
			AuthRpo := repository.Auth{Ctx: sessCtx}
			hashedPassword := passwordService.HashPassword(input.Password)
			if userID, err = AuthRpo.CreateUser(response.User.Email, hashedPassword); err != nil {
				return
			}
			name := response.User.Name
			lastName := name[strings.LastIndex(name, " ")+1:]
			firstName := strings.TrimSuffix(name, lastName)
			if err = AuthRpo.CreateUserProfile(userID, firstName, lastName); err != nil {
				return
			}
			err = a.OldUserVerifySendEmail(response.User.Email, otp)
			return
		}

		var sess mongo.Session
		if sess, err = repository.MongoClient.StartSession(); err != nil {
			log.Error(err.Error())
			return LoginResponse{Error: genericLoginFailureMsg}
		}
		defer sess.EndSession(ctx)

		if _, err = sess.WithTransaction(ctx, createVerificationLink); err != nil {
			log.Error(err.Error())
			return LoginResponse{Error: genericLoginFailureMsg}
		}

		return
	}
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
			"message":       "Dear valued Shopper, We are migrating our system. Please click the link to verify your account.",
			"subject":       "AirBringr Old Shopper Verification",
			"template_code": "signup_verification",
		}).
		String(); code != fiber.StatusOK {
		log.Error(errs)
		return errors.New("Email send failed")
	}
	return nil
}
