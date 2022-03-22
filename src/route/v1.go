package route

import (
	"github.com/emamulandalib/airbringr-auth/config"
	"github.com/emamulandalib/airbringr-auth/handler"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
)

func V1(server *fiber.App, handler *handler.Handler) {
	v1 := server.Group("/v1")
	v1.Get("/", handler.Home)

	v1.Post("/signup", handler.Signup)
	v1.Post("/verify-email", handler.EmailVerification)

	v1.Post("/send-mobile-verification-otp", handler.MobileVerificationOTP)
	v1.Post("/verify-mobile", handler.VerifyMobile)

	v1.Post("/send-sms-otp", handler.SendSmsOtp)
	v1.Post("/send-email-otp", handler.SendEmailOTP)
	v1.Post("/verify-otp", handler.VerifyOtp)

	v1.Post("/login", handler.Login)
	v1.Post("/password-reset-email-link", handler.PasswordResetEmailLink)
	v1.Post("/password-reset", handler.PasswordReset)

	// JWT Middleware
	v1.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte(config.Params.JWTSignatureSecret),
	}))

	v1.Get("/view-profile", handler.ViewProfile)
	v1.Post("/edit-profile", handler.EditProfile)
	v1.Post("/verify-password", handler.VerifyPassword)
	v1.Post("/edit-email-req", handler.EditEmailReq)
	v1.Post("/verify-and-update-email", handler.VerifyAndEditEmail)
	v1.Post("/edit-mobile-send-otp", handler.MobileVerificationOTP)
	v1.Post("/verify-and-update-mobile", handler.VerifyMobile)
}
