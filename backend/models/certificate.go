package models

import "time"

// Certificate is awarded when a user passes the proficiency exam for a language.
// (Both the exam and the certificate are paid features — payment is bypassed for
// now.) Section scores mirror the four exam skills.
type Certificate struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index" json:"userId"`
	UserName  string    `json:"userName"`
	Language  string    `json:"language"`
	Level     string    `json:"level"` // CEFR: A1..C2
	Score     int       `json:"score"` // overall 0..100
	Listening int       `json:"listening"`
	Reading   int       `json:"reading"`
	Writing   int       `json:"writing"`
	Speaking  int       `json:"speaking"`
	Serial    string    `json:"serial"` // human-friendly certificate id
	IssuedAt  time.Time `json:"issuedAt"`
}
