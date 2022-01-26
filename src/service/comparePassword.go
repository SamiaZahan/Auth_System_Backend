package service

import (
	"golang.org/x/crypto/bcrypt"
)

type ComparePassword struct{}

func (c *ComparePassword) ComparePasswords(hashedPwd string, plainPwd []byte) bool {
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		return false
	}

	return true
}
