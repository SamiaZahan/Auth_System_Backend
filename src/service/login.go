package service

import (
	"context"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/emamulandalib/airbringr-auth/config"
	"github.com/gofiber/fiber/v2"
	"time"

	//"fmt"
	//"time"
	//"github.com/emamulandalib/airbringr-auth/config"
	"github.com/emamulandalib/airbringr-auth/dto"
	"github.com/emamulandalib/airbringr-auth/repository"
	//"github.com/gofiber/fiber/v2"
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
	//ctx := context.Background()
	//context with time out (study n input)
	var ctx, cancel = context.WithTimeout(context.Background(), time.Millisecond*000)
	defer cancel()
	authRepo := repository.Auth{Ctx: ctx}

	// try to get existing user
	existingUser, err := authRepo.GetUserByEmailOrMobile(input.EmailOrMobile)
	if err != nil && err != mongo.ErrNoDocuments {
		log.Error(res.Error.Error())
		res.Error = genericLoginFailureMsg
		return res
	}
	hashedPassword, passwordHashingError := authRepo.HashPassword(input.Password)
	if passwordHashingError != nil {
		log.Error(res.Error.Error())
		return
	}
	passwordMatched := authRepo.ComparePasswords(hashedPassword, []byte(existingUser.Password))

	if existingUser != nil && passwordMatched == true {
		code, inputMarshalError := json.Marshal(input)

		return LoginResponse{
			Redirect: true,
			Code:     b64.StdEncoding.EncodeToString([]byte(code)),
			Error:    inputMarshalError,
		}

	}

	//Lookup in Old DB
	doesUserExists := fmt.Sprintf("%s/does-user-exists/?code=%s", config.Params.AirBringrDomain, res.Code)
	statusCode, body, errs := fiber.
		Post(doesUserExists).
		JSON(fiber.Map{
			"emailOrMobile": input.EmailOrMobile,
			"password":      input.Password,
		}).String()
	if statusCode != fiber.StatusOK {
		log.Error(errs)
		res.Error = errors.New("Login Error")
		return res
	}
	if body.status == true {
		authRepo.CreateUser(body.user.Email, body.user.Password)

		code, inputMarshalError := json.Marshal(input)

		return LoginResponse{
			Redirect: true,
			Code:     b64.StdEncoding.EncodeToString([]byte(code)),
			Error:    inputMarshalError,
		}
	}

	if body.status == false {
		res.Error = errors.New("Not an  User. Please Sign Up")
		return res
	}
	return
}
