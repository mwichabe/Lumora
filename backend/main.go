package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"lumora/backend/config"
	"lumora/backend/controllers"
	"lumora/backend/database"
	"lumora/backend/routes"
)

func main() {
	cfg := config.Load()

	database.Connect(cfg.DBPath)

	app := fiber.New(fiber.Config{
		AppName: "Lumora API",
	})

	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: cfg.CORSOrigins,
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, PATCH, DELETE, OPTIONS",
	}))

	// Serve uploaded profile photos.
	app.Static("/uploads", cfg.UploadsDir)

	routes.Register(app, cfg)

	// Background loop that casually pushes tips / announcements to users.
	controllers.StartNotificationScheduler()

	log.Printf("Lumora API listening on :%s", cfg.Port)
	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
