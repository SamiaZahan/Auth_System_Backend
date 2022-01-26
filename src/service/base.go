package service

import (
	"fmt"
	"github.com/emamulandalib/airbringr-auth/config"
	"github.com/gofiber/fiber/v2"
)

type Auth struct{}

func ExistingEmail(email string) (exists bool) {
	if code, _, _ := fiber.
		Post(fmt.Sprintf("%s/helper/exist-email", config.Params.AirBringrDomain)).
		JSON(fiber.Map{
			"email": email,
		}).
		String(); code != fiber.StatusOK {
		exists = true
		return
	}

	exists = false
	return
}

func ExistingMobile(phone string) (exists bool) {
	if code, _, _ := fiber.
		Post(fmt.Sprintf("%s/helper/exist-phone", config.Params.AirBringrDomain)).
		JSON(fiber.Map{
			"phone": phone,
		}).
		String(); code != fiber.StatusOK {
		exists = true
		return
	}

	exists = false
	return
}
