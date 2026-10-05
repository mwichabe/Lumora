package controllers

import (
	"testing"

	"lumora/backend/models"
)

func kinds(ps []writingProblem) map[string]bool {
	out := map[string]bool{}
	for _, p := range ps {
		out[p.Kind] = true
	}
	return out
}

const hotelTask = "Write a short email to a hotel: say when you arrive and ask for a room and the price."
const hotelSample = "Sehr geehrte Damen und Herren, ich komme am Freitag an. Haben Sie ein Zimmer frei? Was kostet das Zimmer pro Nacht? Mit freundlichen Grüßen, Anna"

func TestCopiedExampleIsRejected(t *testing.T) {
	cases := map[string]string{
		"verbatim":           hotelSample,
		"lightly edited":     "sehr geehrte damen und herren ich komme am freitag an haben sie ein zimmer frei was kostet das zimmer pro nacht",
		"example plus extra": hotelSample + " Danke schön und bis bald, ich freue mich sehr.",
		"most of it":         "Ich komme am Freitag an. Haben Sie ein Zimmer frei? Was kostet das Zimmer pro Nacht?",
	}
	for name, text := range cases {
		if !kinds(writingProblems(text, hotelTask, hotelSample, "de", 12))["copied_example"] {
			t.Errorf("%s: copy of the example was accepted", name)
		}
	}
}

func TestOwnAnswerPasses(t *testing.T) {
	text := "Hallo, mein Name ist Peter. Ich reise am Montag nach Berlin und brauche ein Einzelzimmer für drei Nächte. Wie viel kostet es mit Frühstück? Vielen Dank, Peter"
	if ps := writingProblems(text, hotelTask, hotelSample, "de", 12); len(ps) != 0 {
		t.Fatalf("original answer flagged: %+v", ps)
	}
}

func TestShortWrongLanguageAndPaddedAnswersAreRejected(t *testing.T) {
	if !kinds(writingProblems("Hallo, ich komme.", hotelTask, hotelSample, "de", 12))["too_short"] {
		t.Error("a 3-word answer passed a 12-word minimum")
	}
	english := "Hello, my name is Peter. I will arrive on Monday and I need a single room for three nights. How much does it cost with breakfast?"
	if !kinds(writingProblems(english, hotelTask, hotelSample, "de", 12))["wrong_language"] {
		t.Error("an English answer to a German task passed")
	}
	padded := "Zimmer Zimmer Zimmer Zimmer Zimmer Zimmer Zimmer Zimmer Zimmer Zimmer Zimmer Zimmer Zimmer"
	if !kinds(writingProblems(padded, hotelTask, hotelSample, "de", 12))["repetitive"] {
		t.Error("padding passed")
	}
}

func TestMinWordsReadsTheTask(t *testing.T) {
	cases := map[string]int{
		"Write a short postcard from Madrid (30–40 words): greet…": 18,
		"Write your opinion (≈150 words): 'Soll man…'":             90,
		"Write an argumentative essay (250+ words): …":             150,
		"Write a formal letter (around 200 words) giving…":         120,
		"Write a short email to a hotel: say when you arrive.":     12,
	}
	for task, want := range cases {
		if got := minWordsFor(task); got != want {
			t.Errorf("%q: min %d, want %d", task, got, want)
		}
	}
}

func lessonWith(exs ...models.Exercise) *models.Lesson {
	return &models.Lesson{
		Exercises: exs,
		Vocab: []models.VocabItem{
			{Word: "der Mann", Translation: "the man", Example: "Der Mann liest.", ExampleTranslation: "The man reads."},
			{Word: "die Frau", Translation: "the woman", Example: "Die Frau arbeitet.", ExampleTranslation: "The woman works."},
		},
	}
}

func typedCount(l *models.Lesson) (typed, picked int) {
	for _, e := range l.Exercises {
		if e.Type != models.ExerciseTranslate && e.Type != models.ExerciseFill {
			continue
		}
		if len(e.Options) == 0 {
			typed++
		} else {
			picked++
		}
	}
	return
}

func TestA1AlternatesAndAddsVocabWriting(t *testing.T) {
	l := lessonWith(
		models.Exercise{Type: models.ExerciseTranslate, Question: "I am happy", CorrectAnswer: "Ich bin glücklich"},
		models.Exercise{Type: models.ExerciseFill, Question: "Ich ___ müde.", CorrectAnswer: "bin"},
		models.Exercise{Type: models.ExerciseTranslate, Question: "You are tired", CorrectAnswer: "Du bist müde"},
		models.Exercise{Type: models.ExerciseWrite, Question: hotelTask, CorrectAnswer: hotelSample},
	)
	applyAnswerModes(l, "de", "A1")
	typed, picked := typedCount(l)
	// 3 seeded: pick, type, pick — plus 2 typed vocab items.
	if typed != 3 || picked != 2 {
		t.Fatalf("typed %d / picked %d, want 3 / 2", typed, picked)
	}
	if last := l.Exercises[len(l.Exercises)-1]; last.Type != models.ExerciseWrite {
		t.Errorf("vocab writing went after the closing writing task")
	}
}

func TestA2AndAboveAreAllTyped(t *testing.T) {
	l := lessonWith(
		models.Exercise{Type: models.ExerciseTranslate, Question: "I am happy", CorrectAnswer: "Ich bin glücklich"},
		models.Exercise{Type: models.ExerciseFill, Question: "Ich ___ müde.", CorrectAnswer: "bin"},
	)
	applyAnswerModes(l, "de", "B1")
	if _, picked := typedCount(l); picked != 0 {
		t.Fatalf("%d exercises still multiple choice at B1", picked)
	}
}

func TestChineseStaysMultipleChoice(t *testing.T) {
	l := lessonWith(models.Exercise{Type: models.ExerciseTranslate, Question: "hello", CorrectAnswer: "你好"})
	applyAnswerModes(l, "zh", "A2")
	if typed, _ := typedCount(l); typed != 0 || len(l.Exercises) != 1 {
		t.Fatalf("Chinese got %d typed exercises (%d total)", typed, len(l.Exercises))
	}
}
