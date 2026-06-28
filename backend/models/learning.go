package models

import "time"

// Enrollment records a language a user is learning. A user can be enrolled in
// several languages; User.TargetLanguage holds the currently active one.
type Enrollment struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index" json:"userId"`
	Language  string    `json:"language"`
	CreatedAt time.Time `json:"createdAt"`
}

// Mistake records an exercise a user got wrong, so it can be resurfaced in the
// "Review Mistakes" practice mode and cleared once answered correctly.
type Mistake struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	UserID        uint      `gorm:"index" json:"userId"`
	Language      string    `json:"language"`
	Prompt        string    `json:"prompt"`
	Question      string    `json:"question"`
	CorrectAnswer string    `json:"correctAnswer"`
	CreatedAt     time.Time `json:"createdAt"`
}

// VocabItem is a new word/phrase taught at the START of a lesson, before any
// questions — the "learn the words first" phase.
type VocabItem struct {
	ID                 uint   `gorm:"primaryKey" json:"id"`
	LessonID           uint   `json:"lessonId"`
	OrderIndex         int    `json:"orderIndex"`
	Word               string `json:"word"`               // target language
	Translation        string `json:"translation"`        // native meaning
	Example            string `json:"example"`            // example sentence (target)
	ExampleTranslation string `json:"exampleTranslation"` // example sentence (native)
	Speaker            string `json:"speaker"`            // character whose voice introduces it
}

// ListeningSession is a unit-level verbal session: a short dialogue performed by
// the app characters, followed by comprehension questions about what was heard.
type ListeningSession struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Language    string `json:"language"`
	Unit        string `json:"unit"`
	Title       string `json:"title"`
	Description string `json:"description"`
	OrderIndex  int    `json:"orderIndex"`
	XPReward    int    `json:"xpReward"`

	Matches   []ListeningMatch    `gorm:"foreignKey:SessionID" json:"matches,omitempty"`
	Lines     []ListeningLine     `gorm:"foreignKey:SessionID" json:"lines,omitempty"`
	Questions []ListeningQuestion `gorm:"foreignKey:SessionID" json:"questions,omitempty"`
}

// ListeningMatch is a word pair used in the warm-up matching game shown BEFORE
// the conversation plays. The words appear in the dialogue the user will hear.
type ListeningMatch struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	SessionID   uint   `json:"sessionId"`
	OrderIndex  int    `json:"orderIndex"`
	Word        string `json:"word"`        // target language (e.g. Spanish)
	Translation string `json:"translation"` // native (English)
}

// ListeningLine is one spoken line of the dialogue, attributed to a character so
// the frontend can render and voice it distinctly.
type ListeningLine struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	SessionID   uint   `json:"sessionId"`
	OrderIndex  int    `json:"orderIndex"`
	Character   string `json:"character"`
	Text        string `json:"text"`        // spoken line (target language)
	Translation string `json:"translation"` // native translation (revealed in transcript)
}

// ReadingSession is a unit-level reading passage in the target language with
// comprehension questions — a "read it yourself" complement to listening.
type ReadingSession struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Language    string `json:"language"`
	Unit        string `json:"unit"`
	Title       string `json:"title"`
	Description string `json:"description"`
	OrderIndex  int    `json:"orderIndex"`
	XPReward    int    `json:"xpReward"`

	Lines     []ReadingLine     `gorm:"foreignKey:SessionID" json:"lines,omitempty"`
	Questions []ReadingQuestion `gorm:"foreignKey:SessionID" json:"questions,omitempty"`
}

type ReadingLine struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	SessionID   uint   `json:"sessionId"`
	OrderIndex  int    `json:"orderIndex"`
	Text        string `json:"text"`        // target-language sentence
	Translation string `json:"translation"` // native translation (revealable)
}

type ReadingQuestion struct {
	ID            uint     `gorm:"primaryKey" json:"id"`
	SessionID     uint     `json:"sessionId"`
	OrderIndex    int      `json:"orderIndex"`
	Prompt        string   `json:"prompt"`
	Question      string   `json:"question"`
	OptionsJSON   string   `gorm:"column:options" json:"-"`
	Options       []string `gorm:"-" json:"options"`
	CorrectAnswer string   `json:"correctAnswer"`
}

// ListeningQuestion is a comprehension check asked after the dialogue.
type ListeningQuestion struct {
	ID            uint     `gorm:"primaryKey" json:"id"`
	SessionID     uint     `json:"sessionId"`
	OrderIndex    int      `json:"orderIndex"`
	Prompt        string   `json:"prompt"`
	Question      string   `json:"question"`
	OptionsJSON   string   `gorm:"column:options" json:"-"`
	Options       []string `gorm:"-" json:"options"`
	CorrectAnswer string   `json:"correctAnswer"`
}
