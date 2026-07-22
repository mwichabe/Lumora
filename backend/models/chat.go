package models

import "time"

// Message is a direct 1:1 chat message between two users.
//
// Kept separate from the ideas workspace (models/idea.go) on purpose: folding
// team discussion into personal DMs is what turns both into noise.
type Message struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	SenderID    uint   `gorm:"index" json:"senderId"`
	RecipientID uint   `gorm:"index" json:"recipientId"`
	Body        string `json:"body"`
	Read        bool   `json:"read"`

	// Kind is "text" or "image". Attachments live in the database rather than
	// on disk — the production filesystem is ephemeral and would drop them on
	// the next deploy. Served by GET /api/chat/attachments/:id.
	Kind     string `json:"kind"`
	Data     []byte `json:"-"`
	Mime     string `json:"-"`
	FileName string `json:"fileName"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`

	// Language of the message body, detected offline when it was sent, plus
	// the English translation when one was produced. Both are computed once
	// and stored: detection is free but deterministic, and translation costs
	// an API call that shouldn't be repeated on every read.
	DetectedLang   string     `json:"detectedLang"` // ISO 639-1, "" when undetermined
	TranslatedBody string     `json:"-"`
	TranslatedAt   *time.Time `json:"-"`

	// A deleted message leaves a tombstone rather than disappearing, and an
	// edited one is always marked — silently rewriting a conversation the other
	// person has already read is worse than not being able to edit at all.
	EditedAt  *time.Time `json:"editedAt"`
	DeletedAt *time.Time `gorm:"index" json:"deletedAt"`

	CreatedAt time.Time `json:"createdAt"`
}
