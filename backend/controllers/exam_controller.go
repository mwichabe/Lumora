package controllers

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"

	"lumora/backend/database"
	"lumora/backend/middleware"
	"lumora/backend/models"
)

// ExamController scores the proficiency exam and issues certificates. The exam
// and certificate are paid features; payment is bypassed for now.
type ExamController struct{}

type examSubmit struct {
	Language  string `json:"language"`
	Level     string `json:"level"` // the level the user chose to attempt
	Listening int    `json:"listening"`
	Reading   int    `json:"reading"`
	Writing   int    `json:"writing"`
	Speaking  int    `json:"speaking"`
}

// Each level requires a higher overall score to pass — harder as you climb.
// FINAL is the comprehensive A1→C2 mastery exam.
var levelPassMark = map[string]int{
	"A1": 50, "A2": 58, "B1": 65, "B2": 72, "C1": 80, "C2": 88, "FINAL": 80,
}

// Section weights (sum to 100). Listening & Reading carry the most weight; the
// productive skills (Writing, Speaking) are weighted but slightly lighter.
var sectionWeights = map[string]int{
	"listening": 30, "reading": 30, "writing": 20, "speaking": 20,
}

// Per-level time limit (seconds). Higher levels get a longer but more demanding
// paper — combined with a tougher pass mark, the exam gets harder as you climb.
var levelDuration = map[string]int{
	"A1": 600, "A2": 780, "B1": 960, "B2": 1140, "C1": 1320, "C2": 1500,
	"FINAL": 2100, // 35 minutes — the comprehensive mastery paper
}

func passMarkForLevel(level string) int {
	if v, ok := levelPassMark[level]; ok {
		return v
	}
	return 50
}

func durationForLevel(level string) int {
	if v, ok := levelDuration[level]; ok {
		return v
	}
	return 420
}

// weightedOverall combines the four section scores by their weights.
func weightedOverall(l, r, w, s int) int {
	total := l*sectionWeights["listening"] + r*sectionWeights["reading"] +
		w*sectionWeights["writing"] + s*sectionWeights["speaking"]
	return clamp100(total / 100)
}

func clamp100(v int) int {
	if v < 0 {
		return 0
	}
	if v > 100 {
		return 100
	}
	return v
}

// levelFromScore maps an overall exam score to a CEFR level.
func levelFromScore(s int) string {
	switch {
	case s >= 96:
		return "C2"
	case s >= 90:
		return "C1"
	case s >= 80:
		return "B2"
	case s >= 70:
		return "B1"
	case s >= 55:
		return "A2"
	default:
		return "A1"
	}
}

// Submit scores the four sections, and (on a pass) issues a certificate. Each
// attempt consumes one paid attempt for that level — retaking requires paying
// again.
func (ec *ExamController) Submit(c *fiber.Ctx) error {
	user := middleware.CurrentUser(c)

	var in examSubmit
	if err := c.BodyParser(&in); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
	}
	lang := in.Language
	if lang == "" {
		lang = user.TargetLanguage
	}
	level := in.Level
	if _, ok := levelPassMark[level]; !ok {
		level = "A1"
	}

	// Monetization gate: when payments are enabled, each attempt needs a paid,
	// unconsumed token for this level — consume it now (retakes must pay again).
	if PaymentsEnabled() {
		if !consumePaidAttempt(user.ID, level) {
			return c.Status(fiber.StatusPaymentRequired).
				JSON(fiber.Map{"error": "payment required", "code": "payment_required"})
		}
	}

	l := clamp100(in.Listening)
	r := clamp100(in.Reading)
	w := clamp100(in.Writing)
	s := clamp100(in.Speaking)
	overall := weightedOverall(l, r, w, s)

	passMark := passMarkForLevel(level)
	passed := overall >= passMark

	resp := fiber.Map{
		"passed":       passed,
		"alreadyTaken": false,
		"overall":      overall,
		"level":        level,
		"passMark":     passMark,
		"weights":      sectionWeights,
		"sections":     fiber.Map{"listening": l, "reading": r, "writing": w, "speaking": s},
	}

	if passed {
		// Upsert: a retake refreshes the certificate (keeps its serial).
		var cert models.Certificate
		exists := database.DB.
			Where("user_id = ? AND language = ? AND level = ?", user.ID, lang, level).
			First(&cert).Error == nil

		cert.UserID = user.ID
		cert.UserName = displayName(*user)
		cert.Language = lang
		cert.Level = level
		cert.Score = overall
		cert.Listening, cert.Reading, cert.Writing, cert.Speaking = l, r, w, s
		cert.IssuedAt = time.Now()

		if exists {
			database.DB.Save(&cert)
		} else {
			database.DB.Create(&cert)
			cert.Serial = fmt.Sprintf("LUM-%s-%04d", time.Now().Format("2006"), cert.ID)
			database.DB.Save(&cert)
		}
		resp["certificate"] = cert
		DeliverExamPassed(user.ID, level, overall, fmt.Sprintf("/certificates/%d", cert.ID))
	} else {
		DeliverExamFailed(user.ID, level, overall, passMark)
	}

	return c.JSON(resp)
}

