package models

import "time"

// Notification is an in-app message delivered to a single user. Some are sent on
// events (welcome) and others are pushed automatically by the campaign
// scheduler (tips, new features, upcoming languages).
type Notification struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index" json:"userId"`
	Key       string    `gorm:"index" json:"key"` // campaign key, for dedup ("" = one-off)
	Kind      string    `json:"kind"`             // welcome | tip | feature | language | streak
	Emoji     string    `json:"emoji"`
	Tint      string    `json:"tint"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	Link      string    `json:"link"` // optional in-app destination (e.g. "/chat/3")
	Read      bool      `json:"read"`
	CreatedAt time.Time `json:"createdAt"`
}
