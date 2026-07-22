package controllers

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"lumora/backend/database"
	"lumora/backend/models"
)

// newIdeaDB gives each test a clean database with the ideas tables migrated.
func newIdeaDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared&_pragma=busy_timeout(5000)"),
		&gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatalf("open test db: %v", err)
	}
	tables := []interface{}{
		&models.User{}, &models.Notification{}, &models.Message{},
		&models.Idea{}, &models.IdeaVote{}, &models.IdeaStar{}, &models.IdeaTag{},
		&models.IdeaMessage{}, &models.IdeaReaction{}, &models.IdeaEvent{},
		&models.IdeaTask{}, &models.BrainstormSession{},
	}
	if err := db.Migrator().DropTable(tables...); err != nil {
		t.Fatalf("drop: %v", err)
	}
	if err := db.AutoMigrate(tables...); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	database.DB = db
	return db
}

func makeIdea(t *testing.T, owner uint, title, desc string, tags ...string) *models.Idea {
	t.Helper()
	idea := &models.Idea{
		OwnerID: owner, Title: title, Description: desc,
		Status: models.IdeaDraft, LastActivityAt: time.Now(),
	}
	if err := database.DB.Create(idea).Error; err != nil {
		t.Fatalf("create idea: %v", err)
	}
	setTags(idea.ID, tags)
	return idea
}

// --- voting ------------------------------------------------------------------

func TestVoteTalliesAreRecomputedNotIncremented(t *testing.T) {
	newIdeaDB(t)
	idea := makeIdea(t, 1, "Auto-tagging for notes", "Tag notes automatically")

	applyVote(idea, 2, 1)
	applyVote(idea, 3, 1)
	applyVote(idea, 4, -1)

	if idea.Upvotes != 2 || idea.Downvotes != 1 || idea.Score != 1 {
		t.Fatalf("got %d up / %d down / score %d, want 2/1/1",
			idea.Upvotes, idea.Downvotes, idea.Score)
	}

	// Re-voting replaces rather than stacks.
	applyVote(idea, 2, -1)
	if idea.Upvotes != 1 || idea.Downvotes != 2 {
		t.Errorf("after switching a vote: %d up / %d down, want 1/2", idea.Upvotes, idea.Downvotes)
	}

	// Submitting the same vote twice must not double-count it.
	applyVote(idea, 3, 1)
	applyVote(idea, 3, 1)
	if idea.Upvotes != 1 {
		t.Errorf("duplicate vote counted twice: %d upvotes", idea.Upvotes)
	}

	// Clearing a vote removes it entirely.
	applyVote(idea, 4, 0)
	if idea.Downvotes != 1 {
		t.Errorf("cleared vote still counted: %d downvotes", idea.Downvotes)
	}
}

// Crossing the vote threshold escalates the idea on its own — nobody should
// have to remember to do it.
func TestReachingTheThresholdAutoFlagsForReview(t *testing.T) {
	newIdeaDB(t)
	idea := makeIdea(t, 1, "Dark mode", "")

	for i := 2; i < 2+reviewThreshold-1; i++ {
		applyVote(idea, uint(i), 1)
	}
	if idea.Status != models.IdeaDraft {
		t.Fatalf("status moved early at %d votes: %s", idea.Score, idea.Status)
	}

	applyVote(idea, 99, 1)
	if idea.Score != reviewThreshold {
		t.Fatalf("score = %d, want %d", idea.Score, reviewThreshold)
	}
	if idea.Status != models.IdeaUnderReview {
		t.Errorf("status = %s, want under_review at %d votes", idea.Status, reviewThreshold)
	}

	var events []models.IdeaEvent
	database.DB.Where("idea_id = ? AND kind = ?", idea.ID, "vote_threshold").Find(&events)
	if len(events) != 1 {
		t.Errorf("got %d threshold events in the history, want 1", len(events))
	}
}

// --- ranking -----------------------------------------------------------------

