package controllers

import (
	"context"
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"

	"lumora/backend/database"
	"lumora/backend/models"
	"lumora/backend/utils"
)

// Lesson writing: which questions the learner types rather than picks, extra
// typed practice built from the lesson's own vocabulary, and the checks run on
// what they write.

var coach *utils.WritingCoach

// InitWritingCoach wires the Claude-backed writing feedback (it may be a
// disabled coach — every caller copes with that).
func InitWritingCoach(c *utils.WritingCoach) { coach = c }

// --- choose or type --------------------------------------------------------

// applyAnswerModes decides how each translate/fill exercise is answered.
// An exercise with options is answered by choosing; one without is typed.
//
//   - From A2 up every translate/fill is typed: by then the learner has to
//     produce the language, not recognise it.
//   - At A1 they alternate — pick, type, pick, type — so beginners practise
//     writing from their first lesson without being thrown in at the deep end.
//   - Chinese and Japanese stay multiple choice throughout: answers are
//     characters, and a learner can't be assumed to have the keyboard set up.
//     (Their free-writing tasks are still typed.)
//
// Two typed exercises are then added from the lesson's vocabulary (see
// vocabWriting), so every lesson includes writing.
func applyAnswerModes(lesson *models.Lesson, lang, level string) {
	if lang == "zh" || lang == "ja" {
		// Their exercises are all authored with options; this only fills any
		// gaps.
		addChoiceOptions(lesson, lang)
		return
	}
	advanced := cefrIndex(level) >= 1

	eligible := 0
	var pick []int // exercises answered by choosing
	for i := range lesson.Exercises {
		e := &lesson.Exercises[i]
		if len(e.Options) > 0 || (e.Type != models.ExerciseTranslate && e.Type != models.ExerciseFill) {
			continue
		}
		eligible++
		if advanced || eligible%2 == 0 {
			e.Options = nil // typed
			e.Prompt = typedPrompt(e.Prompt, lang)
		} else {
			pick = append(pick, i)
		}
	}
	for _, i := range pick {
		lesson.Exercises[i].Options = choicesFor(lesson, &lesson.Exercises[i], lang)
	}

	lesson.Exercises = insertBeforeWrapUp(lesson.Exercises, vocabWriting(lesson, lang)...)
}

// choicesFor builds the multiple-choice options for one exercise.
func choicesFor(lesson *models.Lesson, e *models.Exercise, lang string) []string {
	tmp := models.Lesson{Vocab: lesson.Vocab, Exercises: []models.Exercise{*e}}
	for _, other := range lesson.Exercises {
		if other.Type == e.Type && other.CorrectAnswer != e.CorrectAnswer {
			tmp.Exercises = append(tmp.Exercises, models.Exercise{Type: other.Type, CorrectAnswer: other.CorrectAnswer})
		}
	}
	addChoiceOptions(&tmp, lang)
	return tmp.Exercises[0].Options
}

// typedPrompt makes a typed exercise's instruction say so.
func typedPrompt(prompt, lang string) string {
	p := strings.TrimSpace(prompt)
	if p == "" || strings.EqualFold(p, "translate") || strings.EqualFold(p, "translate this sentence") {
		return "Type this in " + utils.LanguageName(lang)
	}
	return p
}

// vocabWriting builds up to two typed exercises from the lesson's vocabulary:
// one word (English → target) and one example sentence. Items already used as
// an answer elsewhere in the lesson are skipped.
func vocabWriting(lesson *models.Lesson, lang string) []models.Exercise {
	used := map[string]bool{}
	for _, e := range lesson.Exercises {
		used[strings.ToLower(strings.TrimSpace(e.CorrectAnswer))] = true
	}
	name := utils.LanguageName(lang)
	var out []models.Exercise
	var wordDone, sentenceDone bool
	// Walk from the end: the later words of a lesson are the newest, so they
	// benefit most from being produced rather than recognised.
	for i := len(lesson.Vocab) - 1; i >= 0 && !(wordDone && sentenceDone); i-- {
		v := lesson.Vocab[i]
		if !wordDone && v.Word != "" && v.Translation != "" && !used[strings.ToLower(v.Word)] {
			out = append(out, models.Exercise{
				LessonID: lesson.ID, Type: models.ExerciseTranslate,
				Prompt:        "Type this word in " + name,
				Question:      v.Translation,
				CorrectAnswer: v.Word,
			})
			wordDone = true
			continue
		}
		if !sentenceDone && v.Example != "" && v.ExampleTranslation != "" && !used[strings.ToLower(v.Example)] {
			out = append(out, models.Exercise{
				LessonID: lesson.ID, Type: models.ExerciseTranslate,
				Prompt:        "Write this sentence in " + name,
				Question:      v.ExampleTranslation,
				CorrectAnswer: v.Example,
				Character:     v.Speaker,
			})
			sentenceDone = true
		}
	}
	return out
}

