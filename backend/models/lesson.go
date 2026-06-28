package models

import "time"

// Skill is a node on the "galaxy map" skill tree. Each skill groups a set of
// lessons under a theme (Greetings, Food, Travel...).
type Skill struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Language    string `json:"language"` // target language code this skill belongs to
	Unit        string `json:"unit"`     // section the skill belongs to (e.g. "Basics")
	Title       string `json:"title"`
	Description string `json:"description"`
	Icon        string `json:"icon"`  // lucide icon name used by the frontend
	Color       string `json:"color"` // accent hex for the node
	OrderIndex  int    `json:"orderIndex"`
	RequiredXP  int    `json:"requiredXp"` // XP needed before this node unlocks

	Lessons []Lesson `gorm:"foreignKey:SkillID" json:"lessons,omitempty"`
}

// Lesson is a bite-sized unit (3-7 min) made of several exercises.
type Lesson struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	SkillID    uint   `json:"skillId"`
	Title      string `json:"title"`
	OrderIndex int    `json:"orderIndex"`
	XPReward   int    `json:"xpReward"`

	Vocab     []VocabItem `gorm:"foreignKey:LessonID" json:"vocab,omitempty"`
	Exercises []Exercise  `gorm:"foreignKey:LessonID" json:"exercises,omitempty"`
}

// ExerciseType enumerates the supported exercise formats from the design spec.
type ExerciseType string

const (
	ExerciseTranslate      ExerciseType = "translate"
	ExerciseMultipleChoice ExerciseType = "multiple_choice"
	ExerciseListen         ExerciseType = "listen"
	ExerciseMatch          ExerciseType = "match"
	ExerciseSpeak          ExerciseType = "speak"
	ExerciseFill           ExerciseType = "fill"
	ExerciseWrite          ExerciseType = "write"     // free-text writing (e.g. an email)
	ExerciseCharacter      ExerciseType = "character" // a narrative interjection
)

// Exercise is a single question/interaction inside a lesson.
type Exercise struct {
	ID            uint         `gorm:"primaryKey" json:"id"`
	LessonID      uint         `json:"lessonId"`
	Type          ExerciseType `json:"type"`
	OrderIndex    int          `json:"orderIndex"`
	Prompt        string       `json:"prompt"`                  // instruction shown to the user
	Question      string       `json:"question"`                // sentence / audio caption / phrase
	OptionsJSON   string       `gorm:"column:options" json:"-"` // stored as JSON string
	Options       []string     `gorm:"-" json:"options"`        // hydrated for the API
	CorrectAnswer string       `json:"correctAnswer"`
	Character     string       `json:"character"` // speaking character name, if any
	CreatedAt     time.Time    `json:"-"`
}
