package service

import (
	"encoding/json"
	"fmt"
	"github.com/emamulandalib/airbringr-auth/config"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
)

type DoesUserExists struct{}

type ResponseData struct {
	UserExists bool   `json:"user_exists"`
	PassValid  bool   `json:"pass_valid"`
	Message    string `json:"message"`
	User       []User `json:"user_profile"`
}
type User struct {
	UserId int    `json:"id"`
	Name   string `json:"name"`
	Phone  string `json:"phone"`
	Email  string `json:"email"`
}

func (d *DoesUserExists) DoesUserExists(emailOrMobile string, password string) (resData ResponseData) {
	doesUserExistsURI := fmt.Sprintf("%s/helper/does-user-exist", config.Params.AirBringrDomain)
	statusCode, body, errs := fiber.
		Post(doesUserExistsURI).
		JSON(fiber.Map{
			"email_or_mobile": emailOrMobile,
			"password":        password,
		}).String()
	if statusCode != fiber.StatusOK {
		log.Error(errs)
		return
	}
	_ = json.Unmarshal([]byte(body), &resData)
	return resData
}