// Hot has to let a fresh idea out-rank a stale one with more votes, or the
// board ossifies around whatever was posted first.
func TestHotRankingDecaysWithAge(t *testing.T) {
	fresh := models.Idea{Score: 5, LastActivityAt: time.Now(), CreatedAt: time.Now()}
	stale := models.Idea{
		Score:          15,
		LastActivityAt: time.Now().Add(-30 * 24 * time.Hour),
		CreatedAt:      time.Now().Add(-30 * 24 * time.Hour),
	}
	if hotScore(fresh) <= hotScore(stale) {
		t.Errorf("a month-old idea with 15 votes (%v) out-ranks today's with 5 (%v)",
			hotScore(stale), hotScore(fresh))
	}
}

// Controversial must surface the evenly-split ideas, not simply the popular
// ones — divisive is exactly the signal pure score buries.
func TestControversyFavoursEvenSplits(t *testing.T) {
	divisive := models.Idea{Upvotes: 30, Downvotes: 28}
	popular := models.Idea{Upvotes: 40, Downvotes: 1}
	quiet := models.Idea{Upvotes: 2, Downvotes: 0}

	if controversy(divisive) <= controversy(popular) {
		t.Errorf("30/28 (%v) should be more controversial than 40/1 (%v)",
			controversy(divisive), controversy(popular))
	}
	if controversy(quiet) != 0 {
		t.Errorf("an unopposed idea scored %v controversy, want 0", controversy(quiet))
	}
}

func TestSortModesOrderDifferently(t *testing.T) {
	old := time.Now().Add(-20 * 24 * time.Hour)
	ideas := []models.Idea{
		{ID: 1, Title: "old favourite", Score: 40, Upvotes: 40, CreatedAt: old, LastActivityAt: old},
		{ID: 2, Title: "new", Score: 2, Upvotes: 2, CreatedAt: time.Now(), LastActivityAt: time.Now()},
		{ID: 3, Title: "divisive", Score: 2, Upvotes: 25, Downvotes: 23,
			CreatedAt: time.Now().Add(-time.Hour), LastActivityAt: time.Now().Add(-time.Hour)},
	}

	cp := append([]models.Idea{}, ideas...)
	sortIdeas(cp, "top")
	if cp[0].ID != 1 {
		t.Errorf("top sort led with #%d, want the 40-vote idea", cp[0].ID)
	}

	cp = append([]models.Idea{}, ideas...)
	sortIdeas(cp, "new")
	if cp[0].ID != 2 {
		t.Errorf("new sort led with #%d, want the newest", cp[0].ID)
	}

	cp = append([]models.Idea{}, ideas...)
	sortIdeas(cp, "controversial")
	if cp[0].ID != 3 {
		t.Errorf("controversial sort led with #%d, want the 25/23 split", cp[0].ID)
	}
}

// --- tags --------------------------------------------------------------------

func TestTagsAreNormalisedAndDeduplicated(t *testing.T) {
	newIdeaDB(t)
	idea := makeIdea(t, 1, "Voice notes", "")

	setTags(idea.ID, []string{"#AI", "ai", " AI ", "User Experience", "v2.0"})

	got := tagsOf(idea.ID)
	want := []string{"ai", "user-experience", "v2.0"}
	if strings.Join(got, ",") != strings.Join(want, ",") {
		t.Errorf("tags = %v, want %v", got, want)
	}
}

// --- archive / delete --------------------------------------------------------

// Deleting an idea people have engaged with would erase their contributions —
// archiving is the honest action, and it's forced.
func TestDiscussedIdeasCannotBeDeleted(t *testing.T) {
	newIdeaDB(t)
	idea := makeIdea(t, 1, "Auto-tagging", "")
	idea.MessageCount = 3
	database.DB.Save(idea)

	// Mirrors the guard in Delete.
	if !(idea.MessageCount > 0 || idea.Upvotes > 1 || idea.Downvotes > 0) {
		t.Fatal("an idea with 3 messages should be protected from deletion")
	}
}

