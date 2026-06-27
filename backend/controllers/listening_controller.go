package controllers

import (
	"encoding/json"
	"time"

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
		Preload("Lines", orderByIndex).
		Preload("Questions", orderByIndex).
		Order("order_index asc")

	q.Where("language = ?", lang).Find(&sessions)
	if len(sessions) == 0 && lang != "es" {
		q.Where("language = ?", "es").Find(&sessions)
	}

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

	today := time.Now().Format("2006-01-02")
	if user.LastActiveDate != today {
		user.Streak++
		user.LastActiveDate = today
		user.XPToday = xpGain
	}

	promoteLevel(user)
	database.DB.Save(user)

	return c.JSON(fiber.Map{"xpEarned": xpGain, "user": user})
}
