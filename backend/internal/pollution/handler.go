package pollution

import (
	"context"
	"encoding/json"
	"time"

	"github.com/AkifSahn/pollution-tracker/internal/database"
	"github.com/AkifSahn/pollution-tracker/internal/rabbitmq"
	"github.com/gofiber/fiber/v2"
	amqp "github.com/rabbitmq/amqp091-go"
)

func SetupRoutes(app *fiber.App) {

	api := app.Group("/api")

	api.Post("pollutions", PostPollutionEntry)

	api.Get("pollutions", GetAllPolutions)
	api.Get("pollutions/density/rect", GetPollutionDensityOfRect)
	api.Get("pollutions/:latitude/:longitude", GetPollutionsByLatLon)

	api.Get("anomalies", GetAnomaliesOfRange)
	api.Get("pollutants", GetPollutants)
}

// PostPollutionEntry
//
//	@Summary		Posts pollution entry
//	@Description	Posts a new pollution entry
//	@Tags			pollutions
//	@Accept			json
//	@Produce		json
//	@Param			request	body		Pollution	true	"Request of adding a new pollution entry"
//	@Success		400		{string}	string		"Failed to parse request body"
//	@Success		400		{string}	string		"Failed to marshal request body"
//	@Success		500		{string}	string		"Failed to publish pollution entry to RabbitMQ queue"
//	@Success		200		{string}	string		"Successfully received the pollution entry"
//	@Router			/api/pollutions [post]
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

// GetAllPollutions
//
//	@Summary		Gets pollution values
//	@Description	Gets all pollution values for given time range.
//	@Tags			pollutions
//	@Produce		json
//
//	@Param			from		query		string					true	"Start time"
//	@Param			to			query		string					true	"End time"
//	@Param			pollutant	query		string					false	"Pollutant"
//
//	@Success		200			{object}	map[string][]Pollution	"Pollution values"
//	@Failure		400			{object}	map[string]string		"Invalid params"
//	@Failure		500			{object}	map[string]string		"Failed to fetch pollution entries from database"
//	@Router			/api/pollutions [get]
func GetAllPolutions(c *fiber.Ctx) error {
	fromStr := c.Query("from")
	toStr := c.Query("to")

	// Parse the string times into time.Time
	var from, to time.Time
	ok, msg := ParseTimeRange(fromStr, toStr, &from, &to)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": msg,
		})
	}

	pollutant := c.Query("pollutant")

	repo := NewPollutionRepo(database.DBPool)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pollutions, err := repo.GetAllPolutionWithinTimeRange(ctx, from, to, pollutant)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch pollution entries from database ",
		})
	}

	return c.JSON(fiber.Map{
		"data": pollutions,
	})
}

// GetPollutionsByLatLon
//
//	@Summary		Gets pollution values
//	@Description	Gets pollution values for given location and time range
//	@Tags			pollutions
//	@Produce		json
//
//	@Param			latitude	path		string								true	"latitude"
//	@Param			longitude	path		string								true	"longitude"
//	@Param			from		query		string								false	"Start time"
//	@Param			to			query		string								false	"End time"
//
//	@Failure		400			{object}	map[string]string					"Invalid params"
//	@Failure		500			{object}	map[string]string					"Failed to fetch pollution entries from database"
//	@Success		200			{object}	map[string][]PollutionValueResponse	"Pollution Values"
//	@Router			/api/pollutions/{latitude}/{longitude} [get]
func GetPollutionsByLatLon(c *fiber.Ctx) error {

	latStr := c.Params("latitude")
	longStr := c.Params("longitude")

	var latitude, longitude float64
	ok, msg := ParseLatLon(latStr, longStr, &latitude, &longitude)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": msg,
		})
	}

	// TODO: should we do the 24 hour default time range????
	fromStr := c.Query("from", time.Now().Add(-24*time.Hour).Format(TimeFormat))
	toStr := c.Query("to", time.Now().Format(TimeFormat))

	var from, to time.Time
	ok, msg = ParseTimeRange(fromStr, toStr, &from, &to)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": msg,
		})
	}

	repo := NewPollutionRepo(database.DBPool)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	vals, err := repo.GetPollutionValueByPosition(ctx, latitude, longitude, from, to)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch data from database! " + err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": vals,
	})
}

// GetPollutionDensityOfRect
//
//	@Summary		Gets pollution densities of rect
//	@Description	Gets pollution densities for a given rect and time range
//	@Tags			pollutions
//	@Produce		json
//
//	@Param			latFrom		query		float64							true	"latFrom"
//	@Param			latTo		query		float64							true	"latTo"
//	@Param			longFrom	query		float64							true	"longFrom"
//	@Param			longTo		query		float64							true	"longTo"
//	@Param			from		query		string							true	"from"
//	@Param			to			query		string							true	"to"
//	@Param			pollutant	query		string							false	"pollutant"
//
//	@Failure		400			{object}	map[string]string				"Invalid params"
//	@Failure		500			{object}	map[string]string				"Failed to fetch pollution entries from database"
//	@Success		200			{object}	map[string][]PollutionDensity	"Pollution Densities"
//	@Router			/api/pollutions/density/rect [get]
func GetPollutionDensityOfRect(c *fiber.Ctx) error {
	fromStr := c.Query("from", time.Now().Add(-24*time.Hour).Format(TimeFormat))
	toStr := c.Query("to", time.Now().Format(TimeFormat))

	latFrom := c.QueryFloat("latFrom")
	latTo := c.QueryFloat("latTo")

	longFrom := c.QueryFloat("longFrom")
	longTo := c.QueryFloat("longTo")

	pollutant := c.Query("pollutant")

	var from, to time.Time
	ok, msg := ParseTimeRange(fromStr, toStr, &from, &to)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": msg,
		})
	}

	repo := NewPollutionRepo(database.DBPool)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	densities, err := repo.GetPollutionDensityOfRect(ctx, latFrom, latTo, longFrom, longTo, from, to, 5*time.Minute, pollutant)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch rect densities from database: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": densities,
	})
}

// GetAnomaliesOfRange
//
//	@Summary		Gets anomalies for range
//	@Description	Gets anomalies for a given time range
//	@Tags			anomalies
//	@Produce		json
//
//	@Param			from	query		string					true	"Start time"
//	@Param			to		query		string					true	"End time"
//
//	@Failure		400		{object}	map[string]string		"Invalid params"
//	@Failure		500		{object}	map[string]string		"Failed to fetch pollution entries from database"
//	@Success		200		{object}	map[string][]Pollution	"Anomalies"
//	@Router			/api/anomalies [get]
func GetAnomaliesOfRange(c *fiber.Ctx) error {
	fromStr := c.Query("from")
	toStr := c.Query("to")

	// Parse the string times into time.Time
	var from, to time.Time
	ok, msg := ParseTimeRange(fromStr, toStr, &from, &to)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": msg,
		})
	}

	repo := NewPollutionRepo(database.DBPool)
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
		"data": pollutions,
	})
}

// GetPollutants
//
//	@Summary		Gets pollutants
//	@Description	Gets distinct pollutants that exists in database
//	@Tags			pollutants
//	@Produce		json
//
//	@Failure		500	{object}	map[string]string	"Failed to fetch pollutants entries from database"
//	@Success		200	{object}	map[string]string	"Pollutants"
//	@Router			/api/pollutants [get]
func GetPollutants(c *fiber.Ctx) error {
	repo := NewPollutionRepo(database.DBPool)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pollutants, err := repo.GetDistinctPollutants(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch pollutants from database: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": pollutants,
	})

}
