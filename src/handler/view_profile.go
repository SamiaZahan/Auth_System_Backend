package handler

import (
	"context"
	"github.com/emamulandalib/airbringr-auth/repository"
	"github.com/emamulandalib/airbringr-auth/response"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

func (h *Handler) ViewProfile(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	email := claims["email"].(string)
	var ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	authRepo := repository.Auth{Ctx: ctx}
	userDoc, _ := authRepo.GetUserByEmail(email)
	userProfile, _ := authRepo.GetUserProfileByID(userDoc.ID.Hex())
	res := map[string]string{
		"image":      userProfile.ProfilePicURI,
		"first_name": userProfile.FirstName,
		"last_name":  userProfile.LastName,
		"gender":     userProfile.Gender,
		"age":        userProfile.Age,
		"email":      userDoc.Email,
		"mobile":     userDoc.Mobile,
	}
	address := map[string]string{
		"division": userProfile.Address.Division,
		"district": userProfile.Address.District,
		"area":     userProfile.Address.Area,
		"text":     userProfile.Address.Text,
		"zone":     userProfile.Address.Zone,
	}

	return c.JSON(response.Payload{
		Message: "Data Found",
		Data:    fiber.Map{"user": res, "address": address},
	})
}
