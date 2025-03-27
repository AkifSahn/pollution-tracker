package pollution

import (
	"encoding/json"

	"github.com/AkifSahn/pollution-tracker/internal/rabbitmq"
	"github.com/gofiber/fiber/v2"
	amqp "github.com/rabbitmq/amqp091-go"
)

func SetupRoutes(app *fiber.App) {

	api := app.Group("/api")

	// // Endpoints for auth handlers
	api.Get("air-quality/:location", GetAirQualityLocation)
	api.Get("anomalies", GetAnomaliesOfRange)
	api.Get("regions-density/:region", GetPollutionDensityRegion)

	api.Post("ingest/manual", PostPollutionEntry)
}

func GetAirQualityLocation(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "GetAirQualityLocation",
	})
}

func GetAnomaliesOfRange(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "GetAnomaliesOfRange",
	})
}

func GetPollutionDensityRegion(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "GetPollutionDensityRegion",
	})
}

func PostPollutionEntry(c *fiber.Ctx) error {
	var body Pollution

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request body" + err.Error(),
		})
	}

	var msg []byte
	msg, err := json.Marshal(&body)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to marshal request body" + err.Error(),
		})
	}

	err = rabbitmq.AmqpCh.Publish(
		"",
		"ingest_queue",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        msg,
		})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to publish pollution entry to RabbitMQ queue",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Successfully received the pollution entry",
	})
}