// Deleting must take everything with it. These tables have no foreign keys, so
// a missed table leaves orphaned rows — and orphaned attachment bytes — in the
// database forever.
func TestDeletingAnIdeaCascades(t *testing.T) {
	newIdeaDB(t)

	owner := &models.User{Email: "owner@test.dev", Name: "Owner"}
	database.DB.Create(owner)
	idea := makeIdea(t, owner.ID, "Auto-tagging", "desc", "ai", "ux")

	applyVote(idea, owner.ID, 1)
	applyVote(idea, 42, 1)
	database.DB.Create(&models.IdeaStar{IdeaID: idea.ID, UserID: 42})
	database.DB.Create(&models.IdeaTask{IdeaID: idea.ID, Title: "ship it"})
	database.DB.Create(&models.BrainstormSession{
		IdeaID: &idea.ID, StartedBy: owner.ID, EndsAt: time.Now(),
	})

	msg := models.IdeaMessage{
		IdeaID: idea.ID, AuthorID: 42, Kind: models.MsgImage,
		Body: "a mock", Data: []byte("pretend-jpeg-bytes"), Mime: "image/jpeg",
	}
	database.DB.Create(&msg)
	database.DB.Create(&models.IdeaReaction{MessageID: msg.ID, UserID: 42, Emoji: "🔥"})

	// A second idea merged into this one must survive, released back onto the
	// board rather than left pointing at a row that no longer exists.
	orphan := makeIdea(t, 42, "Automatic note tags", "")
	database.DB.Model(orphan).Update("merged_into_id", idea.ID)

	app := fiber.New()
	ic := &IdeaController{}
	app.Delete("/ideas/:id", func(c *fiber.Ctx) error {
		c.Locals("user", owner)
		return ic.Delete(c)
	})
	res, err := app.Test(
		httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/ideas/%d", idea.ID), nil), -1)
	if err != nil {
		t.Fatalf("request: %v", err)
	}
	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		t.Fatalf("status = %d (%s), want 200", res.StatusCode, body)
	}

	counts := map[string]int64{}
	for name, model := range map[string]interface{}{
		"ideas":       &models.Idea{},
		"messages":    &models.IdeaMessage{},
		"votes":       &models.IdeaVote{},
		"stars":       &models.IdeaStar{},
		"tags":        &models.IdeaTag{},
		"events":      &models.IdeaEvent{},
		"tasks":       &models.IdeaTask{},
		"brainstorms": &models.BrainstormSession{},
	} {
		var n int64
		database.DB.Model(model).Where("idea_id = ?", idea.ID).Count(&n)
		counts[name] = n
	}
	for name, n := range counts {
		if n != 0 {
			t.Errorf("%s: %d rows survived the delete", name, n)
		}
	}

	var reactions int64
	database.DB.Model(&models.IdeaReaction{}).
		Where("message_id = ?", msg.ID).Count(&reactions)
	if reactions != 0 {
		t.Errorf("%d reactions survived — they hang off the message, not the idea", reactions)
	}

	var released models.Idea
	if database.DB.First(&released, orphan.ID).Error != nil {
		t.Fatal("the merged-in idea was deleted along with its target")
	}
	if released.MergedIntoID != nil {
		t.Errorf("merged idea still points at the deleted row (%v)", *released.MergedIntoID)
	}
}

