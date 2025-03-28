package pollution

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/AkifSahn/pollution-tracker/internal/database"
	"github.com/AkifSahn/pollution-tracker/internal/rabbitmq"
	"github.com/gofiber/fiber/v2"
	amqp "github.com/rabbitmq/amqp091-go"
)

func SetupRoutes(app *fiber.App) {

	api := app.Group("/api")

	// // Endpoints for auth handlers
	api.Get("air-quality/:latitude/:longitude", GetAirQualityLocation)
	api.Get("anomalies", GetAnomaliesOfRange)
	api.Get("regions-density/:region", GetPollutionDensityRegion)

	api.Post("ingest/manual", PostPollutionEntry)
}

func GetAirQualityLocation(c *fiber.Ctx) error {

	latStr := c.Params("latitude")
	longStr := c.Params("longitude")

	latitude, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse the given latitude value!",
		})
	}

	longitude, err := strconv.ParseFloat(longStr, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse the given longitude value!",
		})
	}

	repo := NewPollutionRepo(database.DB)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	val, err := repo.GetPollutionValueByPosition(ctx, latitude, longitude)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch data from database! " + err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"value": val,
	})
}

func GetAnomaliesOfRange(c *fiber.Ctx) error {
	fromStr := c.Query("from")
	toStr := c.Query("to")

	// Parse the string times into time.Time
	format := "2006-01-02 15:04:05"
	from, err := time.Parse(format, fromStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse given time value(from)!",
		})
	}

	to, err := time.Parse(format, toStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse given time value(to)!",
		})
	}

	repo := NewPollutionRepo(database.DB)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get anomalies from database
	pollutions, err := repo.GetAnomaliesWithinTimeRange(ctx, from, to)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch pollution entries from database " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"pollutions": pollutions,
	})
}

func GetPollutionDensityRegion(c *fiber.Ctx) error {
	region := c.Params("region")
	fromStr := c.Query("from")
	toStr := c.Query("to")

	// Parse the string times into time.Time
	format := "2006-01-02 15:04:05"
	from, err := time.Parse(format, fromStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse given time value(from)!",
		})
	}

	to, err := time.Parse(format, toStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse given time value(to)!",
		})
	}
	// --------

	repo := NewPollutionRepo(database.DB)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	density, err := repo.GetPollutionDensityByRegion(ctx, region, from, to)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch region density from database",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"region":  region,
		"density": density,
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
