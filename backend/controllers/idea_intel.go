package controllers

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"

	"lumora/backend/database"
	"lumora/backend/middleware"
	"lumora/backend/models"
)

// ============================================================================
// Duplicate detection and thread summarisation
// ============================================================================
//
// Two features the spec calls "AI": "Similar to Idea #3 — view thread?" and a
// one-line summary of a long thread.
//
// Both are implemented here statistically, not with a language model: TF-IDF
// cosine similarity for the first, and extractive summarisation (pick the
// sentences that carry the most of the thread's distinctive vocabulary) for the
// second. That choice is deliberate and worth stating plainly —
//
//   - it runs in single-digit milliseconds with no API key, no per-call cost,
//     no rate limit and no network dependency, which matters on a free-tier host
//     that sleeps;
//   - nothing a team writes in a private idea thread leaves the server;
//   - it degrades honestly. An extractive summary is always a real sentence
//     somebody actually wrote, so it can be unhelpful but never fabricated —
//     unlike a generated summary, which fails by inventing agreement that was
//     never reached.
//
// Both entry points (similarTo, summariseThread) are single functions with
// plain inputs and outputs, so swapping in a call to a language model later is
// a contained change if the quality ceiling here turns out to be too low.

// Words carrying no topical signal. Kept deliberately short: an aggressive
// stop-list on a small corpus throws away the words that actually distinguish
// two ideas.
var stopWords = map[string]bool{
	"the": true, "a": true, "an": true, "and": true, "or": true, "but": true,
	"if": true, "then": true, "to": true, "of": true, "in": true, "on": true,
	"for": true, "with": true, "at": true, "by": true, "from": true, "as": true,
	"is": true, "are": true, "was": true, "were": true, "be": true, "been": true,
	"it": true, "its": true, "this": true, "that": true, "these": true,
	"those": true, "we": true, "our": true, "you": true, "your": true,
	"i": true, "my": true, "they": true, "their": true, "he": true, "she": true,
	"do": true, "does": true, "did": true, "can": true, "could": true,
	"would": true, "should": true, "will": true, "just": true, "not": true,
	"so": true, "have": true, "has": true, "had": true, "what": true,
	"how": true, "why": true, "when": true, "there": true, "here": true,
	"about": true, "into": true, "than": true, "some": true, "any": true,
}

// tokenise lowercases, strips punctuation and drops stop words and very short
// tokens. Words are truncated to a crude stem so "tagging"/"tags"/"tag" match.
func tokenise(s string) []string {
	fields := strings.FieldsFunc(strings.ToLower(s), func(r rune) bool {
		return !(r >= 'a' && r <= 'z') && !(r >= '0' && r <= '9')
	})
	out := make([]string, 0, len(fields))
	for _, f := range fields {
		if len(f) < 3 || stopWords[f] {
			continue
		}
		out = append(out, stem(f))
	}
	return out
}

// stem is a small suffix chopper — enough to collapse the plural and gerund
// forms that make near-identical ideas look different, without the weight of a
// real stemmer.
//
// The doubled-consonant rule matters more than it looks: English doubles the
// final consonant before "-ing", so without it "tagging" stems to "tagg" and
// fails to match "tag" — which is exactly the pair a duplicate-detector on an
// auto-tagging board needs to catch.
func stem(w string) string {
	for _, suffix := range []string{"ing", "ers", "er", "ies", "es", "s"} {
		if len(w) > len(suffix)+3 && strings.HasSuffix(w, suffix) {
			base := strings.TrimSuffix(w, suffix)
			if suffix == "ing" {
				base = undouble(base)
			}
			return base
		}
	}
	return w
}

// undouble collapses a trailing doubled consonant ("tagg" -> "tag"). Vowels are
// left alone: "seeing" -> "see" must not become "se".
func undouble(w string) string {
	n := len(w)
	if n < 3 || w[n-1] != w[n-2] {
		return w
	}
	switch w[n-1] {
	case 'a', 'e', 'i', 'o', 'u':
		return w
	}
	return w[:n-1]
}

func termFreq(tokens []string) map[string]float64 {
	tf := map[string]float64{}
	for _, t := range tokens {
		tf[t]++
	}
	return tf
}

