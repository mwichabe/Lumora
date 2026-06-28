package controllers

import (
	"encoding/json"
	"math/rand"

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

	// The learner can't type the target language yet, so turn typed exercises
	// (translate / fill) into multiple choice by generating plausible options.
	addChoiceOptions(&lesson)

	return c.JSON(fiber.Map{"lesson": lesson})
}

// addChoiceOptions gives translate/fill exercises a set of multiple-choice
// options (correct answer + distractors drawn from the lesson, then padded).
func addChoiceOptions(lesson *models.Lesson) {
	var phrasePool, wordPool []string
	for _, e := range lesson.Exercises {
		switch e.Type {
		case models.ExerciseTranslate:
			phrasePool = append(phrasePool, e.CorrectAnswer)
		case models.ExerciseFill:
			wordPool = append(wordPool, e.CorrectAnswer)
		}
	}
	for _, v := range lesson.Vocab {
		wordPool = append(wordPool, v.Word)
		phrasePool = append(phrasePool, v.Word)
	}

	phraseFiller := []string{"Buenos días", "Por favor", "Hasta luego", "No lo sé", "Mucho gusto"}
	wordFiller := []string{"gracias", "hola", "casa", "agua", "bien", "sí"}

	for i := range lesson.Exercises {
		e := &lesson.Exercises[i]
		if len(e.Options) > 0 {
			continue
		}
		if e.Type != models.ExerciseTranslate && e.Type != models.ExerciseFill {
			continue
		}
		pool, filler := phrasePool, phraseFiller
		if e.Type == models.ExerciseFill {
			pool, filler = wordPool, wordFiller
		}
		e.Options = buildOptions(e.CorrectAnswer, pool, filler, 4)
	}
}

func buildOptions(correct string, pool, filler []string, n int) []string {
	seen := map[string]bool{correct: true}
	out := []string{correct}

	take := func(src []string) {
		for _, s := range src {
			if len(out) >= n {
				return
			}
			if s == "" || seen[s] {
				continue
			}
			seen[s] = true
			out = append(out, s)
		}
	}
	take(pool)
	take(filler)

	// Shuffle so the answer isn't always first.
	rand.Shuffle(len(out), func(i, j int) { out[i], out[j] = out[j], out[i] })
	return out
}
