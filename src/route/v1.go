package route

import (
	"github.com/emamulandalib/airbringr-auth/handler"
	"github.com/gofiber/fiber/v2"
)

func V1(server *fiber.App, handler *handler.Handler) {
	v1 := server.Group("/v1")
	v1.Get("/", handler.Home)
	v1.Post("/signup", handler.Signup)
}