// insertBeforeWrapUp places extra exercises before a lesson's closing
// writing/speaking tasks (or at the end when it has none).
func insertBeforeWrapUp(exs []models.Exercise, extra ...models.Exercise) []models.Exercise {
	if len(extra) == 0 {
		return exs
	}
	at := len(exs)
	for at > 0 && (exs[at-1].Type == models.ExerciseWrite || exs[at-1].Type == models.ExerciseSpeak) {
		at--
	}
	out := make([]models.Exercise, 0, len(exs)+len(extra))
	out = append(out, exs[:at]...)
	out = append(out, extra...)
	out = append(out, exs[at:]...)
	for i := range out {
		out[i].OrderIndex = i
	}
	return out
}

// --- free writing ----------------------------------------------------------

type writingProblem struct {
	Kind    string `json:"kind"` // copied_example | copied_task | too_short | wrong_language | repetitive | empty
	Message string `json:"message"`
}

type checkWritingInput struct {
	Text string `json:"text"`
}

// CheckWriting marks a free-writing answer (POST /exercises/:id/check-writing).
//
// Rule checks run first and are authoritative: an answer that copies the
// example or the task, is too short, is in the wrong language, or pads with
// repetition comes back with `problems` and must be fixed before it counts —
// no heart is spent. Only an answer that passes them is sent to the coach for
// corrections, when the coach is configured.
func (l *LessonController) CheckWriting(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid exercise"})
	}
	var ex models.Exercise
	if database.DB.First(&ex, id).Error != nil || ex.Type != models.ExerciseWrite {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "writing task not found"})
	}
	var in checkWritingInput
	if err := c.BodyParser(&in); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
	}

	lang, level := lessonLanguage(ex.LessonID)
	words := countWritingWords(in.Text)
	minWords := minWordsFor(ex.Question)
	problems := writingProblems(in.Text, ex.Question, ex.CorrectAnswer, lang, minWords)

	res := fiber.Map{
		"ok": len(problems) == 0, "problems": problems,
		"wordCount": words, "minWords": minWords,
		"feedbackAvailable": false,
	}
	if len(problems) > 0 {
		return c.JSON(res)
	}

	if review, err := coach.ReviewWriting(c.Context(), lang, level, ex.Question, in.Text); err == nil {
		if review.Corrections == nil {
			review.Corrections = []utils.WritingCorrection{}
		}
		res["feedbackAvailable"] = true
		res["acceptable"] = review.Acceptable
		res["summary"] = review.Summary
		res["corrections"] = review.Corrections
	} else {
		// No coach (or it failed): the rule checks passed, so accept it.
		res["acceptable"] = true
		res["corrections"] = []utils.WritingCorrection{}
		res["summary"] = "Nice work — your answer meets the task."
	}
	return c.JSON(res)
}

// writingProblems runs the rule checks, in the order a learner should fix them.
func writingProblems(text, task, example, lang string, minWords int) []writingProblem {
	out := []writingProblem{}
	if strings.TrimSpace(text) == "" {
		return append(out, writingProblem{"empty", "Write your answer first."})
	}
	if copiesText(text, example) {
		out = append(out, writingProblem{"copied_example",
			"An example cannot be used here — write your own answer in your own words."})
	} else if copiesText(text, task) {
		out = append(out, writingProblem{"copied_task",
			"That repeats the task itself — write your answer to it."})
	}
	if n := countWritingWords(text); n < minWords {
		out = append(out, writingProblem{"too_short",
			"Write at least " + strconv.Itoa(minWords) + " words — you have " + strconv.Itoa(n) + "."})
	}
	if d := utils.Detect(text); lang != "" && d.Code != "" && d.Code != lang && d.Confidence >= utils.MinConfidence {
		out = append(out, writingProblem{"wrong_language",
			"This looks like " + d.Name + ". Write your answer in " + utils.LanguageName(lang) + "."})
	}
	if isRepetitive(text) {
		out = append(out, writingProblem{"repetitive",
			"Repeating the same words doesn't count — write real sentences."})
	}
	return out
}

// A unit is a word in space-separated languages, and a single character in
// Chinese and Japanese, which don't separate words — otherwise a whole
// sentence would count as one "word" and copy checks couldn't compare it.
var wordRe = regexp.MustCompile(`[\p{Han}\p{Hiragana}\p{Katakana}ー]|[^\s\p{Han}\p{Hiragana}\p{Katakana}ー\p{P}\p{S}]+`)

func writingWords(s string) []string {
	out := []string{}
	for _, w := range wordRe.FindAllString(strings.ToLower(s), -1) {
		// Keep apostrophes and hyphens inside Latin words; drop stray ones.
		if w = strings.Trim(w, "'’-"); w != "" {
			out = append(out, w)
		}
	}
	return out
}

