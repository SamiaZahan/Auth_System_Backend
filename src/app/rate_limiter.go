package app

import (
	"fmt"
	"time"

	"github.com/emamulandalib/airbringr-auth/dto"
	"github.com/emamulandalib/airbringr-auth/response"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

func (a *App) SetRateLimiterMiddleware(app *fiber.App) {
	// rete limit for SEND SMS OTP
	app.Use(limiter.New(limiter.Config{
		Next: func(c *fiber.Ctx) bool {
			if c.Path() == "/v1/send-sms-otp" {
				d := new(dto.SendSmsOtpInput)
				if err := c.BodyParser(d); err != nil {
					return true
				}
				return false
			}

			return true
		},
		Max:        2,
		Expiration: time.Minute * 60,
		LimitReached: func(c *fiber.Ctx) error {
			return c.
				Status(fiber.StatusTooManyRequests).
				JSON(response.Payload{
					Message: "Too many requests. Please try again after 1 hour or you can try to login using email.",
				})
		},
		KeyGenerator: func(c *fiber.Ctx) string {
			userIP := string(c.Request().Header.Peek("X-Forwarded-For"))
			d := new(dto.SendSmsOtpInput)
			c.BodyParser(d)
			return fmt.Sprintf("ab-auth-rl-sms-%s-%s", userIP, d.Mobile)
		},
		Storage: a.RedisStorage(),
	}))

	// rete limit for SEND Email OTP
	app.Use(limiter.New(limiter.Config{
		Next: func(c *fiber.Ctx) bool {
			if c.Path() == "/v1/send-email-otp" {
				d := new(dto.EmailOtpInput)
				if err := c.BodyParser(d); err != nil {
					return true
				}
				return false
			}

			return true
		},
		Max:        3,
		Expiration: time.Minute * 60,
		LimitReached: func(c *fiber.Ctx) error {
			return c.
				Status(fiber.StatusTooManyRequests).
				JSON(response.Payload{
					Message: "Too many requests. Please try again after 1 hour or you can try to login using mobile number",
				})
		},
		KeyGenerator: func(c *fiber.Ctx) string {
			userIP := string(c.Request().Header.Peek("X-Forwarded-For"))
			d := new(dto.EmailOtpInput)
			c.BodyParser(d)
			return fmt.Sprintf("ab-auth-rl-email-%s-%s", userIP, d.Email)
		},
		Storage: a.RedisStorage(),
	}))
}
