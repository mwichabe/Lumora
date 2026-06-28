package controllers

import (
	"encoding/json"
	"time"

	"github.com/gofiber/fiber/v2"

	"lumora/backend/database"
	"lumora/backend/middleware"
	"lumora/backend/models"
)

// ReadingController serves unit-level reading passages + comprehension questions.
type ReadingController struct{}

func hydrateReadingQs(qs []models.ReadingQuestion) {
	for i := range qs {
		var o []string
		if qs[i].OptionsJSON != "" {
			_ = json.Unmarshal([]byte(qs[i].OptionsJSON), &o)
		}
		qs[i].Options = o
	}
}

func (r *ReadingController) List(c *fiber.Ctx) error {
	user := middleware.CurrentUser(c)
	lang := user.TargetLanguage
	if lang == "" {
		lang = "es"
	}

	var sessions []models.ReadingSession
	q := database.DB.
		Preload("Lines", orderByIndex).
		Preload("Questions", orderByIndex).
		Order("order_index asc")

	q.Where("language = ?", lang).Find(&sessions)
	if len(sessions) == 0 && lang != "es" {
		q.Where("language = ?", "es").Find(&sessions)
	}
	for i := range sessions {
		hydrateReadingQs(sessions[i].Questions)
	}
	return c.JSON(fiber.Map{"sessions": sessions})
}

func (r *ReadingController) Get(c *fiber.Ctx) error {
	id := c.Params("id")
	var session models.ReadingSession
	if err := database.DB.
		Preload("Lines", orderByIndex).
		Preload("Questions", orderByIndex).
		First(&session, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "session not found"})
	}
	hydrateReadingQs(session.Questions)
	return c.JSON(fiber.Map{"session": session})
}

func (r *ReadingController) Complete(c *fiber.Ctx) error {
	user := middleware.CurrentUser(c)
	rollOverDay(user)

	id := c.Params("id")
	var session models.ReadingSession
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
