package models

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
