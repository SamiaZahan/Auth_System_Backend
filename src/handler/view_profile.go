package handler

import (
	"context"
	"github.com/emamulandalib/airbringr-auth/repository"
	"github.com/emamulandalib/airbringr-auth/response"
	"github.com/emamulandalib/airbringr-auth/service"
	"github.com/gofiber/fiber/v2"
	"time"
)

func (h *Handler) ViewProfile(c *fiber.Ctx) error {
	token := c.Get("Authorization", "")
	ExtractJwtClaims := service.ExtractJwtClaims{}
	claims, ok := ExtractJwtClaims.ExtractClaims(token)
	if !ok {
		return c.JSON(response.Payload{
			Message: "Something  Wrong",
		})
	}
	email := claims["email"].(string)
	var ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	authRepo := repository.Auth{Ctx: ctx}
	user, _ := authRepo.GetUserByEmail(email)
	userProfile, _ := authRepo.GetUserProfileByID(user.ID.Hex())
	res := map[string]string{
		"image":      userProfile.ProfilePicURI,
		"first_name": userProfile.FirstName,
		"last_name":  userProfile.LastName,
		"gender":     userProfile.Gender,
		"email":      user.Email,
		"mobile":     user.Mobile,
		"address":    userProfile.Address,
	}
	return c.JSON(response.Payload{
		Message: "Data Found",
		Data:    res,
	})
}
