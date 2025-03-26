package controllers

import "github.com/gofiber/fiber/v2"

func GetAirQualityLocation(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "GetAirQualityLocation",
	})
}

func ListAnomaliesOfRange(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "ListAnomaliesOfRange",
	})
}

func GetPollutionDensityRegion(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "GetPollutionDensityRegion",
	})
}
