package models

import "time"

// User is the core account model. It holds learning progress, gamification
// state (XP, gems, streaks) and the gameplay metadata Lumora needs.
type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Email        string    `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash string    `gorm:"not null" json:"-"`
	Name         string    `json:"name"`
	AvatarColor  string    `json:"avatarColor"` // hex used for the placeholder avatar ring
	AvatarURL    string    `json:"avatarUrl"`   // uploaded profile photo (served from /uploads)

	// Learning setup (chosen during onboarding)
	TargetLanguage string `json:"targetLanguage"` // e.g. "es"
	NativeLanguage string `json:"nativeLanguage"` // e.g. "en"
	CEFRLevel      string `json:"cefrLevel"`      // A1..C2
	LevelName      string `json:"levelName"`      // Spark, Glow, Flame...
	DailyGoalXP    int    `json:"dailyGoalXp"`    // 10 / 20 / 30 / 50

	// Gamification state
	XP            int    `json:"xp"`
	XPToday       int    `json:"xpToday"`
	Gems          int    `json:"gems"`
	Hearts        int    `json:"hearts"`
	Streak        int    `json:"streak"`
	FluencyScore  int    `json:"fluencyScore"` // 0..1000
	League        string `json:"league"`       // Bronze..Obsidian

	LastActiveDate string    `json:"lastActiveDate"` // YYYY-MM-DD, drives streak logic
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}
