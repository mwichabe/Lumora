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

func seededJapanese(t *testing.T) *gorm.DB {
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
	seedJapanese(db)
	return db
}

func TestJapaneseCourseCoversN5ToBeyondN1(t *testing.T) {
	db := seededJapanese(t)

	var skills []models.Skill
	db.Where("language = ?", "ja").Preload("Lessons").Order("order_index").Find(&skills)

	// One CEFR-tagged unit per JLPT level, so the level system can promote
	// learners through it (see controllers/level.go).
	want := []string{"A1 · JLPT N5", "A2 · JLPT N4", "B1 · JLPT N3", "B2 · JLPT N2", "C1 · JLPT N1", "C2 · Beyond N1"}
	seen := map[string]int{}
	lastXP := -1
	for _, s := range skills {
		for _, prefix := range want {
			if strings.HasPrefix(s.Unit, prefix) {
				seen[prefix]++
			}
		}
		if len(s.Lessons) == 0 {
			t.Errorf("skill %q has no lessons", s.Title)
		}
		// Skills unlock in order: required XP never goes down.
		if s.RequiredXP < lastXP {
			t.Errorf("skill %q needs %d XP, less than the skill before it (%d)", s.Title, s.RequiredXP, lastXP)
		}
		lastXP = s.RequiredXP
	}
	for _, prefix := range want {
		if seen[prefix] < 5 {
			t.Errorf("unit %q has %d skills, want at least 5", prefix, seen[prefix])
		}
	}

	// N5 starts with the foundations: hiragana, then katakana, open from the start.
	if skills[0].Title != "Hiragana I" || skills[0].RequiredXP != 0 {
		t.Errorf("first skill = %q (req %d XP), want Hiragana I open from the start", skills[0].Title, skills[0].RequiredXP)
	}
	titles := map[string]bool{}
	for _, s := range skills {
		titles[s.Title] = true
	}
	for _, must := range []string{"Katakana", "Sounds & Pitch Accent", "Core Particles", "First Kanji (N5)",
		"Keigo Foundations", "Business Keigo in Full", "Classical Japanese (文語)", "Dialects: Kansai & Beyond"} {
		if !titles[must] {
			t.Errorf("missing skill %q", must)
		}
	}
}

func TestJapaneseExercisesAreSelfContained(t *testing.T) {
	db := seededJapanese(t)

	var exercises []models.Exercise
	db.Joins("JOIN lessons ON lessons.id = exercises.lesson_id").
		Joins("JOIN skills ON skills.id = lessons.skill_id").
		Where("skills.language = ?", "ja").Find(&exercises)
	if len(exercises) < 300 {
		t.Fatalf("only %d exercises seeded", len(exercises))
	}
	writes := 0
	for _, e := range exercises {
		switch e.Type {
		case models.ExerciseTranslate, models.ExerciseFill:
			// These get distractors from a generic (Spanish) pool at serve time.
			t.Errorf("exercise %q uses %s; Japanese exercises must carry their own options", e.Question, e.Type)
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
			if strings.TrimSpace(e.Question) == "" {
				t.Errorf("exercise %q has no question line", e.Prompt)
			}
		case models.ExerciseWrite:
			writes++
			if e.CorrectAnswer == "" {
				t.Errorf("writing task %q has no example answer", e.Question)
			}
		}
	}
	if writes < 5 {
		t.Errorf("only %d writing tasks; want guided composition at N4 and above", writes)
	}

	var vocab []models.VocabItem
	db.Joins("JOIN lessons ON lessons.id = vocab_items.lesson_id").
		Joins("JOIN skills ON skills.id = lessons.skill_id").
		Where("skills.language = ?", "ja").Find(&vocab)
	if len(vocab) < 250 {
		t.Errorf("only %d vocabulary items", len(vocab))
	}
	for _, v := range vocab {
		// Every word shows its reading alongside the meaning.
		if !strings.Contains(v.Translation, " · ") {
			t.Errorf("vocab %q is missing its reading: %q", v.Word, v.Translation)
		}
	}
}

func TestJapaneseListeningAndReading(t *testing.T) {
	db := seededJapanese(t)
	var listening, reading int64
	db.Model(&models.ListeningSession{}).Where("language = ?", "ja").Count(&listening)
	db.Model(&models.ReadingSession{}).Where("language = ?", "ja").Count(&reading)
	if listening < 7 || reading < 7 {
		t.Errorf("listening = %d, reading = %d; want at least one per level", listening, reading)
	}
}