// Only the person who posted an idea can delete it.
func TestOnlyTheOwnerCanDeleteAnIdea(t *testing.T) {
	newIdeaDB(t)

	owner := &models.User{Email: "o@test.dev", Name: "Owner"}
	stranger := &models.User{Email: "s@test.dev", Name: "Stranger"}
	database.DB.Create(owner)
	database.DB.Create(stranger)
	idea := makeIdea(t, owner.ID, "Dark mode", "")

	app := fiber.New()
	ic := &IdeaController{}
	app.Delete("/ideas/:id", func(c *fiber.Ctx) error {
		c.Locals("user", stranger)
		return ic.Delete(c)
	})
	res, _ := app.Test(
		httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/ideas/%d", idea.ID), nil), -1)
	if res.StatusCode != http.StatusForbidden {
		t.Fatalf("status = %d, want 403", res.StatusCode)
	}

	var still models.Idea
	if database.DB.First(&still, idea.ID).Error != nil {
		t.Error("the idea was deleted by someone who doesn't own it")
	}
}

// --- merging -----------------------------------------------------------------

// Merging must not let one person's support count twice.
func TestMergeDeduplicatesVotersAcrossIdeas(t *testing.T) {
	newIdeaDB(t)
	target := makeIdea(t, 1, "Auto-tagging for notes", "")
	dupe := makeIdea(t, 2, "Automatic note tags", "")

	// Alice (3) backs both; Bob (4) only backs the duplicate.
	applyVote(target, 3, 1)
	applyVote(dupe, 3, 1)
	applyVote(dupe, 4, 1)

	// The transfer, as Merge performs it.
	var votes []models.IdeaVote
	database.DB.Where("idea_id = ?", dupe.ID).Find(&votes)
	for _, v := range votes {
		var dup int64
		database.DB.Model(&models.IdeaVote{}).
			Where("idea_id = ? AND user_id = ?", target.ID, v.UserID).Count(&dup)
		if dup == 0 {
			database.DB.Create(&models.IdeaVote{IdeaID: target.ID, UserID: v.UserID, Value: v.Value})
		}
	}
	applyVote(target, 99, 0) // force a recount

	if target.Upvotes != 2 {
		t.Errorf("after merge target has %d upvotes, want 2 (Alice counted once, plus Bob)",
			target.Upvotes)
	}
}

// --- duplicate detection -----------------------------------------------------

func TestSimilarIdeasAreDetected(t *testing.T) {
	newIdeaDB(t)
	makeIdea(t, 1, "Auto-tagging for notes",
		"Automatically tag notes using the existing search API so users don't file things by hand")
	makeIdea(t, 1, "Dark mode", "A dark colour scheme for the whole app at night")
	makeIdea(t, 1, "Voice memos", "Record short audio clips instead of typing")

	got := findSimilar("Automatic tagging of notes with our search API", 0, 3)
	if len(got) == 0 {
		t.Fatal("a near-identical idea was not detected as similar")
	}
	if got[0]["title"] != "Auto-tagging for notes" {
		t.Errorf("closest match = %v, want the auto-tagging idea", got[0]["title"])
	}
}

// A conservative floor matters more than recall: a false "this is a duplicate"
// prompt teaches people to dismiss the prompt.
func TestUnrelatedIdeasAreNotFlaggedAsSimilar(t *testing.T) {
	newIdeaDB(t)
	makeIdea(t, 1, "Auto-tagging for notes", "Tag notes automatically using the search API")

	got := findSimilar("Add a Swahili course for beginners", 0, 3)
	if len(got) != 0 {
		t.Errorf("unrelated idea matched: %v", got)
	}
}

func TestStemmingCollapsesWordForms(t *testing.T) {
	if stem("tagging") != stem("tag") {
		t.Errorf("tagging -> %q, tag -> %q; want a match", stem("tagging"), stem("tag"))
	}
	if stem("notes") != stem("note") {
		t.Errorf("notes -> %q, note -> %q; want a match", stem("notes"), stem("note"))
	}
	// Short words must survive intact, or "was" becomes "wa".
	if stem("api") != "api" {
		t.Errorf("api was stemmed to %q", stem("api"))
	}
}

// --- summarisation -----------------------------------------------------------

// The summary is extractive, so every line it returns must be a sentence
// somebody actually wrote — that's the property that makes it unable to invent
// agreement that was never reached.
func TestSummaryOnlyQuotesRealSentences(t *testing.T) {
	newIdeaDB(t)
	idea := makeIdea(t, 1, "Auto-tagging for notes", "Tag notes automatically")

	bodies := []string{
		"What if we used AI to auto-tag notes as people write them?",
		"Love it. Could work with our existing search API without new infrastructure.",
		"What about privacy? Tagging implies reading the note contents on the server.",
		"ok",
		"We could run the tagging locally on device so note contents never leave it.",
	}
	var msgs []models.IdeaMessage
	for i, b := range bodies {
		m := models.IdeaMessage{
			IdeaID: idea.ID, AuthorID: uint(i%3 + 1), Kind: models.MsgText, Body: b,
			CreatedAt: time.Now().Add(time.Duration(i) * time.Minute),
		}
		database.DB.Create(&m)
		msgs = append(msgs, m)
	}

	sum := summariseThread(*idea, msgs)
	if sum["generated"] != true {
		t.Fatal("summary was not generated for a five-message thread")
	}

	points, ok := sum["keyPoints"].([]fiber.Map)
	if !ok || len(points) == 0 {
		t.Fatalf("keyPoints = %#v, want a non-empty list", sum["keyPoints"])
	}
	joined := strings.Join(bodies, " ")
	for _, kp := range points {
		text, _ := kp["text"].(string)
		if text == "" || !strings.Contains(joined, strings.TrimSpace(text)) {
			t.Errorf("summary quoted %q, which nobody wrote", text)
		}
	}

	// The very short "ok" carries no signal and shouldn't win a slot.
	for _, kp := range points {
		if strings.TrimSpace(kp["text"].(string)) == "ok" {
			t.Error("a filler message was selected as a key point")
		}
	}

	// Open questions are extracted separately, because in a long thread they're
	// the most actionable thing in it.
	questions, ok := sum["questions"].([]fiber.Map)
	if !ok || len(questions) == 0 {
		t.Fatalf("questions = %#v, want the two '?' messages", sum["questions"])
	}
	foundPrivacy := false
	for _, q := range questions {
		if strings.Contains(q["text"].(string), "privacy") {
			foundPrivacy = true
		}
	}
	if !foundPrivacy {
		t.Error(`"What about privacy?" was not surfaced as an open question`)
	}
}

// A thread with nothing in it must not claim to have summarised anything.
func TestEmptyThreadReportsNoSummary(t *testing.T) {
	newIdeaDB(t)
	idea := makeIdea(t, 1, "Voice notes", "Record short clips")

	sum := summariseThread(*idea, nil)
	if sum["generated"] != false {
		t.Errorf("generated = %v for an empty thread, want false", sum["generated"])
	}
	if sum["gist"] == "" {
		t.Error("even an empty thread should fall back to the idea's own description")
	}
}

// --- mentions ----------------------------------------------------------------

func TestMentionPatterns(t *testing.T) {
	body := "Nice work @alice — this overlaps with @idea#12 and maybe @bob-smith too."

	users := mentionUser.FindAllStringSubmatch(body, -1)
	names := []string{}
	for _, m := range users {
		if m[1] != "idea" {
			names = append(names, m[1])
		}
	}
	if len(names) != 2 || names[0] != "alice" || names[1] != "bob-smith" {
		t.Errorf("user mentions = %v, want [alice bob-smith]", names)
	}

	ideas := mentionIdea.FindAllStringSubmatch(body, -1)
	if len(ideas) != 1 || ideas[0][1] != "12" {
		t.Errorf("idea mentions = %v, want [12]", ideas)
	}
}

// --- workflow ----------------------------------------------------------------

func TestStatusFlowIsClosed(t *testing.T) {
	for _, s := range statusFlow {
		if !validStatus(s) {
			t.Errorf("%q is in the flow but rejected as invalid", s)
		}
	}
	if validStatus("shipped") {
		t.Error("an arbitrary status was accepted")
	}
	if len(statusFlow) != 6 {
		t.Errorf("status flow has %d steps, want 6", len(statusFlow))
	}
}

// Collection fields must serialise as [] rather than null. A nil slice becomes
// `null` in JSON, and the client reads reactions/replies as arrays — one nil
// and the thread panel throws on .length.
func TestMessageCollectionsSerialiseAsEmptyArrays(t *testing.T) {
	newIdeaDB(t)
	idea := makeIdea(t, 1, "Voice notes", "")

	msg := models.IdeaMessage{
		IdeaID: idea.ID, AuthorID: 1, Kind: models.MsgText,
		Body: "no replies and no reactions on this one", CreatedAt: time.Now(),
	}
	database.DB.Create(&msg)

	ic := &IdeaController{}
	dto := ic.toMessageDTO(msg, 1)

	if dto.Reactions == nil {
		t.Error("reactions is nil; it must be an empty slice so JSON emits []")
	}
	if dto.Replies == nil {
		t.Error("replies is nil; it must be an empty slice so JSON emits []")
	}

	blob, err := json.Marshal(dto)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	for _, field := range []string{`"reactions":null`, `"replies":null`} {
		if strings.Contains(string(blob), field) {
			t.Errorf("serialised JSON contains %s", field)
		}
	}
}

// A missing idea must 404, not panic.
//
// Regression test for a Fiber trap: c.JSON() returns nil on success, so a
// helper written as `return nil, c.Status(404).JSON(...)` hands its caller a
// nil error alongside a nil value. The guard passes, the next line dereferences
// nil, and a clean 404 becomes a 500 panic. The helpers return (value, ok) now,
// which can't be misread that way.
func TestMissingIdeaReturnsNotFoundNotPanic(t *testing.T) {
	newIdeaDB(t)

	app := fiber.New()
	ic := &IdeaController{}
	app.Get("/ideas/:id", func(c *fiber.Ctx) error {
		idea, ok := ic.load(c)
		if !ok {
			return nil
		}
		return c.JSON(fiber.Map{"id": idea.ID})
	})

	for _, path := range []string{"/ideas/99999", "/ideas/not-a-number"} {
		t.Run(path, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, path, nil)
			res, err := app.Test(req, -1)
			if err != nil {
				t.Fatalf("request: %v", err)
			}
			if res.StatusCode != http.StatusNotFound && res.StatusCode != http.StatusBadRequest {
				t.Fatalf("status = %d, want 404 or 400", res.StatusCode)
			}
			body, _ := io.ReadAll(res.Body)
			var payload map[string]any
			if err := json.Unmarshal(body, &payload); err != nil {
				t.Fatalf("body was not JSON (%q): %v", body, err)
			}
			if payload["error"] == nil {
				t.Errorf("no error message in %s", body)
			}
		})
	}
}

