package main

import (
	"github.com/AkifSahn/pollution-tracker/config"
	"github.com/AkifSahn/pollution-tracker/routes"
	"github.com/gofiber/fiber/v2"
)

func main() {

	app := fiber.New(fiber.Config{
		StrictRouting: true,
	})

	cfg := config.LoadConfig()
	// database.InitDB(cfg)

	routes.SetupRoutes(app)

	err := app.Listen(":" + cfg.ServerPort)
	if err != nil {
		panic(err)
	}
}
