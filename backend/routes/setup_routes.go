package routes

import (
	"github.com/AkifSahn/pollution-tracker/controllers"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {

	// // Endpoints for auth handlers
	app.Get("/api/air-quality", controllers.GetAirQualityLocation)
	app.Get("/api/anomalies-range", controllers.ListAnomaliesOfRange)
	app.Get("api/pollution-density-region", controllers.GetPollutionDensityRegion)
}
