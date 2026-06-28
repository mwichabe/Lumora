package controllers

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	"lumora/backend/database"
	"lumora/backend/middleware"
	"lumora/backend/models"
)

// EnrollmentController manages the set of languages a user is learning and which
// one is currently active.
type EnrollmentController struct{}

// EnsureEnrollment makes sure a (user, language) enrollment row exists.
func EnsureEnrollment(userID uint, lang string) {
	if lang == "" {
		return
	}
	var e models.Enrollment
	err := database.DB.Where("user_id = ? AND language = ?", userID, lang).First(&e).Error
	if err != nil {
		database.DB.Create(&models.Enrollment{UserID: userID, Language: lang})
	}
}

func (ec *EnrollmentController) list(userID uint, active string) []string {
	var rows []models.Enrollment
	database.DB.Where("user_id = ?", userID).Order("created_at asc").Find(&rows)
	langs := make([]string, 0, len(rows))
	for _, r := range rows {
		langs = append(langs, r.Language)
	}
	// Make sure the active language is always present.
	if active != "" {
		found := false
		for _, l := range langs {
			if l == active {
				found = true
				break
			}
		}
		if !found {
			langs = append(langs, active)
		}
	}
	return langs
}

// List returns the user's enrolled languages and the active one.
func (ec *EnrollmentController) List(c *fiber.Ctx) error {
	user := middleware.CurrentUser(c)
	if user.TargetLanguage != "" {
		EnsureEnrollment(user.ID, user.TargetLanguage)
	}
	return c.JSON(fiber.Map{
		"languages": ec.list(user.ID, user.TargetLanguage),
		"active":    user.TargetLanguage,
	})
}

type enrollInput struct {
	Language string `json:"language"`
}

// Enroll adds a language and makes it the active course.
func (ec *EnrollmentController) Enroll(c *fiber.Ctx) error {
	user := middleware.CurrentUser(c)
	var in enrollInput
	if err := c.BodyParser(&in); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
	}
	in.Language = strings.ToLower(strings.TrimSpace(in.Language))
	if in.Language == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "language is required"})
	}

	EnsureEnrollment(user.ID, in.Language)
	user.TargetLanguage = in.Language
	database.DB.Save(user)

	return c.JSON(fiber.Map{
		"languages": ec.list(user.ID, user.TargetLanguage),
		"active":    user.TargetLanguage,
		"user":      user,
	})
}

// SetActive switches the active course to an already-enrolled language.
func (ec *EnrollmentController) SetActive(c *fiber.Ctx) error {
	user := middleware.CurrentUser(c)
	var in enrollInput
	if err := c.BodyParser(&in); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
	}
	in.Language = strings.ToLower(strings.TrimSpace(in.Language))

	// Must already be enrolled.
	var e models.Enrollment
	if err := database.DB.Where("user_id = ? AND language = ?", user.ID, in.Language).First(&e).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "not enrolled in that language"})
	}

	user.TargetLanguage = in.Language
	database.DB.Save(user)

	return c.JSON(fiber.Map{
		"languages": ec.list(user.ID, user.TargetLanguage),
		"active":    user.TargetLanguage,
		"user":      user,
	})
}
