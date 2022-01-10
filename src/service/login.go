package service

import (
	"context"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/emamulandalib/airbringr-auth/config"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
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

	//input.EmailOrMobile check is email??
	err := validation.Validate(input.EmailOrMobile, is.Email)
	if err != nil {
		if input.CountryPrefix == "" {
			return LoginResponse{Error: errors.New("Country Prefix Required")}
		}
		phoneNumberMap, err := authRepo.GetInfoByCountryPrefix(input.CountryPrefix)
		if err != nil {
			return LoginResponse{Error: errors.New("Not a valid Country Prefix")}
		}
		countryCode := phoneNumberMap.CountryCode
		//check valid phone number
		phoneValidate := PhoneValidate{}
		valid, _ := phoneValidate.ValidatePhoneNumber(input.EmailOrMobile, countryCode)
		if !valid {
			return LoginResponse{Error: errors.New("Not a Valid Phone Number")}
		}

	}
	// try to get existing user
	existingUser, err := authRepo.GetUserByEmailOrMobile(input.EmailOrMobile)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			return LoginResponse{Error: errors.New("User not found")}
		}
		log.Error(err.Error())
		return LoginResponse{Error: genericLoginFailureMsg}
	}

	passwordMatched := authRepo.ComparePasswords(existingUser.Password, []byte(input.Password))

	if passwordMatched {
		code, inputMarshalError := json.Marshal(input)

		return LoginResponse{
			Redirect: true,
			Code:     b64.StdEncoding.EncodeToString([]byte(code)),
			Error:    inputMarshalError,
		}

	}

	//Lookup in Old DB
	doesUserExistsURI := fmt.Sprintf("%s/does-user-exists/?code=%s", config.Params.AirBringrDomain, res.Code)
	statusCode, body, errs := fiber.
		Post(doesUserExistsURI).
		JSON(fiber.Map{
			"emailOrMobile": input.EmailOrMobile,
			"password":      input.Password,
		}).String()
	if statusCode != fiber.StatusOK {
		log.Error(errs)
		res.Error = errors.New("Failed to login")
		return res
	}

	type Response struct {
		status  bool   `json:"status"`
		message string `json:"message"`
		user    struct {
			userId    string `json:"user_id"`
			email     string `json:"email"`
			firstName string `json:"first_name"`
			lastName  string `json:"last_name"`
		} `json:"user"`
	}
	var data Response
	_ = json.Unmarshal([]byte(body), &data)

	if !data.status {
		res.Error = errors.New("Not an  User. Please Sign Up")
		return res
	}

	hashedPassword, passwordHasingError := authRepo.HashPassword(input.Password)
	if err != nil {
		log.Error(passwordHasingError.Error())
		return
	}

	//Transaction
	wc := writeconcern.New(writeconcern.WMajority())
	rc := readconcern.Snapshot()
	txnOpts := options.Transaction().SetWriteConcern(wc).SetReadConcern(rc)
	if err != nil {
		panic(err)
	}
	session, err := repository.MongoClient.StartSession()
	if err != nil {
		log.Fatal(err)
	}
	defer session.EndSession(context.Background())
	err = mongo.WithSession(context.Background(), session, func(sessionContext mongo.SessionContext) error {
		if err = session.StartTransaction(txnOpts); err != nil {
			return err
		}
		AuthRepo := repository.Auth{Ctx: sessionContext}
		_, err = AuthRepo.CreateUser(data.user.email, hashedPassword)
		if err != nil {
			return err
		}
		err = AuthRepo.CreateUserProfile(data.user.userId, data.user.firstName, data.user.lastName)
		if err != nil {
			return err
		}
		if err = session.CommitTransaction(sessionContext); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		if abortErr := session.AbortTransaction(context.Background()); abortErr != nil {
			panic(abortErr)
		}
		panic(err)
	}

	code, inputMarshalError := json.Marshal(input)

	return LoginResponse{
		Redirect: true,
		Code:     b64.StdEncoding.EncodeToString([]byte(code)),
		Error:    inputMarshalError,
	}

	return
}
