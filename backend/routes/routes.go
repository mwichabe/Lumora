package routes

import (
	"github.com/gofiber/fiber/v2"

	"lumora/backend/config"
	"lumora/backend/controllers"
	"lumora/backend/middleware"
)

// Register wires every API route onto the Fiber app.
func Register(app *fiber.App, cfg config.Config) {
	api := app.Group("/api")

	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok", "service": "lumora"})
	})

	auth := &controllers.AuthController{Cfg: cfg}
	api.Post("/auth/register", auth.Register)
	api.Post("/auth/login", auth.Login)

	// Protected routes require a valid Bearer token.
	protected := api.Group("", middleware.Protected(cfg.JWTSecret))

	protected.Get("/auth/me", auth.Me)
	protected.Post("/auth/setup", auth.Setup)

	lessons := &controllers.LessonController{}
	protected.Get("/skills", lessons.GalaxyMap)
	protected.Get("/lessons/:id", lessons.GetLesson)

	progress := &controllers.ProgressController{}
	protected.Get("/home", progress.Home)
	protected.Post("/lessons/:id/complete", progress.CompleteLesson)

	listening := &controllers.ListeningController{}
	protected.Get("/listening", listening.List)
	protected.Get("/listening/:id", listening.Get)
	protected.Post("/listening/:id/complete", listening.Complete)

	reading := &controllers.ReadingController{}
	protected.Get("/reading", reading.List)
	protected.Get("/reading/:id", reading.Get)
	protected.Post("/reading/:id/complete", reading.Complete)

	quests := &controllers.QuestController{}
	protected.Get("/quests/daily", quests.Daily)

	characters := &controllers.CharacterController{}
	protected.Get("/characters", characters.List)

	leaderboard := &controllers.LeaderboardController{}
	protected.Get("/leaderboard", leaderboard.League)
}
