package database

import (
	"encoding/json"
	"strings"
	"testing"
	"unicode"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"lumora/backend/models"
)

func seededHindi(t *testing.T) *gorm.DB {
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
	seedHindi(db)
	return db
}

func hasDevanagari(s string) bool {
	for _, r := range s {
		if unicode.Is(unicode.Devanagari, r) {
			return true
		}
	}
	return false
}

func TestHindiCourseCoversA1ToC2(t *testing.T) {
	db := seededHindi(t)

	var skills []models.Skill
	db.Where("language = ?", "hi").Preload("Lessons").Order("order_index").Find(&skills)

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

	// Script first: the first three skills are Devanagari, the first open
	// from the start.
	for i, want := range []string{"Devanagari I: Vowels", "Devanagari II: Consonants", "Devanagari III: Matras & Conjuncts"} {
		if skills[i].Title != want {
			t.Errorf("skill %d = %q, want %q", i+1, skills[i].Title, want)
		}
	}
	if skills[0].RequiredXP != 0 {
		t.Errorf("the first skill needs %d XP; it must be open from the start", skills[0].RequiredXP)
	}
	titles := map[string]bool{}
	for _, s := range skills {
		titles[s.Title] = true
	}
	for _, must := range []string{"Gender & Adjectives", "Postpositions", "The ने Construction",
		"Honorifics in Practice", "Idioms (मुहावरे)", "Hindi, Urdu & Hindustani"} {
		if !titles[must] {
			t.Errorf("missing skill %q", must)
		}
	}
}

func TestHindiExercisesAreSelfContained(t *testing.T) {
	db := seededHindi(t)

	var exercises []models.Exercise
	db.Joins("JOIN lessons ON lessons.id = exercises.lesson_id").
		Joins("JOIN skills ON skills.id = lessons.skill_id").
		Where("skills.language = ?", "hi").Find(&exercises)
	if len(exercises) < 280 {
		t.Fatalf("only %d exercises seeded", len(exercises))
	}
	writes := 0
	for _, e := range exercises {
		switch e.Type {
		case models.ExerciseTranslate, models.ExerciseFill:
			t.Errorf("exercise %q uses %s; Hindi exercises must carry their own options", e.Question, e.Type)
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
			if !hasDevanagari(e.CorrectAnswer) {
				t.Errorf("writing task %q: its example must be in Devanagari", e.Question)
			}
		}
	}
	if writes < 6 {
		t.Errorf("only %d writing tasks", writes)
	}

	var vocab []models.VocabItem
	db.Joins("JOIN lessons ON lessons.id = vocab_items.lesson_id").
		Joins("JOIN skills ON skills.id = lessons.skill_id").
		Where("skills.language = ?", "hi").Find(&vocab)
	if len(vocab) < 200 {
		t.Errorf("only %d vocabulary items", len(vocab))
	}
	for _, v := range vocab {
		// Devanagari is what's spoken; the transliteration rides along.
		if !hasDevanagari(v.Word) {
			t.Errorf("vocab %q is not in Devanagari", v.Word)
		}
		if !strings.Contains(v.Translation, " · ") {
			t.Errorf("vocab %q is missing its transliteration: %q", v.Word, v.Translation)
		}
	}
}

func TestHindiListeningAndReading(t *testing.T) {
	db := seededHindi(t)
	var listening, reading int64
	db.Model(&models.ListeningSession{}).Where("language = ?", "hi").Count(&listening)
	db.Model(&models.ReadingSession{}).Where("language = ?", "hi").Count(&reading)
	if listening < 7 || reading < 7 {
		t.Errorf("listening = %d, reading = %d; want at least one per level", listening, reading)
	}
}
