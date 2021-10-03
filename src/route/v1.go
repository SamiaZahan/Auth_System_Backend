package route

import (
	"github.com/emamulandalib/airbringr-auth/handler"
	"github.com/gofiber/fiber/v2"
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
	v1.Post("/verify-sms-otp", handler.VerifyOtp)
}