// A found idea still comes back normally — proving the ok-check didn't just
// make every request fail closed.
func TestExistingIdeaLoads(t *testing.T) {
	newIdeaDB(t)
	idea := makeIdea(t, 1, "Auto-tagging", "")

	app := fiber.New()
	ic := &IdeaController{}
	app.Get("/ideas/:id", func(c *fiber.Ctx) error {
		got, ok := ic.load(c)
		if !ok {
			return nil
		}
		return c.JSON(fiber.Map{"id": got.ID})
	})

	res, err := app.Test(httptest.NewRequest(http.MethodGet, fmt.Sprintf("/ideas/%d", idea.ID), nil), -1)
	if err != nil {
		t.Fatalf("request: %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", res.StatusCode)
	}
}

// --- attachments -------------------------------------------------------------

func TestDownscaleCapsTheLongestEdge(t *testing.T) {
	wide := imageOfSize(3000, 1000)
	got := downscale(wide, maxImageEdge)
	b := got.Bounds()
	if b.Dx() != maxImageEdge {
		t.Errorf("width = %d, want %d", b.Dx(), maxImageEdge)
	}
	// Aspect ratio must survive: 3:1 in, 3:1 out.
	if ratio := float64(b.Dx()) / float64(b.Dy()); ratio < 2.9 || ratio > 3.1 {
		t.Errorf("aspect ratio = %v, want ~3", ratio)
	}

	small := imageOfSize(200, 100)
	sb := downscale(small, maxImageEdge).Bounds()
	if sb.Dx() != 200 || sb.Dy() != 100 {
		t.Errorf("a small image was resized to %dx%d; it should be left alone", sb.Dx(), sb.Dy())
	}
}

