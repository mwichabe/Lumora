package database

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"lumora/backend/models"
)

func seededItalian(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&models.Skill{}, &models.Lesson{}, &models.Exercise{}, &models.VocabItem{},
		&models.ListeningSession{}, &models.ListeningMatch{}, &models.ListeningLine{}, &models.ListeningQuestion{},
		&models.ReadingSession{}, &models.ReadingLine{}, &models.ReadingQuestion{}); err != nil {
		t.Fatal(err)
	}
	seedItalian(db)
	return db
}

func TestItalianCourseCoversA1ToC2(t *testing.T) {
	db := seededItalian(t)

	var skills []models.Skill
	db.Where("language = ?", "it").Preload("Lessons").Order("order_index").Find(&skills)

	seen := map[string]int{}
	lastXP := -1
	for _, s := range skills {
		seen[s.Unit[:2]]++
		if len(s.Lessons) == 0 {
			t.Errorf("skill %q has no lessons", s.Title)
		}
		if s.RequiredXP < lastXP {
			t.Errorf("skill %q needs %d XP, less than the skill before it (%d)", s.Title, s.RequiredXP, lastXP)
		}
		lastXP = s.RequiredXP
	}
	for _, level := range []string{"A1", "A2", "B1", "B2", "C1", "C2"} {
		if seen[level] < 5 {
			t.Errorf("level %s has %d skills, want at least 5", level, seen[level])
		}
	}

	// Sounds first; essere/avere and gender + articles early; the spoken past
	// (A2) before the subjunctive (B1); the literary passato remoto only at C1.
	if skills[0].Title != "Sounds & Spelling" || skills[0].RequiredXP != 0 {
		t.Errorf("first skill = %q (req %d XP), want Sounds & Spelling open from the start", skills[0].Title, skills[0].RequiredXP)
	}
	unitOf := map[string]string{}
	for _, s := range skills {
		unitOf[s.Title] = s.Unit[:2]
	}
	for title, level := range map[string]string{
		"Essere & Avere": "A1", "Gender & Articles": "A1",
		"Passato Prossimo with Avere": "A2", "Passato Prossimo with Essere": "A2", "The Imperfetto": "A2",
		"The Present Subjunctive": "B1", "Hypothetical Sentences": "B2",
		"Passato Remoto & Literary Italian": "C1", "Regional Varieties & Dialects": "C2",
	} {
		if got, ok := unitOf[title]; !ok {
			t.Errorf("missing skill %q", title)
		} else if got != level {
			t.Errorf("skill %q is at %s, want %s", title, got, level)
		}
	}
}

func TestItalianExercisesMixChoosingAndTyping(t *testing.T) {
	db := seededItalian(t)

	var exercises []models.Exercise
	db.Joins("JOIN lessons ON lessons.id = exercises.lesson_id").
		Joins("JOIN skills ON skills.id = lessons.skill_id").
		Where("skills.language = ?", "it").Find(&exercises)
	if len(exercises) < 250 {
		t.Fatalf("only %d exercises seeded", len(exercises))
	}
	typed, writes := 0, 0
	for _, e := range exercises {
		switch e.Type {
		case models.ExerciseTranslate, models.ExerciseFill:
			typed++
			if strings.TrimSpace(e.CorrectAnswer) == "" || strings.TrimSpace(e.Question) == "" {
				t.Errorf("typed exercise %q has no question or answer", e.Question)
			}
		case models.ExerciseMultipleChoice, models.ExerciseListen, models.ExerciseMatch:
			var opts []string
			_ = json.Unmarshal([]byte(e.OptionsJSON), &opts)
			found, dup := false, map[string]bool{}
			for _, o := range opts {
				found = found || o == e.CorrectAnswer
				if dup[o] {
					t.Errorf("exercise %q repeats option %q", e.Question, o)
				}
				dup[o] = true
			}
			if !found || len(opts) < 3 {
				t.Errorf("exercise %q: answer %q not among %d options", e.Question, e.CorrectAnswer, len(opts))
			}
		case models.ExerciseWrite:
			writes++
			if e.CorrectAnswer == "" {
				t.Errorf("writing task %q has no example answer", e.Question)
			}
		}
	}
	// Swahili is typed in the Latin alphabet, so production practice is
	// built in from A1 — unlike the character-based courses.
	if typed < 50 {
		t.Errorf("only %d typed (translate/fill) exercises", typed)
	}
	if writes < 6 {
		t.Errorf("only %d writing tasks", writes)
	}

	var vocab int64
	db.Model(&models.VocabItem{}).
		Joins("JOIN lessons ON lessons.id = vocab_items.lesson_id").
		Joins("JOIN skills ON skills.id = lessons.skill_id").
		Where("skills.language = ?", "it").Count(&vocab)
	if vocab < 200 {
		t.Errorf("only %d vocabulary items", vocab)
	}
}

func TestItalianListeningAndReading(t *testing.T) {
	db := seededItalian(t)
	var listening, reading int64
	db.Model(&models.ListeningSession{}).Where("language = ?", "it").Count(&listening)
	db.Model(&models.ReadingSession{}).Where("language = ?", "it").Count(&reading)
	if listening < 7 || reading < 7 {
		t.Errorf("listening = %d, reading = %d; want at least one per level", listening, reading)
	}
}
