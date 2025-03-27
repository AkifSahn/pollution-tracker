package main

import (
	"log"

	"github.com/AkifSahn/pollution-tracker/config"
	"github.com/AkifSahn/pollution-tracker/internal/database"
	"github.com/AkifSahn/pollution-tracker/internal/pollution"
	"github.com/AkifSahn/pollution-tracker/internal/rabbitmq"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

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

	msgs, err := rabbitmq.AmqpCh.Consume(
		"ingest_queue", // queue
		"",             // consumer
		true,           // auto-ack
		false,          // exclusive
		false,          // no-local
		false,          // no-wait
		nil,            // args
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %s", err.Error())
	}

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	app.Use(logger.New())
	pollution.SetupRoutes(app)

	err = app.Listen(":" + cfg.ServerPort)
	if err != nil {
		panic(err)
	}
}
