package controllers

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"

	"lumora/backend/database"
	"lumora/backend/middleware"
	"lumora/backend/models"
)

// This file builds the proficiency-exam "paper" the frontend runs. German has a
// hand-authored, level-scaled bank with lengthy listening/reading passages and
// progressively harder comprehension questions. Other languages fall back to a
// paper assembled from their seeded course content.

type PaperQuestion struct {
	Question      string   `json:"question"`
	Options       []string `json:"options"`
	CorrectAnswer string   `json:"correctAnswer"`
}

type PaperLine struct {
	Character   string `json:"character"`
	Text        string `json:"text"`
	Translation string `json:"translation"`
}

type PaperListening struct {
	Title     string          `json:"title"`
	Lines     []PaperLine     `json:"lines"`
	Questions []PaperQuestion `json:"questions"`
}

type PaperReading struct {
	Title      string          `json:"title"`
	Paragraphs []string        `json:"paragraphs"`
	Questions  []PaperQuestion `json:"questions"`
}

type PaperWriting struct {
	Prompt   string `json:"prompt"`
	MinWords int    `json:"minWords"`
}

type PaperSpeaking struct {
	Phrase      string `json:"phrase"`
	Speaker     string `json:"speaker"`
	Translation string `json:"translation"`
}

type ExamPaperDTO struct {
	Ready           bool            `json:"ready"`
	Language        string          `json:"language"`
	Level           string          `json:"level"`
	DurationSeconds int             `json:"durationSeconds"`
	PassMark        int             `json:"passMark"`
	Weights         map[string]int  `json:"weights"`
	Listening       *PaperListening `json:"listening"`
	Reading         *PaperReading   `json:"reading"`
	Writing         PaperWriting    `json:"writing"`
	Speaking        PaperSpeaking   `json:"speaking"`
}

type paperContent struct {
	Listening PaperListening
	Reading   PaperReading
	Writing   PaperWriting
	Speaking  PaperSpeaking
}

// Paper returns the exam paper for the user's language at the requested level.
func (ec *ExamController) Paper(c *fiber.Ctx) error {
	user := middleware.CurrentUser(c)
	lang := c.Query("language", user.TargetLanguage)
	if lang == "" {
		lang = "es"
	}
	level := c.Query("level", "A1")
	if _, ok := levelPassMark[level]; !ok {
		level = "A1"
	}

	dto := ExamPaperDTO{
		Language:        lang,
		Level:           level,
		DurationSeconds: durationForLevel(level),
		PassMark:        passMarkForLevel(level),
		Weights:         sectionWeights,
	}

	// German & Spanish have purpose-built, advanced paper banks.
	var bank map[string]paperContent
	switch lang {
	case "de":
		bank = germanPapers
	case "es":
		bank = spanishPapers
	}
	if bank != nil {
		if pc, ok := bank[level]; ok {
			l := pc.Listening
			r := pc.Reading
			dto.Listening = &l
			dto.Reading = &r
			dto.Writing = pc.Writing
			dto.Speaking = pc.Speaking
			dto.Ready = true
			return c.JSON(dto)
		}
	}

	// Any other language: assemble a paper from its seeded course content.
	assembleFromDB(&dto, lang, level)
	return c.JSON(dto)
}

// levelOrder maps a CEFR code to an index 0..5 for band selection. FINAL (the
// comprehensive mastery exam) maps to the hardest band for the DB fallback.
var levelOrder = map[string]int{"A1": 0, "A2": 1, "B1": 2, "B2": 3, "C1": 4, "C2": 5, "FINAL": 5}

var defaultWritingPrompts = []string{
	"a short email to a friend: introduce yourself, say where you live and what you do every day.",
	"an email to a colleague describing your last weekend and your plans for the next one.",
	"a forum post giving your opinion on whether young people should learn to cook, with reasons.",
	"a formal email of complaint to a hotel about problems during your stay, requesting a solution.",
	"an essay weighing the advantages and disadvantages of working from home, with examples.",
	"a structured argument on how technology is reshaping the way societies communicate.",
}
var defaultWritingMin = []int{30, 50, 80, 120, 170, 220}

// assembleFromDB builds a fallback paper from a language's seeded sessions.
func assembleFromDB(dto *ExamPaperDTO, lang, level string) {
	idx := levelOrder[level]

	var ls []models.ListeningSession
	database.DB.
		Preload("Lines", orderByIndex).
		Preload("Questions", orderByIndex).
		Where("language = ?", lang).Order("order_index asc").Find(&ls)
	if pick := pickSessionByLevel(len(ls), level); pick >= 0 {
		s := ls[pick]
		hydrateQuestions(s.Questions)
		pl := PaperListening{Title: s.Title}
		for _, ln := range s.Lines {
			pl.Lines = append(pl.Lines, PaperLine{Character: ln.Character, Text: ln.Text, Translation: ln.Translation})
		}
		for _, q := range s.Questions {
			pl.Questions = append(pl.Questions, PaperQuestion{Question: q.Question, Options: q.Options, CorrectAnswer: q.CorrectAnswer})
		}
		dto.Listening = &pl
	}

	var rs []models.ReadingSession
	database.DB.
		Preload("Lines", orderByIndex).
		Preload("Questions", orderByIndex).
		Where("language = ?", lang).Order("order_index asc").Find(&rs)
	if pick := pickSessionByLevel(len(rs), level); pick >= 0 {
		s := rs[pick]
		hydrateReadingQuestions(s.Questions)
		pr := PaperReading{Title: s.Title}
		for _, ln := range s.Lines {
			pr.Paragraphs = append(pr.Paragraphs, ln.Text)
		}
		for _, q := range s.Questions {
			pr.Questions = append(pr.Questions, PaperQuestion{Question: q.Question, Options: q.Options, CorrectAnswer: q.CorrectAnswer})
		}
		dto.Reading = &pr
	}

	dto.Writing = PaperWriting{Prompt: defaultWritingPrompts[idx], MinWords: defaultWritingMin[idx]}

	// Speaking phrase from a mid-course vocab item.
	var vocab []models.VocabItem
	database.DB.
		Joins("JOIN lessons ON lessons.id = vocab_items.lesson_id").
		Joins("JOIN skills ON skills.id = lessons.skill_id").
		Where("skills.language = ?", lang).
		Order("skills.order_index asc, vocab_items.id asc").Find(&vocab)
	if len(vocab) > 0 {
		v := vocab[len(vocab)/2]
		phrase := v.Word
		if idx >= 2 && v.Example != "" {
			phrase = v.Example
		}
		dto.Speaking = PaperSpeaking{Phrase: phrase, Speaker: v.Speaker, Translation: v.Translation}
	}

	dto.Ready = dto.Listening != nil && dto.Reading != nil
}

// pickSessionByLevel returns the index of the session to use for a level band,
// or -1 if there are none. Higher levels map to later (harder) sessions.
func pickSessionByLevel(n int, level string) int {
	if n == 0 {
		return -1
	}
	idx := levelOrder[level]
	pos := int(float64(idx) / 5.0 * float64(n-1))
	if pos < 0 {
		pos = 0
	}
	if pos > n-1 {
		pos = n - 1
	}
	return pos
}

func hydrateReadingQuestions(qs []models.ReadingQuestion) {
	for i := range qs {
		var o []string
		if qs[i].OptionsJSON != "" {
			_ = json.Unmarshal([]byte(qs[i].OptionsJSON), &o)
		}
		qs[i].Options = o
	}
}
