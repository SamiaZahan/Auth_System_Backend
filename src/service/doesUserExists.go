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
	status  bool   `json:"status"`
	message string `json:"message"`
	user    struct {
		userId    string `json:"user_id"`
		email     string `json:"email"`
		phone     string `json:"phone"`
		password  string `json:"password"`
		firstName string `json:"first_name"`
		lastName  string `json:"last_name"`
	} `json:"user"`
}

func (d *DoesUserExists) DoesUserExists(emailOrMobile string) (resData ResponseData) {
	doesUserExistsURI := fmt.Sprintf("%s/does_user_exists", config.Params.AirBringrDomain)
	statusCode, body, errs := fiber.
		Post(doesUserExistsURI).
		JSON(fiber.Map{
			"emailOrMobile": emailOrMobile,
		}).String()
	if statusCode != fiber.StatusOK {
		log.Error(errs)
		return
	}
	_ = json.Unmarshal([]byte(body), &resData)
	return resData
}
