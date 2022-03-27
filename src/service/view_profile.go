package service

import (
	"context"
	"errors"
	"github.com/emamulandalib/airbringr-auth/repository"
	"github.com/gofiber/fiber/v2"
	"time"
)

type ViewProfile struct{}

func (p ViewProfile) ViewUserProfile(email string, c *fiber.Ctx) (err error, res interface{}, address interface{}) {
	var ctx, cancel = context.WithTimeout(c.Context(), 2*time.Second)
	defer cancel()
	authRepo := repository.Auth{Ctx: ctx}
	userDoc, err := authRepo.GetUserByEmail(email)
	if err != nil {
		return errors.New("error in  loading data"), nil, nil
	}
	userProfile, err := authRepo.GetUserProfileByID(userDoc.ID.Hex())
	if err != nil {
		return errors.New("error in  loading data"), nil, nil
	}
	res = map[string]string{
		"image":      userProfile.ProfilePicURI,
		"first_name": userProfile.FirstName,
		"last_name":  userProfile.LastName,
		"gender":     userProfile.Gender,
		"age":        userProfile.Age,
		"email":      userDoc.Email,
		"mobile":     userDoc.Mobile,
	}
	address = map[string]string{
		"division": userProfile.Address.Division,
		"district": userProfile.Address.District,
		"area":     userProfile.Address.Area,
		"text":     userProfile.Address.Text,
		"zone":     userProfile.Address.Zone,
	}
	return
}
