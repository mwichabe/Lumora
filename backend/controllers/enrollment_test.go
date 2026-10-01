package controllers

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"

	"lumora/backend/database"
	"lumora/backend/models"
)

func newEnrollmentFixture(t *testing.T, active string, langs ...string) (*fiber.App, *models.User) {
	t.Helper()
	db := newTestDB(t)
	extra := []interface{}{&models.Enrollment{}, &models.Skill{}, &models.Lesson{}, &models.LessonProgress{}, &models.Certificate{}}
	if err := db.Migrator().DropTable(extra...); err != nil {
		t.Fatalf("drop: %v", err)
	}
	if err := db.AutoMigrate(extra...); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	user := &models.User{Email: "learner@example.com", TargetLanguage: active, CEFRLevel: "A1"}
	database.DB.Create(user)
	for _, l := range langs {
		database.DB.Create(&models.Enrollment{UserID: user.ID, Language: l})
	}

	ec := &EnrollmentController{}
	app := fiber.New()
	app.Delete("/enrollments/:language", func(c *fiber.Ctx) error {
		c.Locals("user", user)
		return c.Next()
	}, ec.Remove)
	return app, user
}

type removeResponse struct {
	Languages []string `json:"languages"`
	Active    string   `json:"active"`
	Error     string   `json:"error"`
}

func remove(t *testing.T, app *fiber.App, lang string) (int, removeResponse) {
	t.Helper()
	res, err := app.Test(httptest.NewRequest("DELETE", "/enrollments/"+lang, nil))
	if err != nil {
		t.Fatal(err)
	}
	var body removeResponse
	_ = json.NewDecoder(res.Body).Decode(&body)
	return res.StatusCode, body
}

func TestRemoveInactiveLanguage(t *testing.T) {
	app, user := newEnrollmentFixture(t, "de", "es", "de", "fr")

	status, body := remove(t, app, "fr")
	if status != fiber.StatusOK {
		t.Fatalf("status %d: %s", status, body.Error)
	}
	if body.Active != "de" || user.TargetLanguage != "de" {
		t.Errorf("active course changed to %q; removing another language must not switch it", body.Active)
	}
	if len(body.Languages) != 2 || body.Languages[0] != "es" || body.Languages[1] != "de" {
		t.Errorf("languages = %v, want [es de]", body.Languages)
	}
}

func TestRemoveActiveLanguageSwitchesCourse(t *testing.T) {
	app, user := newEnrollmentFixture(t, "de", "es", "de", "fr")

	status, body := remove(t, app, "de")
	if status != fiber.StatusOK {
		t.Fatalf("status %d: %s", status, body.Error)
	}
	// The most recently added remaining course takes over.
	if body.Active != "fr" || user.TargetLanguage != "fr" {
		t.Errorf("active = %q, want fr", body.Active)
	}
	var saved models.User
	database.DB.First(&saved, user.ID)
	if saved.TargetLanguage != "fr" {
		t.Errorf("new active course not saved: %q", saved.TargetLanguage)
	}
}

func TestCannotRemoveLastLanguage(t *testing.T) {
	app, _ := newEnrollmentFixture(t, "es", "es")

	status, body := remove(t, app, "es")
	if status != fiber.StatusBadRequest {
		t.Fatalf("status %d, want 400", status)
	}
	if body.Error == "" {
		t.Error("expected a message explaining why")
	}
	var count int64
	database.DB.Model(&models.Enrollment{}).Count(&count)
	if count != 1 {
		t.Errorf("the only course was removed anyway (%d left)", count)
	}
}

func TestRemoveUnknownLanguage(t *testing.T) {
	app, _ := newEnrollmentFixture(t, "es", "es", "de")
	if status, _ := remove(t, app, "ja"); status != fiber.StatusNotFound {
		t.Errorf("status %d, want 404", status)
	}
}