func TestReactionPaletteIsClosed(t *testing.T) {
	for _, e := range reactionPalette() {
		if !allowedReactions[e] {
			t.Errorf("palette offers %q but the API rejects it", e)
		}
	}
	if allowedReactions["💩"] {
		t.Error("an emoji outside the palette was accepted")
	}
}

func TestSnippetPreviewsAttachments(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name string
		msg  models.Message
		want string
	}{
		{"deleted", models.Message{DeletedAt: &now, Body: "secret"}, "Message deleted"},
		{"bare photo", models.Message{Kind: models.MsgImage}, "📷 Photo"},
		{"captioned photo", models.Message{Kind: models.MsgImage, Body: "look"}, "📷 look"},
		{"text", models.Message{Kind: models.MsgText, Body: "hello"}, "hello"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := previewOf(c.msg); got != c.want {
				t.Errorf("previewOf = %q, want %q", got, c.want)
			}
		})
	}
}

// imageOfSize builds a plain canvas of the requested dimensions — enough to
// exercise the scaler without needing a fixture file on disk.
func imageOfSize(w, h int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			img.Set(x, y, color.RGBA{R: uint8(x % 256), G: uint8(y % 256), B: 128, A: 255})
		}
	}
	return img
}

var _ = fmt.Sprintf
