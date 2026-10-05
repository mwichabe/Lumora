package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"lumora/backend/config"
	"lumora/backend/controllers"
	"lumora/backend/database"
	"lumora/backend/routes"
	"lumora/backend/utils"
)

func main() {
	cfg := config.Load()

	database.Connect(cfg.DatabaseURL, cfg.DBPath)

	// Undo levels and level-up notifications issued by the old XP-based rule.
	controllers.RepairLevels()

	app := fiber.New(fiber.Config{
		AppName: "Lumora API",
	})

	app.Use(recover.New())
	app.Use(logger.New())
	// Course payloads are repetitive JSON that shrinks ~4x under gzip — most of
	// the wait on a phone's mobile-data connection is the download itself.
	app.Use(compress.New(compress.Config{Level: compress.LevelBestSpeed}))
	app.Use(cors.New(cors.Config{
		AllowOrigins: cfg.CORSOrigins,
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, PATCH, DELETE, OPTIONS",
	}))

	// Profile photos live in the database and are served by GET /api/avatars/:id,
	// so there is no upload directory to expose here.

	// Message translation. Offline language detection always runs; the model
	// call only happens when an API key is configured.
	controllers.InitTranslation(utils.NewTranslator(cfg))
	controllers.InitWritingCoach(utils.NewWritingCoach(cfg))

	routes.Register(app, cfg)

	// Background loop that casually pushes tips / announcements to users.
	controllers.StartNotificationScheduler()

	// Closes out finished league weeks. Seasons also settle lazily when anyone
	// opens the league, so a sleeping free instance still resolves correctly.
	controllers.StartLeagueScheduler()

	log.Printf("Lumora API listening on :%s", cfg.Port)
	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
