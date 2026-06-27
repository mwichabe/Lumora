package controllers

import (
	"time"

	"github.com/gofiber/fiber/v2"

	"lumora/backend/database"
	"lumora/backend/middleware"
	"lumora/backend/models"
)

// QuestController exposes daily quests.
type QuestController struct{}

// Daily returns today's quests for the user, creating them if needed.
func (q *QuestController) Daily(c *fiber.Ctx) error {
	user := middleware.CurrentUser(c)
	return c.JSON(fiber.Map{"quests": ensureDailyQuests(user.ID)})
}

// ensureDailyQuests materialises a UserQuest row per template quest for today.
func ensureDailyQuests(userID uint) []models.UserQuest {
	today := time.Now().Format("2006-01-02")

	var existing []models.UserQuest
	database.DB.Where("user_id = ? AND date = ?", userID, today).
		Preload("Quest").Find(&existing)
	if len(existing) > 0 {
		return existing
	}

	var templates []models.Quest
	database.DB.Find(&templates)
	for _, t := range templates {
		uq := models.UserQuest{UserID: userID, QuestID: t.ID, Date: today}
		database.DB.Create(&uq)
	}

	database.DB.Where("user_id = ? AND date = ?", userID, today).
		Preload("Quest").Find(&existing)
	return existing
}

// updateQuestsOnLesson advances quest progress after a lesson completion.
func updateQuestsOnLesson(userID uint, accuracy int) {
	uqs := ensureDailyQuests(userID)
	for i := range uqs {
		uq := &uqs[i]
		if uq.Completed || uq.Quest == nil {
			continue
		}
		// Simple heuristic: a completed lesson counts toward every quest type.
		uq.Progress++
		if uq.Progress >= uq.Quest.Target {
			uq.Completed = true
			uq.Progress = uq.Quest.Target
		}
		database.DB.Save(uq)
	}
}
