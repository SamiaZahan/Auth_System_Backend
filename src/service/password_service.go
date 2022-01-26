package service

import (
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type PasswordService struct{}

func (c *PasswordService) ComparePasswords(hashedPwd string, plainPwd []byte) bool {
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		log.Error(err)
		return false
	}
	return true
}

func (c *PasswordService) HashPassword(password string) (hashedPassword string) {
	hashed, passwordHashingError := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if passwordHashingError != nil {
		log.Error(passwordHashingError.Error())
		return
	}
	hashedPassword = string(hashed)
	return
}