// cosine similarity over IDF-weighted term vectors. IDF is what stops every
// idea on a language-learning board scoring 0.4 against every other one purely
// because they all say "lesson" and "learner".
func cosineIDF(a, b map[string]float64, idf map[string]float64) float64 {
	var dot, na, nb float64
	for t, av := range a {
		w := idf[t]
		na += (av * w) * (av * w)
		if bv, ok := b[t]; ok {
			dot += (av * w) * (bv * w)
		}
	}
	for t, bv := range b {
		w := idf[t]
		nb += (bv * w) * (bv * w)
	}
	if na == 0 || nb == 0 {
		return 0
	}
	return dot / (math.Sqrt(na) * math.Sqrt(nb))
}

// buildIDF computes inverse document frequency across a corpus of documents.
func buildIDF(docs [][]string) map[string]float64 {
	n := float64(len(docs))
	seen := map[string]float64{}
	for _, doc := range docs {
		uniq := map[string]bool{}
		for _, t := range doc {
			uniq[t] = true
		}
		for t := range uniq {
			seen[t]++
		}
	}
	idf := map[string]float64{}
	for t, df := range seen {
		idf[t] = math.Log((n+1)/(df+1)) + 1
	}
	return idf
}

// --- duplicate detection -----------------------------------------------------

// similarityFloor is the point below which a "similar idea" suggestion is more
// annoying than useful. Tuned to be conservative: a missed duplicate costs a
// little redundancy, a false one trains people to ignore the prompt.
const similarityFloor = 0.18

// similarTo finds the ideas closest to the given one. Used both for the "you
// might be duplicating this" prompt while typing a new idea, and for the
// "related" list in the details panel.
func (ic *IdeaController) similarTo(idea models.Idea, limit int) []fiber.Map {
	return findSimilar(idea.Title+" "+idea.Description, idea.ID, limit)
}

func findSimilar(text string, excludeID uint, limit int) []fiber.Map {
	query := tokenise(text)
	if len(query) == 0 {
		return []fiber.Map{}
	}

	var ideas []models.Idea
	database.DB.Where("merged_into_id IS NULL").Limit(500).Find(&ideas)

	docs := make([][]string, 0, len(ideas)+1)
	tokensByIdea := make(map[uint][]string, len(ideas))
	for _, i := range ideas {
		tk := tokenise(i.Title + " " + i.Description)
		tokensByIdea[i.ID] = tk
		docs = append(docs, tk)
	}
	docs = append(docs, query)
	idf := buildIDF(docs)
	qv := termFreq(query)

	type scored struct {
		idea models.Idea
		sim  float64
	}
	ranked := make([]scored, 0, len(ideas))
	for _, i := range ideas {
		if i.ID == excludeID {
			continue
		}
		sim := cosineIDF(qv, termFreq(tokensByIdea[i.ID]), idf)
		if sim >= similarityFloor {
			ranked = append(ranked, scored{i, sim})
		}
	}
	sort.SliceStable(ranked, func(a, b int) bool { return ranked[a].sim > ranked[b].sim })

	out := make([]fiber.Map, 0, limit)
	for i, r := range ranked {
		if i >= limit {
			break
		}
		out = append(out, fiber.Map{
			"id": r.idea.ID, "title": r.idea.Title, "status": r.idea.Status,
			"score": r.idea.Score, "messageCount": r.idea.MessageCount,
			"similarity": math.Round(r.sim*100) / 100,
		})
	}
	return out
}

// Similar backs the live "this looks like Idea #3" prompt in the composer.
func (ic *IdeaController) Similar(c *fiber.Ctx) error {
	text := strings.TrimSpace(c.Query("q"))
	if len(text) < 8 {
		return c.JSON(fiber.Map{"similar": []fiber.Map{}})
	}
	return c.JSON(fiber.Map{"similar": findSimilar(text, 0, 3)})
}

// --- thread summarisation ----------------------------------------------------

// Summary condenses a long thread: a one-line gist, the points that carried the
// most of the discussion, the open questions, and who took part.
func (ic *IdeaController) Summary(c *fiber.Ctx) error {
	user := middleware.CurrentUser(c)

	idea, ok := ic.load(c)
	if !ok {
		return nil // the helper already wrote the error response
	}

	var msgs []models.IdeaMessage
	database.DB.Where("idea_id = ? AND deleted_at IS NULL", idea.ID).
		Order("created_at asc").Find(&msgs)

	sum := summariseThread(*idea, msgs)
	sum["viewer"] = user.ID
	return c.JSON(sum)
}

