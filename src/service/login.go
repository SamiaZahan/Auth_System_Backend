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
	var ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	authRepo := repository.Auth{Ctx: ctx}

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

	type response struct {
		status  bool
		message string
		user    struct {
			email    string
			password string
		}
	}
	data := response{}
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
	_, err = authRepo.CreateUser(data.user.email, hashedPassword)
	if err != nil {
		return LoginResponse{}
	}

	code, inputMarshalError := json.Marshal(input)

	return LoginResponse{
		Redirect: true,
		Code:     b64.StdEncoding.EncodeToString([]byte(code)),
		Error:    inputMarshalError,
	}

	return
}