// examLangDisplay maps a language code to a friendly name for notifications.
var examLangDisplay = map[string]string{
	"es": "Spanish", "de": "German", "fr": "French", "it": "Italian",
	"pt": "Portuguese", "en": "English",
}

func langDisplay(code string) string {
	if n, ok := examLangDisplay[code]; ok {
		return n
	}
	return code
}

// Start notifies the user that they've begun an exam attempt.
func (ec *ExamController) Start(c *fiber.Ctx) error {
	user := middleware.CurrentUser(c)
	var in struct {
		Level    string `json:"level"`
		Language string `json:"language"`
	}
	_ = c.BodyParser(&in)
	level := in.Level
	if _, ok := levelPassMark[level]; !ok {
		level = "A1"
	}
	lang := in.Language
	if lang == "" {
		lang = user.TargetLanguage
	}
	DeliverExamStarted(user.ID, level, langDisplay(lang))
	return c.JSON(fiber.Map{"ok": true})
}

// Meta describes the exam rules the frontend shows before the user starts:
// section weights, per-level pass marks and time limits.
func (ec *ExamController) Meta(c *fiber.Ctx) error {
	levels := fiber.Map{}
	for code, pm := range levelPassMark {
		levels[code] = fiber.Map{
			"passMark":        pm,
			"durationSeconds": durationForLevel(code),
		}
	}
	return c.JSON(fiber.Map{
		"weights": sectionWeights,
		"levels":  levels,
	})
}

// Verify is a PUBLIC endpoint: anyone with a certificate serial can confirm it
// is genuine. Returns only the holder's name and result — no private data.
func (ec *ExamController) Verify(c *fiber.Ctx) error {
	serial := c.Params("serial")
	var cert models.Certificate
	if err := database.DB.Where("serial = ?", serial).First(&cert).Error; err != nil {
		return c.JSON(fiber.Map{"valid": false})
	}
	return c.JSON(fiber.Map{
		"valid": true,
		"certificate": fiber.Map{
			"userName": cert.UserName,
			"language": cert.Language,
			"level":    cert.Level,
			"score":    cert.Score,
			"serial":   cert.Serial,
			"issuedAt": cert.IssuedAt,
		},
	})
}

// ListCertificates returns the user's earned certificates (newest first).
func (ec *ExamController) ListCertificates(c *fiber.Ctx) error {
	user := middleware.CurrentUser(c)
	var certs []models.Certificate
	database.DB.Where("user_id = ?", user.ID).Order("issued_at desc").Find(&certs)
	return c.JSON(fiber.Map{"certificates": certs})
}

// DeleteCertificate removes a certificate the user owns. Because each level can
// only be passed once, deleting a certificate also frees that level to retake.
func (ec *ExamController) DeleteCertificate(c *fiber.Ctx) error {
	user := middleware.CurrentUser(c)
	id := c.Params("id")
	var cert models.Certificate
	if err := database.DB.Where("id = ? AND user_id = ?", id, user.ID).First(&cert).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "certificate not found"})
	}
	database.DB.Delete(&cert)
	return c.JSON(fiber.Map{"ok": true})
}

// GetCertificate returns a single certificate the user owns.
func (ec *ExamController) GetCertificate(c *fiber.Ctx) error {
	user := middleware.CurrentUser(c)
	id := c.Params("id")
	var cert models.Certificate
	if err := database.DB.Where("id = ? AND user_id = ?", id, user.ID).First(&cert).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "certificate not found"})
	}
	return c.JSON(fiber.Map{"certificate": cert})
}
