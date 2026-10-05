package controllers

import (
	"testing"
	"unicode"
)

func TestHindiExamBankIsComplete(t *testing.T) {
	levels := []string{"A1", "A2", "B1", "B2", "C1", "C2", "FINAL"}
	prev := 0
	for _, level := range levels {
		p, ok := hindiPapers[level]
		if !ok {
			t.Errorf("no Hindi paper for %s", level)
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
		// The texts themselves are in Devanagari.
		devanagari := false
		for _, r := range p.Speaking.Phrase + p.Reading.Paragraphs[0] {
			devanagari = devanagari || unicode.Is(unicode.Devanagari, r)
		}
		if !devanagari {
			t.Errorf("%s texts are not in Devanagari", level)
		}
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

// Hindi is written in Devanagari, so its lessons are answered by choosing —
// the same rule as Chinese and Japanese.
func TestHindiLessonsAreChoiceOnly(t *testing.T) {
	l := lessonWith()
	applyAnswerModes(l, "hi", "B2")
	if len(l.Exercises) != 0 {
		t.Errorf("Hindi got %d generated typed exercises", len(l.Exercises))
	}
}

// Hindi writing is counted in words: postpositions are separate words, and
// danda punctuation isn't a word.
func TestHindiWritingIsCountedByWord(t *testing.T) {
	if n := countWritingWords("मेरा नाम अन्ना है। मैं केन्या से हूँ।"); n != 8 {
		t.Errorf("counted %d words, want 8", n)
	}
	sample := "मैं नैरोबी में रहती हूँ। सुबह थोड़ी सर्दी होती है, लेकिन दोपहर में धूप निकलती है।"
	if !kinds(writingProblems(sample, "Write about the weather…", sample, "hi", 12))["copied_example"] {
		t.Error("a Hindi copy of the example was accepted")
	}
}
