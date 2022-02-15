package service

import (
	"encoding/json"
	"fmt"
	"github.com/emamulandalib/airbringr-auth/config"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
)

type Auth struct{}

type ExistUserResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Error   bool   `json:"error"`
}

func ExistingEmail(email string) (resData ExistUserResponse) {
	_, body, errs := fiber.
		Post(fmt.Sprintf("%s/helper/exist-email", config.Params.AirBringrDomain)).
		JSON(fiber.Map{
			"email": email,
		}).
		String()
	if errs != nil {
		log.Error(errs)
		return ExistUserResponse{
			Status: true,
		}
	}
	err := json.Unmarshal([]byte(body), &resData)
	if err != nil {
		log.Error(err)
		return ExistUserResponse{
			Status: true,
		}
	}
	return resData
}

func ExistingMobile(phone string) (resData ExistUserResponse) {
	_, body, errs := fiber.
		Post(fmt.Sprintf("%s/helper/exist-phone", config.Params.AirBringrDomain)).
		JSON(fiber.Map{
			"phone": phone,
		}).
		String()
	if errs != nil {
		log.Error(errs)
		return ExistUserResponse{
			Status: true,
		}
	}
	err := json.Unmarshal([]byte(body), &resData)
	if err != nil {
		log.Error(err)
		return ExistUserResponse{
			Status: true,
		}
	}
	return resData
}
