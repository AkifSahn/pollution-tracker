package main

import (
	"log"

	"github.com/AkifSahn/pollution-tracker/config"
	"github.com/AkifSahn/pollution-tracker/internal/database"
	"github.com/AkifSahn/pollution-tracker/internal/ingest"
	"github.com/AkifSahn/pollution-tracker/internal/notification"
	"github.com/AkifSahn/pollution-tracker/internal/pollution"
	"github.com/AkifSahn/pollution-tracker/internal/rabbitmq"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
	"github.com/gofiber/websocket/v2"

	_ "github.com/AkifSahn/pollution-tracker/docs"
)

// @title			pollution-tracker API
// @description	API documentation for pollution-tracker app
func main() {

	app := fiber.New(fiber.Config{
		StrictRouting: true,
	})

	cfg := config.LoadConfig()

	log.Printf("Connecting to database")
	database.InitDB(cfg) // add defer close

	// Init rabbitmq
	log.Printf("Initializing RabbitMQ")
	rabbitmq.Connect(cfg)
	rabbitmq.DeclareQueues()
	defer rabbitmq.AmqpConn.Close()
	defer rabbitmq.AmqpCh.Close()

	go ingest.ListenIngestion()

	hub := notification.NewHub()
	go hub.Run()

	go notification.ListenAndConsumeNotifications(hub)

	app.Use(cors.New())
	app.Use(logger.New())

	// Middleware to upgrade HTTP to WebSocket if requested
	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	pollution.SetupRoutes(app)
	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		notification.NewWs(hub, c)
	}))

	app.Get("/swagger/*", swagger.HandlerDefault)

	err := app.Listen(":" + cfg.ServerPort)
	if err != nil {
		panic(err)
	}
}
