package controllers

import (
	"strings"
	"testing"

	"lumora/backend/models"
)

func TestSwahiliExamBankIsComplete(t *testing.T) {
	levels := []string{"A1", "A2", "B1", "B2", "C1", "C2", "FINAL"}
	prev := 0
	for _, level := range levels {
		p, ok := swahiliPapers[level]
		if !ok {
			t.Errorf("no Swahili paper for %s", level)
			continue
		}
		if len(p.Listening.Lines) == 0 || len(p.Listening.Questions) < 4 {
			t.Errorf("%s listening is too thin", level)
		}
		if len(p.Reading.Paragraphs) == 0 || len(p.Reading.Questions) < 4 {
			t.Errorf("%s reading is too thin", level)
		}
		if p.Writing.MinWords <= prev || p.Speaking.Phrase == "" {
			t.Errorf("%s: writing target %d must rise above %d, and speaking must be set", level, p.Writing.MinWords, prev)
		}
		prev = p.Writing.MinWords
		for _, q := range append(p.Listening.Questions, p.Reading.Questions...) {
			found, seen := false, map[string]bool{}
			for _, o := range q.Options {
				found = found || o == q.CorrectAnswer
				if seen[o] {
					t.Errorf("%s: %q repeats option %q", level, q.Question, o)
				}
				seen[o] = true
			}
			if !found {
				t.Errorf("%s: answer %q is not one of its options", level, q.CorrectAnswer)
			}
		}
	}
}

// A Swahili question answered by choosing must never be offered Spanish
// options — the old filler list was Spanish for every language.
func TestSwahiliChoicesUseSwahiliFillers(t *testing.T) {
	l := &models.Lesson{Exercises: []models.Exercise{
		{Type: models.ExerciseTranslate, Question: "Thank you very much", CorrectAnswer: "Asante sana"},
		{Type: models.ExerciseTranslate, Question: "Good night", CorrectAnswer: "Lala salama"},
	}}
	applyAnswerModes(l, "sw", "A1")
	spanish := []string{"Buenos días", "Por favor", "Hasta luego", "No lo sé", "Mucho gusto", "gracias", "hola"}
	for _, e := range l.Exercises {
		for _, o := range e.Options {
			for _, es := range spanish {
				if strings.EqualFold(o, es) {
					t.Errorf("Swahili exercise %q offered Spanish option %q", e.Question, o)
				}
			}
		}
	}
	if len(l.Exercises[0].Options) < 3 {
		t.Errorf("the A1 choice exercise has only %d options", len(l.Exercises[0].Options))
	}
}