// summariseThread is extractive: it scores every sentence by how much of the
// thread's distinctive vocabulary it carries, and returns the strongest ones.
// Nothing is generated, so nothing can be invented.
func summariseThread(idea models.Idea, msgs []models.IdeaMessage) fiber.Map {
	type sentence struct {
		text   string
		author uint
		at     time.Time
		score  float64
	}

	var sentences []sentence
	docs := [][]string{tokenise(idea.Title + " " + idea.Description)}
	questions := []fiber.Map{}

	for _, m := range msgs {
		if m.Kind != models.MsgText && m.Kind != models.MsgCode {
			continue
		}
		for _, raw := range splitSentences(m.Body) {
			s := strings.TrimSpace(raw)
			if len(s) < 15 {
				continue
			}
			tk := tokenise(s)
			if len(tk) == 0 {
				continue
			}
			docs = append(docs, tk)
			sentences = append(sentences, sentence{text: s, author: m.AuthorID, at: m.CreatedAt})

			// Unanswered questions are the most actionable thing in a long
			// thread, so they're pulled out separately rather than competing
			// with statements for a summary slot.
			if strings.HasSuffix(s, "?") {
				questions = append(questions, fiber.Map{
					"text": s, "authorId": m.AuthorID,
					"at": m.CreatedAt.Format(time.RFC3339),
				})
			}
		}
	}

	if len(sentences) == 0 {
		gist := idea.Description
		if strings.TrimSpace(gist) == "" {
			gist = idea.Title
		}
		return fiber.Map{
			"gist": snippet(gist), "keyPoints": []fiber.Map{}, "questions": questions,
			"messageCount": len(msgs), "generated": false,
		}
	}

	idf := buildIDF(docs)
	for i := range sentences {
		tk := tokenise(sentences[i].text)
		var score float64
		for t, f := range termFreq(tk) {
			score += f * idf[t]
		}
		// Normalise by length so a rambling paragraph doesn't automatically win
		// over a sharp one-liner.
		sentences[i].score = score / math.Sqrt(float64(len(tk))+1)
	}

	ranked := make([]sentence, len(sentences))
	copy(ranked, sentences)
	sort.SliceStable(ranked, func(a, b int) bool { return ranked[a].score > ranked[b].score })

	limit := 3
	if len(ranked) < limit {
		limit = len(ranked)
	}
	top := ranked[:limit]
	// Present them in the order they were said, so the summary reads as a
	// narrative rather than a leaderboard of sentences.
	sort.SliceStable(top, func(a, b int) bool { return top[a].at.Before(top[b].at) })

	keyPoints := make([]fiber.Map, 0, limit)
	for _, s := range top {
		var u models.User
		database.DB.First(&u, s.author)
		keyPoints = append(keyPoints, fiber.Map{
			"text": s.text, "author": toChatUser(u),
			"at": s.at.Format(time.RFC3339),
		})
	}

	if len(questions) > 3 {
		questions = questions[:3]
	}

	return fiber.Map{
		"gist":         gistLine(idea, msgs, ranked[0].text),
		"keyPoints":    keyPoints,
		"questions":    questions,
		"messageCount": len(msgs),
		"generated":    true,
	}
}

// gistLine is the one-liner at the top of the summary: what the thread is
// about, how busy it is, and its single strongest sentence.
func gistLine(idea models.Idea, msgs []models.IdeaMessage, best string) string {
	people := map[uint]bool{}
	images := 0
	for _, m := range msgs {
		people[m.AuthorID] = true
		if m.Kind == models.MsgImage {
			images++
		}
	}
	parts := []string{fmt.Sprintf("%d messages from %d people", len(msgs), len(people))}
	if images > 0 {
		parts = append(parts, fmt.Sprintf("%d attachments", images))
	}
	return fmt.Sprintf("%s — %s. Strongest point so far: \"%s\"",
		idea.Title, strings.Join(parts, ", "), snippet(best))
}

// splitSentences breaks text on terminators while keeping the terminator, so an
// extracted question still reads as a question.
func splitSentences(s string) []string {
	var out []string
	var cur strings.Builder
	for _, r := range s {
		cur.WriteRune(r)
		if r == '.' || r == '!' || r == '?' || r == '\n' {
			out = append(out, cur.String())
			cur.Reset()
		}
	}
	if strings.TrimSpace(cur.String()) != "" {
		out = append(out, cur.String())
	}
	return out
}
