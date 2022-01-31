package service

import (
	"encoding/json"
	"fmt"
	"github.com/emamulandalib/airbringr-auth/config"
	"github.com/gofiber/fiber/v2"
)

type Auth struct{}

type ExistUserResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Error   bool   `json:"error"`
}

func ExistingEmail(email string) (resData ExistUserResponse) {
	_, body, _ := fiber.
		Post(fmt.Sprintf("%s/helper/exist-email", config.Params.AirBringrDomain)).
		JSON(fiber.Map{
			"email": email,
		}).
		String()

	_ = json.Unmarshal([]byte(body), &resData)
	return resData
	//if code != fiber.StatusOK {
	//	log.Error(errs)
	//	exists = true
	//	return
	//}
	//if errs != nil {
	//	log.Error(errs)
	//}
	//exists = false
	//return
}

func ExistingMobile(phone string) (resData ExistUserResponse) {
	_, body, _ := fiber.
		Post(fmt.Sprintf("%s/helper/exist-phone", config.Params.AirBringrDomain)).
		JSON(fiber.Map{
			"phone": phone,
		}).
		String()
	_ = json.Unmarshal([]byte(body), &resData)
	return resData
	//if code != fiber.StatusOK {
	//	log.Error(errs)
	//	exists = true
	//	return
	//}
	//if errs != nil {
	//	log.Error(errs)
	//}
	//exists = false
	//return
}
