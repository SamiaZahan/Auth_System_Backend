package service

import (
	"errors"
	"github.com/emamulandalib/airbringr-auth/repository"
	"github.com/gofiber/fiber/v2"
)

type VerifyUserPassword struct{}

func (v VerifyUserPassword) VerifyPassword(password string, email string, c *fiber.Ctx) (err error) {
	ctx := c.Context()
	aRepo := repository.Auth{Ctx: ctx}
	passwordService := PasswordService{}
	user, err := aRepo.GetUserByEmail(email)
	if err != nil {
		return errors.New("user not found")
	}
	matched := passwordService.ComparePasswords(user.Password, []byte(password))
	if !matched {
		return errors.New("password didn't match")
	}
	return nil
}
