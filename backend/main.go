package main

import (
	"github.com/AkifSahn/pollution-tracker/config"
	"github.com/AkifSahn/pollution-tracker/internal/database"
	"github.com/AkifSahn/pollution-tracker/internal/pollution"
	"github.com/gofiber/fiber/v2"
)

func main() {

	app := fiber.New(fiber.Config{
		StrictRouting: true,
	})

	cfg := config.LoadConfig()
	database.InitDB(cfg)
	pollution.SetupRoutes(app)

	err := app.Listen(":" + cfg.ServerPort)
	if err != nil {
		panic(err)
	}
}
