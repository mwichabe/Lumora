package controllers

import (
	"github.com/gofiber/fiber/v2"

	"lumora/backend/database"
	"lumora/backend/middleware"
	"lumora/backend/models"
)

// PracticeController powers the Practice tab: it serves a pool of vocabulary to
// build drills from, tracks mistakes, and awards XP for finished sessions.
type PracticeController struct{}

// Pool returns the vocabulary for the active language plus the user's open
// mistakes — the frontend turns these into quiz/listening/speaking drills.
func (p *PracticeController) Pool(c *fiber.Ctx) error {
	user := middleware.CurrentUser(c)
	lang := user.TargetLanguage
	if lang == "" {
		lang = "es"
	}

	// Only words from skills the learner has unlocked — so the pool grows
	// automatically as they progress, and never quizzes locked content.
	var vocab []models.VocabItem
	database.DB.
		Joins("JOIN lessons ON lessons.id = vocab_items.lesson_id").
		Joins("JOIN skills ON skills.id = lessons.skill_id").
		Where("skills.language = ? AND skills.required_xp <= ?", lang, user.XP).
		Order("skills.order_index desc, vocab_items.id asc"). // newest skills first
		Find(&vocab)

	// Fallback: a brand-new learner (no unlocked skills yet) still gets the
	// very first words so Practice isn't empty.
	if len(vocab) == 0 {
		database.DB.
			Joins("JOIN lessons ON lessons.id = vocab_items.lesson_id").
			Joins("JOIN skills ON skills.id = lessons.skill_id").
			Where("skills.language = ?", lang).
			Order("skills.order_index asc, vocab_items.id asc").
			Limit(8).
			Find(&vocab)
	}

	var mistakes []models.Mistake
	database.DB.Where("user_id = ? AND language = ?", user.ID, lang).
		Order("created_at desc").Find(&mistakes)

	// Counts of comprehension practice the learner has unlocked, so the Practice
	// tab can show/enable the Listening and Reading modes.
	unlocked := unitUnlocker(lang, user.XP)
	listeningCount := countUnlockedSessions("listening_sessions", lang, unlocked)
	readingCount := countUnlockedSessions("reading_sessions", lang, unlocked)

	return c.JSON(fiber.Map{
		"vocab":          vocab,
		"mistakes":       mistakes,
		"listeningCount": listeningCount,
		"readingCount":   readingCount,
	})
}

// unitUnlocker returns a predicate telling whether a unit has been reached. A
// unit is reached once the user's XP meets the lowest skill gate in that unit,
// matching the lock logic on the Learn tab — so practice only ever serves
// content the learner has actually unlocked.
func unitUnlocker(lang string, xp int) func(string) bool {
	var skills []models.Skill
	database.DB.Where("language = ?", lang).Find(&skills)
	minXP := map[string]int{}
	for _, s := range skills {
		if v, ok := minXP[s.Unit]; !ok || s.RequiredXP < v {
			minXP[s.Unit] = s.RequiredXP
		}
	}
	return func(unit string) bool {
		req, ok := minXP[unit]
		if !ok {
			return true // no gate info → treat as available
		}
		return xp >= req
	}
}

// countUnlockedSessions counts listening/reading sessions in unlocked units.
func countUnlockedSessions(table, lang string, unlocked func(string) bool) int {
	type unitRow struct{ Unit string }
	var rows []unitRow
	database.DB.Table(table).Select("unit").Where("language = ?", lang).Find(&rows)
	n := 0
	for _, r := range rows {
		if unlocked(r.Unit) {
			n++
		}
	}
	return n
}

// Listening returns the unlocked listening sessions for comprehension practice.
func (p *PracticeController) Listening(c *fiber.Ctx) error {
	user := middleware.CurrentUser(c)
	lang := user.TargetLanguage
	if lang == "" {
		lang = "es"
	}
	unlocked := unitUnlocker(lang, user.XP)

	var sessions []models.ListeningSession
	database.DB.
		Preload("Lines", orderByIndex).
		Preload("Questions", orderByIndex).
		Where("language = ?", lang).Order("order_index asc").Find(&sessions)

	out := make([]models.ListeningSession, 0, len(sessions))
	for i := range sessions {
		if unlocked(sessions[i].Unit) {
			hydrateQuestions(sessions[i].Questions)
			out = append(out, sessions[i])
		}
	}
	return c.JSON(fiber.Map{"sessions": out})
}

// Reading returns the unlocked reading sessions for comprehension practice.
func (p *PracticeController) Reading(c *fiber.Ctx) error {
	user := middleware.CurrentUser(c)
	lang := user.TargetLanguage
	if lang == "" {
		lang = "es"
	}
	unlocked := unitUnlocker(lang, user.XP)

	var sessions []models.ReadingSession
	database.DB.
		Preload("Lines", orderByIndex).
		Preload("Questions", orderByIndex).
		Where("language = ?", lang).Order("order_index asc").Find(&sessions)

	out := make([]models.ReadingSession, 0, len(sessions))
	for i := range sessions {
		if unlocked(sessions[i].Unit) {
			hydrateReadingQuestions(sessions[i].Questions)
			out = append(out, sessions[i])
		}
	}
	return c.JSON(fiber.Map{"sessions": out})
}

type mistakeInput struct {
	Prompt        string `json:"prompt"`
	Question      string `json:"question"`
	CorrectAnswer string `json:"correctAnswer"`
}

// RecordMistake stores a missed exercise (deduped per user+language+question).
func (p *PracticeController) RecordMistake(c *fiber.Ctx) error {
	user := middleware.CurrentUser(c)
	lang := user.TargetLanguage
	if lang == "" {
		lang = "es"
	}
	var in mistakeInput
	if err := c.BodyParser(&in); err != nil || in.Question == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
	}

	var existing models.Mistake
	err := database.DB.Where(
		"user_id = ? AND language = ? AND question = ? AND correct_answer = ?",
		user.ID, lang, in.Question, in.CorrectAnswer,
	).First(&existing).Error
	if err != nil {
		database.DB.Create(&models.Mistake{
			UserID: user.ID, Language: lang,
			Prompt: in.Prompt, Question: in.Question, CorrectAnswer: in.CorrectAnswer,
		})
	}
	return c.JSON(fiber.Map{"ok": true})
}

type resolveInput struct {
	IDs []uint `json:"ids"`
}

// ResolveMistakes removes mistakes the user has now answered correctly.
func (p *PracticeController) ResolveMistakes(c *fiber.Ctx) error {
	user := middleware.CurrentUser(c)
	var in resolveInput
	if err := c.BodyParser(&in); err != nil || len(in.IDs) == 0 {
		return c.JSON(fiber.Map{"ok": true})
	}
	database.DB.Where("user_id = ? AND id IN ?", user.ID, in.IDs).Delete(&models.Mistake{})
	return c.JSON(fiber.Map{"ok": true})
}

type practiceCompleteInput struct {
	XP int `json:"xp"`
}

// Complete awards XP for a finished practice session.
func (p *PracticeController) Complete(c *fiber.Ctx) error {
	user := middleware.CurrentUser(c)
	rollOverDay(user)

	var in practiceCompleteInput
	_ = c.BodyParser(&in)
	xpGain := in.XP
	if xpGain <= 0 {
		xpGain = 10
	}
	if xpGain > 100 {
		xpGain = 100 // safety cap
	}

	user.XP += xpGain
	user.XPToday += xpGain
	user.Gems += 2

	touchStreak(user)

	promoteLevel(user)
	database.DB.Save(user)

	return c.JSON(fiber.Map{"xpEarned": xpGain, "user": user})
}
