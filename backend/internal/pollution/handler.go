package pollution

import "github.com/gofiber/fiber/v2"

func SetupRoutes(app *fiber.App) {

	api := app.Group("/api")

	// // Endpoints for auth handlers
	api.Get("/air-quality/:location", GetAirQualityLocation)
	api.Get("anomalies", GetAnomaliesOfRange)
	api.Get("regions-density/:region", GetPollutionDensityRegion)
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
