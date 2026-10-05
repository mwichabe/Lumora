package controllers

import "testing"

func TestJapaneseExamBankIsComplete(t *testing.T) {
	levels := []string{"A1", "A2", "B1", "B2", "C1", "C2", "FINAL"}
	for _, level := range levels {
		p, ok := japanesePapers[level]
		if !ok {
			t.Errorf("no Japanese paper for %s", level)
			continue
		}
		if len(p.Listening.Lines) == 0 || len(p.Listening.Questions) < 4 {
			t.Errorf("%s listening is too thin", level)
		}
		if len(p.Reading.Paragraphs) == 0 || len(p.Reading.Questions) < 4 {
			t.Errorf("%s reading is too thin", level)
		}
		if p.Writing.MinWords <= 0 || p.Writing.Prompt == "" || p.Speaking.Phrase == "" {
			t.Errorf("%s is missing writing or speaking", level)
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
	// Writing targets are character counts and must rise with the level.
	prev := 0
	for _, level := range levels {
		if n := japanesePapers[level].Writing.MinWords; n <= prev {
			t.Errorf("%s writing target %d doesn't rise above %d", level, n, prev)
		} else {
			prev = n
		}
	}
}
