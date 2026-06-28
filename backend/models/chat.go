package models

import "time"

// Message is a direct 1:1 chat message between two users.
type Message struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	SenderID    uint      `gorm:"index" json:"senderId"`
	RecipientID uint      `gorm:"index" json:"recipientId"`
	Body        string    `json:"body"`
	Read        bool      `json:"read"`
	CreatedAt   time.Time `json:"createdAt"`
}
