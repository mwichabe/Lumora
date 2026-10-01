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

func seededMandarin(t *testing.T) *gorm.DB {
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
	seedMandarin(db)
	return db
}

func TestMandarinCourseCoversHSK1To6(t *testing.T) {
	db := seededMandarin(t)

	var skills []models.Skill
	db.Where("language = ?", "zh").Preload("Lessons").Order("order_index").Find(&skills)

	// One CEFR-tagged unit per HSK level, so the level system can promote
	// learners through it (see controllers/level.go).
	want := []string{"A1 · HSK 1", "A2 · HSK 2", "B1 · HSK 3", "B2 · HSK 4", "C1 · HSK 5", "C2 · HSK 6"}
	seen := map[string]int{}
	for _, s := range skills {
		for _, prefix := range want {
			if strings.HasPrefix(s.Unit, prefix) {
				seen[prefix]++
			}
		}
		if len(s.Lessons) == 0 {
			t.Errorf("skill %q has no lessons", s.Title)
		}
	}
	for _, prefix := range want {
		if seen[prefix] < 4 {
			t.Errorf("unit %q has %d skills, want at least 4", prefix, seen[prefix])
		}
	}

	// HSK 1 starts with the foundations: Pinyin and tones come first.
	if skills[0].Title != "Pinyin & the Four Tones" || skills[0].RequiredXP != 0 {
		t.Errorf("first skill = %q (req %d XP), want Pinyin & the Four Tones open from the start", skills[0].Title, skills[0].RequiredXP)
	}
}

func TestMandarinExercisesAreSelfContained(t *testing.T) {
	db := seededMandarin(t)

	var exercises []models.Exercise
	db.Joins("JOIN lessons ON lessons.id = exercises.lesson_id").
		Joins("JOIN skills ON skills.id = lessons.skill_id").
		Where("skills.language = ?", "zh").Find(&exercises)
	if len(exercises) < 150 {
		t.Fatalf("only %d exercises seeded", len(exercises))
	}
	for _, e := range exercises {
		switch e.Type {
		case models.ExerciseTranslate, models.ExerciseFill:
			// These get distractors from a generic (Spanish) pool at serve time.
			t.Errorf("exercise %q uses %s; Mandarin exercises must carry their own options", e.Question, e.Type)
		case models.ExerciseMultipleChoice, models.ExerciseListen:
			var opts []string
			_ = json.Unmarshal([]byte(e.OptionsJSON), &opts)
			found := false
			for _, o := range opts {
				found = found || o == e.CorrectAnswer
			}
			if !found || len(opts) < 3 {
				t.Errorf("exercise %q: answer %q not among %d options", e.Question, e.CorrectAnswer, len(opts))
			}
			if strings.TrimSpace(e.Question) == "" {
				t.Errorf("exercise %q has no question line", e.Prompt)
			}
		}
	}

	var vocab []models.VocabItem
	db.Joins("JOIN lessons ON lessons.id = vocab_items.lesson_id").
		Joins("JOIN skills ON skills.id = lessons.skill_id").
		Where("skills.language = ?", "zh").Find(&vocab)
	for _, v := range vocab {
		// Every word shows its Pinyin alongside the meaning.
		if !strings.Contains(v.Translation, " · ") {
			t.Errorf("vocab %q is missing Pinyin: %q", v.Word, v.Translation)
		}
	}
}

func TestMandarinListeningAndReading(t *testing.T) {
	db := seededMandarin(t)
	var listening, reading int64
	db.Model(&models.ListeningSession{}).Where("language = ?", "zh").Count(&listening)
	db.Model(&models.ReadingSession{}).Where("language = ?", "zh").Count(&reading)
	if listening < 7 || reading < 7 {
		t.Errorf("listening = %d, reading = %d; want at least one per HSK level", listening, reading)
	}
}
