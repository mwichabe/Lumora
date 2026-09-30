package controllers

import (
	"testing"

	"lumora/backend/database"
	"lumora/backend/models"
)

// levelFixture is a German course with two A1 lessons and one A2 lesson, plus
// an untagged Spanish course, on a fresh database.
type levelFixture struct {
	a1a, a1b, a2, es models.Lesson
}

func newLevelFixture(t *testing.T) levelFixture {
	t.Helper()
	db := newTestDB(t)
	extra := []interface{}{&models.Skill{}, &models.Lesson{}, &models.LessonProgress{}, &models.Certificate{}}
	if err := db.Migrator().DropTable(extra...); err != nil {
		t.Fatalf("drop: %v", err)
	}
	if err := db.AutoMigrate(extra...); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	lesson := func(lang, unit string) models.Lesson {
		s := models.Skill{Language: lang, Unit: unit}
		db.Create(&s)
		l := models.Lesson{SkillID: s.ID}
		db.Create(&l)
		return l
	}
	return levelFixture{
		a1a: lesson("de", "A1 · Grundlagen"),
		a1b: lesson("de", "A1 · Grundlagen"),
		a2:  lesson("de", "A2 · Aufbaustufe"),
		es:  lesson("es", "Basics"),
	}
}

func complete(t *testing.T, userID uint, lessons ...models.Lesson) {
	t.Helper()
	for _, l := range lessons {
		database.DB.Create(&models.LessonProgress{UserID: userID, LessonID: l.ID, Completed: true})
	}
}

func TestEarnedLevelIgnoresXP(t *testing.T) {
	newLevelFixture(t)
	// The bug: 400+ XP used to mean C1. XP is not an input at all now.
	u := models.User{Email: "xp@example.com", TargetLanguage: "de", CEFRLevel: "A1", XP: 5000}
	database.DB.Create(&u)
	if got := earnedLevel(u.ID, "de"); got != "A1" {
		t.Errorf("a learner with 5000 XP and no completed level is %s, want A1", got)
	}
}

func TestEarnedLevelFollowsCourseCompletion(t *testing.T) {
	f := newLevelFixture(t)
	const user = 7

	complete(t, user, f.a1a)
	if got := earnedLevel(user, "de"); got != "A1" {
		t.Errorf("with A1 half done: %s, want A1", got)
	}

	complete(t, user, f.a1b)
	if got := earnedLevel(user, "de"); got != "A2" {
		t.Errorf("with every A1 lesson done: %s, want A2", got)
	}

	// A2 is the last level this course has content for; finishing it moves the
	// learner on to B1, and no further since there is nothing there to finish.
	complete(t, user, f.a2)
	if got := earnedLevel(user, "de"); got != "B1" {
		t.Errorf("with A1 and A2 done: %s, want B1", got)
	}
}

func TestEarnedLevelIsPerLanguage(t *testing.T) {
	f := newLevelFixture(t)
	const user = 7
	complete(t, user, f.a1a, f.a1b, f.es)

	if got := earnedLevel(user, "de"); got != "A2" {
		t.Errorf("German: %s, want A2", got)
	}
	// Spanish units aren't tagged with a level, so finishing them proves nothing
	// beyond beginner — and German progress must not leak across.
	if got := earnedLevel(user, "es"); got != "A1" {
		t.Errorf("Spanish: %s, want A1", got)
	}
}

func TestEarnedLevelCountsCertificates(t *testing.T) {
	newLevelFixture(t)
	const user = 7
	database.DB.Create(&models.Certificate{UserID: user, Language: "es", Level: "B1"})

	if got := earnedLevel(user, "es"); got != "B1" {
		t.Errorf("with a B1 certificate: %s, want B1", got)
	}
	if got := earnedLevel(user, "de"); got != "A1" {
		t.Errorf("a Spanish certificate must not count for German: got %s", got)
	}
}

func TestSyncLevelReportsOnlyRealAdvances(t *testing.T) {
	f := newLevelFixture(t)
	u := models.User{Email: "sync@example.com", TargetLanguage: "de", CEFRLevel: "A1", LevelName: "Spark"}
	database.DB.Create(&u)

	if syncLevel(&u) {
		t.Error("nothing completed yet: must not report a level-up")
	}
	complete(t, u.ID, f.a1a, f.a1b)
	if !syncLevel(&u) || u.CEFRLevel != "A2" || u.LevelName != "Glow" {
		t.Errorf("after finishing A1: advanced to %s (%s), want A2 (Glow)", u.CEFRLevel, u.LevelName)
	}
	if syncLevel(&u) {
		t.Error("a second sync with no new progress must not report a level-up again")
	}
}

func TestRepairLevelsUndoesXPBasedLevels(t *testing.T) {
	newLevelFixture(t)
	db := database.DB
	// A learner the old rule promoted to C1 on XP alone, and told so.
	u := models.User{Email: "inflated@example.com", TargetLanguage: "es", CEFRLevel: "C1", LevelName: "Aurora", XP: 450}
	db.Create(&u)
	db.Create(&models.Notification{UserID: u.ID, Key: "levelup_C1", Title: "Level up! You reached C1"})
	db.Create(&models.Notification{UserID: u.ID, Key: "levelup_de_A2", Title: "You've moved up to A2!"})
	db.Create(&models.Notification{UserID: u.ID, Key: "welcome", Title: "Welcome to Lumora!"})

	RepairLevels()

	var after models.User
	db.First(&after, u.ID)
	if after.CEFRLevel != "A1" || after.LevelName != "Spark" {
		t.Errorf("level after repair: %s (%s), want A1 (Spark)", after.CEFRLevel, after.LevelName)
	}
	if after.XP != 450 {
		t.Errorf("repair must not touch XP: got %d", after.XP)
	}

	var keys []string
	db.Model(&models.Notification{}).Where("user_id = ?", u.ID).Order("key").Pluck("key", &keys)
	if len(keys) != 2 || keys[0] != "levelup_de_A2" || keys[1] != "welcome" {
		t.Errorf("only the XP-based level-up should be withdrawn; remaining keys: %v", keys)
	}
}

func TestUnitLevel(t *testing.T) {
	cases := map[string]string{
		"A1 · Grundlagen": "A1",
		"C2 · Feinheiten": "C2",
		"B1":              "B1",
		"Basics":          "",
		"A1rport phrases": "",
		"Everyday Life":   "",
		"":                "",
	}
	for unit, want := range cases {
		if got := unitLevel(unit); got != want {
			t.Errorf("unitLevel(%q) = %q, want %q", unit, got, want)
		}
	}
}
