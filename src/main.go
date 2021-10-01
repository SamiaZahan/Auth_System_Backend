package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/emamulandalib/airbringr-auth/app"
	"github.com/emamulandalib/airbringr-auth/config"
	"github.com/emamulandalib/airbringr-auth/handler"
	"github.com/emamulandalib/airbringr-auth/route"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	log "github.com/sirupsen/logrus"
)

const idleTimeout = 5 * time.Second

func init() {
	// log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.InfoLevel)
	log.SetReportCaller(true)
}

func main() {
	config.New()
	App := app.New()
	App.Bootstrap()
	defer func() { _ = App.Mongo.Disconnect() }()

	Handler := handler.New()

	// setup App
	server := fiber.New(fiber.Config{
		IdleTimeout: idleTimeout,
	})

	// setup middlewares
	server.Use(requestid.New())
	server.Use(recover.New())
	server.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000, https://glistening-sack.surge.sh",
	}))

	server.Use(logger.New(logger.Config{
		Format:   "[${time}] ${status} ${locals:requestid} - ${latency} ${method} ${path}\n",
		TimeZone: "Asia/Dhaka",
	}))

	//routes
	server.Get("/", Handler.Home)
	route.V1(server, Handler)

	// 404
	server.Use(Handler.NotFound)

	// Listen from a different goroutine
	go func() {
		if err := server.Listen(fmt.Sprintf("0.0.0.0:%d", config.Params.Port)); err != nil {
			log.Fatal(err.Error())
		}
	}()

	rand.Seed(time.Now().UnixNano())
	c := make(chan os.Signal, 1)                    // Create channel to signify a signal being sent
	signal.Notify(c, os.Interrupt, syscall.SIGTERM) // When an interrupt or termination signal is sent, notify the channel

	<-c // This blocks the main thread until an interrupt is received
	fmt.Println("Gracefully shutting down...")
	_ = server.Shutdown()

	fmt.Println("Running cleanup tasks...")

	// Your cleanup tasks go here

	fmt.Println("Fiber was successful shutdown.")
}
