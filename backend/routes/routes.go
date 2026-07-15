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
	// Public: rendered by a plain <img>, which can't send a bearer token.
	api.Get("/avatars/:id", auth.GetAvatar)
	api.Post("/auth/register", auth.Register)
	api.Post("/auth/login", auth.Login)
	api.Post("/auth/forgot-password", auth.ForgotPassword)
	api.Post("/auth/reset-password", auth.ResetPassword)

	// Public certificate verification — anyone with a serial can confirm it.
	publicExam := &controllers.ExamController{}
	api.Get("/verify/:serial", publicExam.Verify)

	// Paystack webhook — public, but authenticated by the signature header.
	payments := &controllers.PaymentController{Cfg: cfg}
	api.Post("/paystack/webhook", payments.Webhook)

	// Protected routes require a valid Bearer token.
	protected := api.Group("", middleware.Protected(cfg.JWTSecret))

	protected.Get("/auth/me", auth.Me)
	protected.Post("/auth/setup", auth.Setup)
	protected.Patch("/auth/profile", auth.UpdateProfile)
	protected.Post("/auth/avatar", auth.UploadAvatar)
	protected.Delete("/auth/avatar", auth.RemoveAvatar)
	protected.Post("/auth/password", auth.ChangePassword)
	protected.Delete("/auth/account", auth.DeleteAccount)

	lessons := &controllers.LessonController{}
	protected.Get("/skills", lessons.GalaxyMap)
	protected.Get("/lessons/:id", lessons.GetLesson)

	progress := &controllers.ProgressController{}
	protected.Get("/home", progress.Home)
	protected.Post("/lessons/:id/complete", progress.CompleteLesson)

	hearts := &controllers.HeartsController{}
	protected.Get("/hearts", hearts.Status)
	protected.Post("/hearts/lose", hearts.Lose)

	listening := &controllers.ListeningController{}
	protected.Get("/listening", listening.List)
	protected.Get("/listening/:id", listening.Get)
	protected.Post("/listening/:id/complete", listening.Complete)

	reading := &controllers.ReadingController{}
	protected.Get("/reading", reading.List)
	protected.Get("/reading/:id", reading.Get)
	protected.Post("/reading/:id/complete", reading.Complete)

	enroll := &controllers.EnrollmentController{}
	protected.Get("/enrollments", enroll.List)
	protected.Post("/enrollments", enroll.Enroll)
	protected.Post("/enrollments/active", enroll.SetActive)

	practice := &controllers.PracticeController{}
	protected.Get("/practice", practice.Pool)
	protected.Get("/practice/listening", practice.Listening)
	protected.Get("/practice/reading", practice.Reading)
	protected.Post("/practice/complete", practice.Complete)
	protected.Post("/mistakes", practice.RecordMistake)
	protected.Post("/mistakes/resolve", practice.ResolveMistakes)

	quests := &controllers.QuestController{}
	protected.Get("/quests/daily", quests.Daily)

	characters := &controllers.CharacterController{}
	protected.Get("/characters", characters.List)

	notifications := &controllers.NotificationController{}
	protected.Get("/notifications", notifications.List)
	protected.Post("/notifications/read", notifications.MarkRead)
	protected.Post("/notifications/:id/read", notifications.MarkOneRead)
	protected.Delete("/notifications/:id", notifications.Delete)

	chat := &controllers.ChatController{}
	protected.Get("/chat/contacts", chat.Contacts)
	protected.Get("/chat/threads", chat.Threads)
	protected.Get("/chat/unread", chat.Unread)
	protected.Get("/chat/with/:id", chat.Messages)
	protected.Post("/chat/with/:id", chat.Send)

	protected.Get("/payments/status", payments.Status)
	protected.Post("/payments/initialize", payments.Initialize)
	protected.Get("/payments/verify", payments.Verify)

	exam := &controllers.ExamController{}
	protected.Get("/exam/meta", exam.Meta)
	protected.Get("/exam/paper", exam.Paper)
	protected.Post("/exam/start", exam.Start)
	protected.Post("/exam/submit", exam.Submit)
	protected.Get("/certificates", exam.ListCertificates)
	protected.Get("/certificates/:id", exam.GetCertificate)
	protected.Delete("/certificates/:id", exam.DeleteCertificate)

	leaderboard := &controllers.LeaderboardController{}
	protected.Get("/leaderboard", leaderboard.League)
}
