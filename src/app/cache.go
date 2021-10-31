package app

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
)

func (a *App) SetCacheMiddleware(s *fiber.App) {
	s.Use(cache.New(cache.Config{
		Next:    func(c *fiber.Ctx) bool { return false },
		Storage: a.RedisStorage(),
		KeyGenerator: func(c *fiber.Ctx) string {
			userIP := string(c.Request().Header.Peek("X-Forwarded-For"))
			return fmt.Sprintf("ab-auth-cache-%s-%s", userIP, c.Path())
		},
	}))
}
