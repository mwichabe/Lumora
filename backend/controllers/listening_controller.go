package controllers

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"lumora/backend/database"
	"lumora/backend/middleware"
	"lumora/backend/models"
)

// ListeningController serves unit-level listening sessions (character dialogues
// + comprehension questions) and records their completion.
type ListeningController struct{}

func orderByIndex(db *gorm.DB) *gorm.DB { return db.Order("order_index asc") }

func hydrateQuestions(qs []models.ListeningQuestion) {
	for i := range qs {
		var o []string
		if qs[i].OptionsJSON != "" {
			_ = json.Unmarshal([]byte(qs[i].OptionsJSON), &o)
		}
		qs[i].Options = o
	}
}

// List returns every listening session for the user's target language (falling
// back to the seeded Spanish course), with lines and questions hydrated.
func (l *ListeningController) List(c *fiber.Ctx) error {
	user := middleware.CurrentUser(c)
	lang := user.TargetLanguage
	if lang == "" {
		lang = "es"
	}

	var sessions []models.ListeningSession
	q := database.DB.
		Preload("Matches", orderByIndex).
		Preload("Lines", orderByIndex).
		Preload("Questions", orderByIndex).
		Order("order_index asc")

	q.Where("language = ?", lang).Find(&sessions)

	for i := range sessions {
		hydrateQuestions(sessions[i].Questions)
	}
	return c.JSON(fiber.Map{"sessions": sessions})
}

// Get returns a single listening session by id.
func (l *ListeningController) Get(c *fiber.Ctx) error {
	id := c.Params("id")

	var session models.ListeningSession
	if err := database.DB.
		Preload("Matches", orderByIndex).
		Preload("Lines", orderByIndex).
		Preload("Questions", orderByIndex).
		First(&session, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "session not found"})
	}
	hydrateQuestions(session.Questions)
	return c.JSON(fiber.Map{"session": session})
}

// Complete awards the session's XP to the user (generous, demo-friendly loop).
func (l *ListeningController) Complete(c *fiber.Ctx) error {
	user := middleware.CurrentUser(c)
	rollOverDay(user)

	id := c.Params("id")
	var session models.ListeningSession
	if err := database.DB.First(&session, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "session not found"})
	}

	xpGain := session.XPReward
	if xpGain <= 0 {
		xpGain = 15
	}

	user.XP += xpGain
	user.XPToday += xpGain
	user.Gems += 3

	touchStreak(user)

	promoteLevel(user)
	database.DB.Save(user)

	// Listening sessions sit at unit level, so they weight a little above a
	// single lesson. There's no per-question accuracy here, so it scores clean.
	points := AwardLeaguePoints(user, LeagueAward{
		Source: "listening", RawXP: xpGain, Accuracy: 100, Difficulty: 1.3,
	})

	return c.JSON(fiber.Map{"xpEarned": xpGain, "leaguePoints": points, "user": user})
}
