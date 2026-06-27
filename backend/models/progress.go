package models

import "time"

// LessonProgress records a user's completion of a lesson.
type LessonProgress struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `gorm:"index" json:"userId"`
	LessonID    uint      `gorm:"index" json:"lessonId"`
	Completed   bool      `json:"completed"`
	Accuracy    int       `json:"accuracy"`  // percentage 0..100
	XPEarned    int       `json:"xpEarned"`
	CompletedAt time.Time `json:"completedAt"`
}

// Character is a member of the Lumora cast (Lumora, Professor Finch, Cora...).
type Character struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Name        string `json:"name"`
	Species     string `json:"species"`
	Role        string `json:"role"`
	Personality string `json:"personality"`
	Color       string `json:"color"`
	Emoji       string `json:"emoji"`
}

// Friendship tracks a user's relationship level (1-10) with a character.
type Friendship struct {
	ID          uint `gorm:"primaryKey" json:"id"`
	UserID      uint `gorm:"index" json:"userId"`
	CharacterID uint `json:"characterId"`
	Level       int  `json:"level"`
	XP          int  `json:"xp"`

	Character *Character `gorm:"foreignKey:CharacterID" json:"character,omitempty"`
}

// Quest is a daily goal delivered by Pip.
type Quest struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	XPReward    int    `json:"xpReward"`
	Target      int    `json:"target"` // e.g. complete 2 lessons -> target 2
}

// UserQuest is a per-user, per-day instance of a quest with live progress.
type UserQuest struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	UserID    uint   `gorm:"index" json:"userId"`
	QuestID   uint   `json:"questId"`
	Date      string `gorm:"index" json:"date"` // YYYY-MM-DD
	Progress  int    `json:"progress"`
	Completed bool   `json:"completed"`

	Quest *Quest `gorm:"foreignKey:QuestID" json:"quest,omitempty"`
}