// countWritingWords counts words, and each Chinese/Japanese character as one.
func countWritingWords(s string) int { return len(writingWords(s)) }

// copiesText reports whether answer is (largely) a copy of source: it
// contains source outright, or most of its three-word runs appear in source.
// Short answers that merely share common phrases ("ich bin", "es una") pass.
func copiesText(answer, source string) bool {
	a, s := writingWords(answer), writingWords(source)
	if len(a) == 0 || len(s) < 3 {
		return false
	}
	aj, sj := " "+strings.Join(a, " ")+" ", " "+strings.Join(s, " ")+" "
	if strings.Contains(aj, sj) {
		return true // the whole source pasted in, perhaps with extra words
	}
	if len(a) >= 5 && strings.Contains(sj, aj) {
		return true // a substantial chunk of the source and nothing else
	}
	if len(a) < 4 {
		return false
	}
	grams := map[string]bool{}
	for i := 0; i+3 <= len(s); i++ {
		grams[strings.Join(s[i:i+3], " ")] = true
	}
	total, hit := 0, 0
	for i := 0; i+3 <= len(a); i++ {
		total++
		if grams[strings.Join(a[i:i+3], " ")] {
			hit++
		}
	}
	return total > 0 && float64(hit)/float64(total) >= 0.6
}

// isRepetitive catches padding: very few distinct words, or one sentence
// repeated three or more times.
func isRepetitive(text string) bool {
	w := writingWords(text)
	if len(w) >= 12 {
		uniq := map[string]bool{}
		for _, x := range w {
			uniq[x] = true
		}
		if float64(len(uniq))/float64(len(w)) < 0.35 {
			return true
		}
	}
	seen := map[string]int{}
	for _, s := range regexp.MustCompile(`[.!?。！？\n]+`).Split(strings.ToLower(text), -1) {
		s = strings.Join(writingWords(s), " ")
		if len(strings.Fields(s)) >= 2 {
			seen[s]++
			if seen[s] >= 3 {
				return true
			}
		}
	}
	return false
}

var lengthRange = regexp.MustCompile(`(\d+)\s*(?:[–-]\s*\d+\s*|\+\s*)?(?:words|characters)`)
var approxLength = regexp.MustCompile(`[≈~]\s*(\d+)\s*(?:words|characters)`)

// minWordsFor reads the length a task asks for ("30–40 words", "≈150 words",
// "250+ words", "around 200 words", "about 80 characters" for Chinese and
// Japanese, counted per character) and requires 60% of its lower bound — a
// lesson is practice, not the exam — and never fewer than 12 words.
func minWordsFor(task string) int {
	t := strings.ToLower(task)
	n := 0
	if m := approxLength.FindStringSubmatch(t); m != nil {
		n, _ = strconv.Atoi(m[1])
	} else if m := lengthRange.FindStringSubmatch(t); m != nil {
		n, _ = strconv.Atoi(m[1])
	}
	if req := int(math.Round(float64(n) * 0.6)); req > 12 {
		return req
	}
	return 12
}

// lessonLanguage returns the language and CEFR level of a lesson's skill.
func lessonLanguage(lessonID uint) (lang, level string) {
	var lesson models.Lesson
	if database.DB.First(&lesson, lessonID).Error != nil {
		return "", ""
	}
	var skill models.Skill
	if database.DB.First(&skill, lesson.SkillID).Error != nil {
		return "", ""
	}
	return skill.Language, unitLevel(skill.Unit)
}

// --- typed answers ---------------------------------------------------------

type checkAnswerInput struct {
	LessonID uint   `json:"lessonId"`
	Question string `json:"question"`
	Expected string `json:"expected"`
	Answer   string `json:"answer"`
}

// CheckAnswer is a second opinion on a typed answer the app's exact comparison
// marked wrong (POST /lessons/check-answer). The clients only ask about
// multi-word answers, where valid alternative wordings are common.
//
// Returns {available:false} when the coach isn't configured or fails — the
// client then keeps its own verdict.
func (l *LessonController) CheckAnswer(c *fiber.Ctx) error {
	var in checkAnswerInput
	if err := c.BodyParser(&in); err != nil || strings.TrimSpace(in.Answer) == "" || in.Expected == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
	}
	lang, _ := lessonLanguage(in.LessonID)
	ctx, cancel := context.WithCancel(c.Context())
	defer cancel()
	j, err := coach.JudgeAnswer(ctx, lang, in.Question, in.Expected, in.Answer)
	if err != nil {
		return c.JSON(fiber.Map{"available": false})
	}
	return c.JSON(fiber.Map{"available": true, "correct": j.Correct, "explanation": j.Explanation})
}
