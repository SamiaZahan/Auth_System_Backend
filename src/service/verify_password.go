package service

import (
	"context"
	"errors"
	"github.com/emamulandalib/airbringr-auth/repository"
)

type VerifyUserPassword struct{}

func (v VerifyUserPassword) VerifyPassword(password string, email string) (err error) {
	//ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	//defer cancel()
	ctx := context.Background()
	//genericEditFailureMsg := errors.New("Profile Edit failed for some technical reason.")
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
