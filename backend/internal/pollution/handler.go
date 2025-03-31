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

// GetAirQualityLocation
//
//	@Summary		Gets air quality value
//	@Description	Gets air quality value for given location. Optional time range can be given
//	@Tags			pollution
//	@Produce		json
//	@Param			latitude	path		string						true	"latitude"
//	@Param			longitude	path		string						true	"longitude"
//
//	@Param			from		query		string						false	"Starting time"
//	@Param			to			query		string						false	"End time"
//	@Failure		400			{string}	string						"Failed to parse given time - from"
//	@Failure		400			{string}	string						"Failed to parse given time - to"
//	@Failure		400			{string}	string						"Cannot parse the given latitude/longitude value!"
//	@Failure		500			{string}	string						"Failed to fetch data from database!"
//	@Success		200			{object}	[]PollutionValueResponse	"values: []PollutionValueResponse"
//	@Router			/api/air-quality/{latitude}/{longitude} [get]
func GetAirQualityLocation(c *fiber.Ctx) error {

	latStr := c.Params("latitude")
	longStr := c.Params("longitude")

	// Get `latitude` and `longitude` from URL params and parse them
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
	// --------

	format := "2006-01-02 15:04:05"
	// Get `from` and `to` from URL query and parse them
	fromStr := c.Query("from", time.Now().Add(-24*time.Hour).Format(format))
	toStr := c.Query("to", time.Now().Format(format))

	from, err := time.Parse(format, fromStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse given time - from",
		})
	}

	to, err := time.Parse(format, toStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse given time - to",
		})
	}
	// --------

	repo := NewPollutionRepo(database.DB)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	vals, err := repo.GetPollutionValueByPosition(ctx, latitude, longitude, from, to)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch data from database! " + err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"values": vals,
	})
}

// GetAnomaliesOfRange
//
//	@Summary		Gets anomalies for range
//	@Description	Gets anomalies for a given time frange
//	@Tags			pollution
//	@Produce		json
//	@Param			from	query		string		true	"from"
//	@Param			to		query		string		true	"to"
//	@Success		400		{string}	string		"Failed to parse given time value(from)!"
//	@Success		400		{string}	string		"Failed to parse given time value(from)!"
//	@Success		500		{string}	string		"Failed to fetch pollution entries from database"
//	@Success		200		{object}	[]Pollution	"pollution"
//	@Router			/api/anomalies [get]
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

// GetPollutionDensityRegion
//
//	@Summary		Gets pollution density
//	@Description	Gets pollution density for a given region
//	@Tags			pollution
//	@Produce		json
//	@Param			region	path		string	true	"region"
//	@Param			from	query		string	true	"from"
//	@Param			to		query		string	true	"to"
//
//	@Success		400		{string}	string	"Failed to parse given time value(from)!"
//	@Success		400		{string}	string	"Failed to parse given time value(to)!"
//	@Success		500		{string}	string	"Failed to fetch region density from database"
//	@Success		200		{string}	string	"density"
//	@Router			/api/regions-density/{region} [get]
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

// PostPollutionEntry
//
//	@Summary		Posts pollution entry
//	@Description	Posts a new pollution entry
//	@Tags			pollution
//	@Accept			json
//	@Produce		json
//	@Param			request	body		Pollution	true	"Request of adding a new pollution entry"
//	@Success		400		{string}	string		"Failed to parse request body"
//	@Success		400		{string}	string		"Failed to marshal request body"
//	@Success		500		{string}	string		"Failed to publish pollution entry to RabbitMQ queue"
//	@Success		200		{string}	string		"Successfully received the pollution entry"
//	@Router			/api/ingest/manual [post]
func PostPollutionEntry(c *fiber.Ctx) error {
	var body Pollution

	body.Time = time.Now()
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
