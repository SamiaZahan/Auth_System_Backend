package handler

import (
	"errors"
	"fmt"
	"github.com/emamulandalib/airbringr-auth/config"
	"github.com/emamulandalib/airbringr-auth/dto"
	"github.com/emamulandalib/airbringr-auth/response"
	"github.com/emamulandalib/airbringr-auth/service"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
)

func (h *Handler) Login(c *fiber.Ctx) error {
	input := new(dto.LoginInput)

	if err := c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Payload{
			Message: response.BodyParseFailedErrorMsg,
			Errors:  errors.New(response.BodyParseFailedErrorMsg),
		})
	}

	err := input.Validate()

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Payload{
			Message: response.ValidationFailedMsg,
			Errors:  err,
		})
	}

	svc := service.Auth{}
	res := svc.Login(*input)
	log.Print(res.Code)
	if res.Redirect {
		err = c.Redirect(fmt.Sprintf("%s/forced-login/?code=%s", config.Params.AirBringrDomain, res.Code))
		if err != nil {
			log.Error(err.Error())
			return c.Status(fiber.StatusBadRequest).JSON(response.Payload{
				Message: "Failed to Login",
				Errors:  err,
			})
		}
	}

	if res.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Payload{
			Message: res.Error.Error(),
			Errors:  err,
		})
	}

	return nil
}
