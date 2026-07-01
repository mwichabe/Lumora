package models

import "time"

// PasswordReset is a single-use, time-limited token that lets a user set a new
// password without being logged in (the "forgot password" flow).
type PasswordReset struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index" json:"userId"`
	Token     string    `gorm:"uniqueIndex" json:"-"`
	ExpiresAt time.Time `json:"-"`
	Used      bool      `json:"-"`
	CreatedAt time.Time `json:"-"`
}
