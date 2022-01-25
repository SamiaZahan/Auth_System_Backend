package service

import (
	"context"
	b64 "encoding/base64"
	"encoding/json"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"

	"github.com/emamulandalib/airbringr-auth/dto"
	"github.com/emamulandalib/airbringr-auth/repository"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
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
	comparePassword := ComparePassword{}

	// input modify for force login
	modifiedInput := map[string]string{
		"email_or_mobile": input.EmailOrMobile,
		"password":        input.Password,
	}
	code, inputMarshalError := json.Marshal(modifiedInput)

	//input.EmailOrMobile check is email??
	err := validation.Validate(input.EmailOrMobile, is.Email)
	if err != nil {
		phoneNumberMap, err := authRepo.GetInfoByCountryPrefix(input.CountryPrefix)
		if err != nil {
			return LoginResponse{Error: errors.New("Not a valid Country Prefix")}
		}
		countryCode := phoneNumberMap.CountryCode
		//check valid phone number
		phoneValidate := PhoneValidate{}
		valid, _ := phoneValidate.Validate(input.EmailOrMobile, countryCode)
		if !valid {
			return LoginResponse{Error: errors.New("Not a Valid Phone Number")}
		}

	}
	// try to get existing user
	existingUser, err := authRepo.GetUserByEmailOrMobile(input.EmailOrMobile)

	if err == nil {
		passwordMatched := comparePassword.ComparePasswords(existingUser.Password, []byte(input.Password))

		if passwordMatched {
			return LoginResponse{
				Redirect: true,
				Code:     b64.StdEncoding.EncodeToString([]byte(code)),
				Error:    inputMarshalError,
			}
		}
		if !passwordMatched {
			return LoginResponse{Error: errors.New("Wrong Password")}
		}

	}
	if err != nil {
		if err != mongo.ErrNoDocuments {
			return LoginResponse{Error: errors.New("User not found")}
		}
		log.Error(err.Error())

		//Lookup in Old DB
		doesUserExists := DoesUserExists{}
		response := doesUserExists.DoesUserExists(input.EmailOrMobile, input.Password)

		if !response.UserExists {
			return LoginResponse{Error: errors.New("Not an  User. Please Sign Up")}
		}
		if !response.PassValid {
			return LoginResponse{Error: errors.New("Wrong Password")}
		}

		hashedPassword, passwordHasingError := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)
		if passwordHasingError != nil {
			log.Error(passwordHasingError.Error())
			return
		}

		//Transaction
		wc := writeconcern.New(writeconcern.WMajority())
		rc := readconcern.Snapshot()
		txnOpts := options.Transaction().SetWriteConcern(wc).SetReadConcern(rc)
		//if err != nil {
		//	panic(err)
		//}
		insertUser := func(sessionContext mongo.SessionContext) (i interface{}, err error) {
			AuthRepo := repository.Auth{Ctx: sessionContext}
			var number dto.SendSmsOtpInput
			_, err = AuthRepo.CreateUser(response.User[0].Email, string(hashedPassword), number.Mobile)
			if err != nil {
				return
			}

			//splitting username
			name := response.User[0].Name
			lastName := name[strings.LastIndex(name, " ")+1:]
			firstName := strings.TrimSuffix(name, lastName)

			err = AuthRepo.CreateUserProfile(string(response.User[0].UserId), firstName, lastName)
			if err != nil {
				return
			}
			if err != nil {
				return
			}
			return
		}

		var session mongo.Session
		if session, err = repository.MongoClient.StartSession(); err != nil {
			log.Error(err.Error())
			return LoginResponse{Error: genericLoginFailureMsg}
		}
		defer session.EndSession(context.Background())

		if _, err = session.WithTransaction(context.Background(), insertUser, txnOpts); err != nil {
			log.Error(err.Error())
			return LoginResponse{Error: genericLoginFailureMsg}
		}

		return LoginResponse{
			Redirect: true,
			Code:     b64.StdEncoding.EncodeToString([]byte(code)),
			Error:    inputMarshalError,
		}
		return
	}
	return
}
