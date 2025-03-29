package main

import (
	"log"

	"github.com/AkifSahn/pollution-tracker/config"
	"github.com/AkifSahn/pollution-tracker/internal/database"
	"github.com/AkifSahn/pollution-tracker/internal/ingest"
	"github.com/AkifSahn/pollution-tracker/internal/pollution"
	"github.com/AkifSahn/pollution-tracker/internal/rabbitmq"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"

	_ "github.com/AkifSahn/pollution-tracker/docs"
)

//	@title			pollution-tracker API
//	@description	API documentation for pollution-tracker app
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

	app.Use(logger.New())
	pollution.SetupRoutes(app)

	app.Get("/swagger/*", swagger.HandlerDefault)

	err := app.Listen(":" + cfg.ServerPort)
	if err != nil {
		panic(err)
	}
}
