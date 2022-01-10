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

	//client, err := mongo.NewClient(options.Client().
	//	ApplyURI("mongodb+srv://airbringr:EumNfKThcgIeqz8o@cluster0.nqgzx.mongodb.net/myFirstDatabase?retryWrites=true&w=majority"))
	//if err != nil {
	//	log.Fatal(err)
	//}
	//ctx, _ := context.WithTimeout(context.Background(), 20*time.Second)
	//err = client.Connect(ctx)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer client.Disconnect(ctx)
	//err = client.Ping(ctx, readpref.Primary())
	//if err != nil {
	//	log.Fatal(err)
	//}
	//databases, err := client.ListDatabaseNames(ctx, bson.M{})
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Println(databases)
	//fmt.Println("Mongo Cloud Connected")

	config.New()
	App := app.New()
	App.Bootstrap()
	defer func() { _ = App.Mongo.Disconnect() }()

	// setup App
	server := fiber.New(fiber.Config{
		IdleTimeout: idleTimeout,
	})

	// setup middlewares
	server.Use(requestid.New())
	server.Use(recover.New())
	server.Use(cors.New(cors.Config{
		AllowOrigins: config.Params.CORSPermitted,
	}))

	server.Use(logger.New(logger.Config{
		Format:   "[${time}] ${status} ${locals:requestid} - ${latency} ${method} ${path}\n",
		TimeZone: "Asia/Dhaka",
	}))

	App.SetCacheMiddleware(server)
	App.SetRateLimiterMiddleware(server)

	//routes
	Handler := handler.New()
	server.Get("/", Handler.Home)
	route.V1(server, Handler)
	server.Use(Handler.NotFound) // 404

	// Listen from a different goroutine
	go func() {
		if err := server.Listen(fmt.Sprintf(":%d", config.Params.Port)); err != nil {
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
