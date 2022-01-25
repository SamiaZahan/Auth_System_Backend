package service

import (
	"fmt"
	"math/rand"

	"github.com/emamulandalib/airbringr-auth/config"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
)

type Auth struct{}

func GenerateRandNum() int {
	min := 1000
	max := 9999
	return rand.Intn(max-min) + min
}

func ExistingEmail(email string) (exists bool) {
	if code, _, errs := fiber.
		Post(fmt.Sprintf("%s/helper/exist-email", config.Params.AirBringrDomain)).
		JSON(fiber.Map{
			"email": email,
		}).
		String(); code != fiber.StatusOK {
		log.Error(errs)
		exists = false
		return
	}

	exists = true
	return
}

func ExisitingMobile(phone string) (exists bool) {
	if code, _, errs := fiber.
		Post(fmt.Sprintf("%s/helper/exist-phone", config.Params.AirBringrDomain)).
		JSON(fiber.Map{
			"phone": phone,
		}).
		String(); code != fiber.StatusOK {
		log.Error(errs)
		exists = true
		return
	}

	exists = false
	return
}
