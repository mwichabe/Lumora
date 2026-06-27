package controllers

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"lumora/backend/database"
	"lumora/backend/middleware"
	"lumora/backend/models"
)

// LessonController serves skill-tree and lesson content.
type LessonController struct{}

// skillNode is the galaxy-map representation returned to the frontend, enriched
// with per-user unlock and completion state.
type skillNode struct {
	models.Skill
	Unlocked       bool `json:"unlocked"`
	Completed      bool `json:"completed"`
	LessonCount    int  `json:"lessonCount"`
	CompletedCount int  `json:"completedCount"`
}

// GalaxyMap returns all skills for the user's target language with unlock state.
func (l *LessonController) GalaxyMap(c *fiber.Ctx) error {
	user := middleware.CurrentUser(c)
	lang := user.TargetLanguage
	if lang == "" {
		lang = "es"
	}

	var skills []models.Skill
	database.DB.Where("language = ?", lang).
		Preload("Lessons").
		Order("order_index asc").
		Find(&skills)

	// Fall back to the fully-authored Spanish course for languages that are
	// still previews in the MVP, so the learner always has content.
	if len(skills) == 0 && lang != "es" {
		database.DB.Where("language = ?", "es").
			Preload("Lessons").
			Order("order_index asc").
			Find(&skills)
	}

	// Gather the set of lessons this user has completed.
	var done []models.LessonProgress
	database.DB.Where("user_id = ? AND completed = ?", user.ID, true).Find(&done)
	completedLessons := map[uint]bool{}
	for _, d := range done {
		completedLessons[d.LessonID] = true
	}

	nodes := make([]skillNode, 0, len(skills))
	for _, s := range skills {
		completed := 0
		for _, ls := range s.Lessons {
			if completedLessons[ls.ID] {
				completed++
			}
		}
		nodes = append(nodes, skillNode{
			Skill:          s,
			Unlocked:       user.XP >= s.RequiredXP,
			Completed:      len(s.Lessons) > 0 && completed == len(s.Lessons),
			LessonCount:    len(s.Lessons),
			CompletedCount: completed,
		})
	}

	return c.JSON(fiber.Map{"skills": nodes})
}

// GetLesson returns a lesson with hydrated exercises (options decoded).
func (l *LessonController) GetLesson(c *fiber.Ctx) error {
	id := c.Params("id")

	var lesson models.Lesson
	if err := database.DB.
		Preload("Vocab", func(db *gorm.DB) *gorm.DB {
			return db.Order("order_index asc")
		}).
		Preload("Exercises", func(db *gorm.DB) *gorm.DB {
			return db.Order("order_index asc")
		}).First(&lesson, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "lesson not found"})
	}

	// Order + hydrate options from the stored JSON column.
	for i := range lesson.Exercises {
		var o []string
		if lesson.Exercises[i].OptionsJSON != "" {
			_ = json.Unmarshal([]byte(lesson.Exercises[i].OptionsJSON), &o)
		}
		lesson.Exercises[i].Options = o
	}

	return c.JSON(fiber.Map{"lesson": lesson})
}
