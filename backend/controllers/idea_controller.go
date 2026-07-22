package controllers

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"lumora/backend/database"
	"lumora/backend/middleware"
	"lumora/backend/models"
)

// IdeaController powers the ideas board: the left and right panels of the
// workspace (the centre panel is idea_message_controller.go).
type IdeaController struct{}

const (
	// An idea that clears this much net support is pushed into review
	// automatically. Nobody has to remember to escalate it.
	reviewThreshold = 20

	// Hot ranking half-life. Support decays so a two-year-old idea with 200
	// votes doesn't permanently outrank this morning's good one.
	hotHalfLife = 36 * time.Hour

	maxOpenIdeas = 60 // soft cap; past this the board nags about archiving
)

// statusFlow is the workflow, in order. Movement is validated against it so
// status stays meaningful rather than becoming a free-text field.
var statusFlow = []string{
	models.IdeaDraft,
	models.IdeaUnderReview,
	models.IdeaApproved,
	models.IdeaInProgress,
	models.IdeaCompleted,
	models.IdeaArchived,
}

func validStatus(s string) bool {
	for _, v := range statusFlow {
		if v == s {
			return true
		}
	}
	return false
}

// --- payloads ----------------------------------------------------------------

type ideaDTO struct {
	ID           uint     `json:"id"`
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	Status       string   `json:"status"`
	Owner        chatUser `json:"owner"`
	Upvotes      int      `json:"upvotes"`
	Downvotes    int      `json:"downvotes"`
	Score        int      `json:"score"`
	MyVote       int      `json:"myVote"`
	Starred      bool     `json:"starred"`
	Tags         []string `json:"tags"`
	MessageCount int      `json:"messageCount"`
	CreatedAt    string   `json:"createdAt"`
	LastActivity string   `json:"lastActivity"`
	Archived     bool     `json:"archived"`
	ArchiveNote  string   `json:"archiveReason"`
	MergedInto   *uint    `json:"mergedIntoId"`
	Heat         float64  `json:"heat"`
}

func (ic *IdeaController) toDTO(idea models.Idea, viewerID uint) ideaDTO {
	var owner models.User
	database.DB.First(&owner, idea.OwnerID)

	var tags []models.IdeaTag
	database.DB.Where("idea_id = ?", idea.ID).Order("tag asc").Find(&tags)
	tagNames := make([]string, 0, len(tags))
	for _, t := range tags {
		tagNames = append(tagNames, t.Tag)
	}

	myVote := 0
	var v models.IdeaVote
	if database.DB.Where("idea_id = ? AND user_id = ?", idea.ID, viewerID).
		First(&v).Error == nil {
		myVote = v.Value
	}
	var star int64
	database.DB.Model(&models.IdeaStar{}).
		Where("idea_id = ? AND user_id = ?", idea.ID, viewerID).Count(&star)

	return ideaDTO{
		ID: idea.ID, Title: idea.Title, Description: idea.Description,
		Status: idea.Status, Owner: toChatUser(owner),
		Upvotes: idea.Upvotes, Downvotes: idea.Downvotes, Score: idea.Score,
		MyVote: myVote, Starred: star > 0, Tags: tagNames,
		MessageCount: idea.MessageCount,
		CreatedAt:    idea.CreatedAt.Format(time.RFC3339),
		LastActivity: idea.LastActivityAt.Format(time.RFC3339),
		Archived:     idea.ArchivedAt != nil, ArchiveNote: idea.ArchiveReason,
		MergedInto: idea.MergedIntoID,
		Heat:       hotScore(idea),
	}
}

// --- the board ---------------------------------------------------------------

// List returns the board, filtered and sorted. This is the left panel.
func (ic *IdeaController) List(c *fiber.Ctx) error {
	user := middleware.CurrentUser(c)

	q := database.DB.Model(&models.Idea{})

	// Archived ideas are hidden unless explicitly asked for — that's the point
	// of archiving. Merged ones always hide from the board.
	status := c.Query("status")
	switch {
	case status != "" && validStatus(status):
		q = q.Where("status = ?", status)
	case status == "starred", status == "mine":
		// handled below, after the base filter
		q = q.Where("archived_at IS NULL")
	default:
		q = q.Where("archived_at IS NULL")
	}
	q = q.Where("merged_into_id IS NULL")

	if owner := c.Query("owner"); owner != "" {
		if id, err := strconv.Atoi(owner); err == nil {
			q = q.Where("owner_id = ?", id)
		}
	}
	if status == "mine" {
		q = q.Where("owner_id = ?", user.ID)
	}
	if tag := strings.TrimSpace(c.Query("tag")); tag != "" {
		var ids []uint
		database.DB.Model(&models.IdeaTag{}).Where("tag = ?", normaliseTag(tag)).
			Pluck("idea_id", &ids)
		if len(ids) == 0 {
			return c.JSON(fiber.Map{"ideas": []ideaDTO{}, "tags": ic.allTags(), "counts": ic.counts(user.ID)})
		}
		q = q.Where("id IN ?", ids)
	}
	if search := strings.TrimSpace(c.Query("q")); search != "" {
		like := "%" + strings.ToLower(search) + "%"
		q = q.Where("LOWER(title) LIKE ? OR LOWER(description) LIKE ?", like, like)
	}

	var ideas []models.Idea
	q.Limit(300).Find(&ideas)

	if status == "starred" {
		var starred []uint
		database.DB.Model(&models.IdeaStar{}).Where("user_id = ?", user.ID).
			Pluck("idea_id", &starred)
		keep := map[uint]bool{}
		for _, id := range starred {
			keep[id] = true
		}
		filtered := ideas[:0]
		for _, i := range ideas {
			if keep[i.ID] {
				filtered = append(filtered, i)
			}
		}
		ideas = filtered
	}

	sortIdeas(ideas, c.Query("sort"))

	out := make([]ideaDTO, 0, len(ideas))
	for _, i := range ideas {
		out = append(out, ic.toDTO(i, user.ID))
	}

	var open int64
	database.DB.Model(&models.Idea{}).
		Where("archived_at IS NULL AND merged_into_id IS NULL AND status NOT IN ?",
			[]string{models.IdeaCompleted, models.IdeaArchived}).Count(&open)

	return c.JSON(fiber.Map{
		"ideas":  out,
		"tags":   ic.allTags(),
		"counts": ic.counts(user.ID),
		// Surfaced so the UI can nudge toward archiving before the board becomes
		// a wall nobody reads.
		"openIdeas":    open,
		"maxOpenIdeas": maxOpenIdeas,
		"crowded":      open > maxOpenIdeas,
	})
}

// hotScore ranks by support, decayed by age — the "trending" sort. An idea
// posted an hour ago with 5 votes should be able to out-rank one from last
// month with 15.
func hotScore(i models.Idea) float64 {
	base := float64(i.Score)
	sign := 1.0
	if base < 0 {
		sign, base = -1.0, -base
	}
	age := time.Since(i.LastActivityAt)
	if i.LastActivityAt.IsZero() {
		age = time.Since(i.CreatedAt)
	}
	decay := math.Pow(0.5, age.Hours()/hotHalfLife.Hours())
	// Discussion is a support signal in its own right: an idea people are still
	// arguing about is live even before the votes land.
	engagement := math.Log1p(float64(i.MessageCount)) * 2
	return (sign*math.Log1p(base)*10 + engagement) * decay
}

// controversy peaks when support and opposition are both high and evenly
// matched. A 30/28 split scores far above a 2/0 one — which is the whole point:
// the best ideas often start divisive, and pure score buries them.
func controversy(i models.Idea) float64 {
	up, down := float64(i.Upvotes), float64(i.Downvotes)
	total := up + down
	if total < 2 || up == 0 || down == 0 {
		return 0
	}
	balance := math.Min(up, down) / math.Max(up, down)
	return total * math.Pow(balance, 1.5)
}

func sortIdeas(ideas []models.Idea, mode string) {
	switch mode {
	case "top":
		sort.SliceStable(ideas, func(a, b int) bool { return ideas[a].Score > ideas[b].Score })
	case "new":
		sort.SliceStable(ideas, func(a, b int) bool {
			return ideas[a].CreatedAt.After(ideas[b].CreatedAt)
		})
	case "controversial":
		sort.SliceStable(ideas, func(a, b int) bool {
			return controversy(ideas[a]) > controversy(ideas[b])
		})
	default: // hot
		sort.SliceStable(ideas, func(a, b int) bool {
			return hotScore(ideas[a]) > hotScore(ideas[b])
		})
	}
}

func (ic *IdeaController) allTags() []fiber.Map {
	type row struct {
		Tag string
		N   int
	}
	var rows []row
	database.DB.Model(&models.IdeaTag{}).
		Select("tag, COUNT(*) as n").Group("tag").Order("n desc").Limit(40).Scan(&rows)
	out := make([]fiber.Map, 0, len(rows))
	for _, r := range rows {
		out = append(out, fiber.Map{"tag": r.Tag, "count": r.N})
	}
	return out
}

func (ic *IdeaController) counts(userID uint) fiber.Map {
	count := func(where string, args ...interface{}) int64 {
		var n int64
		database.DB.Model(&models.Idea{}).Where(where, args...).Count(&n)
		return n
	}
	var starred int64
	database.DB.Model(&models.IdeaStar{}).Where("user_id = ?", userID).Count(&starred)

	return fiber.Map{
		"all":         count("archived_at IS NULL AND merged_into_id IS NULL"),
		"draft":       count("status = ? AND archived_at IS NULL", models.IdeaDraft),
		"underReview": count("status = ? AND archived_at IS NULL", models.IdeaUnderReview),
		"approved":    count("status = ? AND archived_at IS NULL", models.IdeaApproved),
		"inProgress":  count("status = ? AND archived_at IS NULL", models.IdeaInProgress),
		"completed":   count("status = ?", models.IdeaCompleted),
		"archived":    count("archived_at IS NOT NULL"),
		"mine":        count("owner_id = ? AND archived_at IS NULL", userID),
		"starred":     starred,
	}
}

// --- create / read / update --------------------------------------------------

type ideaInput struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
	Status      string   `json:"status"`
}

// Create posts a new idea. Only the title is required — every extra mandatory
// field is a reason not to bother, and friction is what kills a board.
func (ic *IdeaController) Create(c *fiber.Ctx) error {
	user := middleware.CurrentUser(c)

	var in ideaInput
	if err := c.BodyParser(&in); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
	}
	title := strings.TrimSpace(in.Title)
	if title == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "an idea needs a title"})
	}
	if len(title) > 160 {
		title = title[:160]
	}

	now := time.Now()
	idea := models.Idea{
		OwnerID: user.ID, Title: title,
		Description:    strings.TrimSpace(in.Description),
		Status:         models.IdeaDraft,
		LastActivityAt: now,
	}
	if err := database.DB.Create(&idea).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not save idea"})
	}
	setTags(idea.ID, in.Tags)
	logIdeaEvent(idea.ID, user.ID, "created", "", "", title, "")

	// The author's own idea starts with their vote. Posting something is a
	// stronger endorsement than clicking an arrow.
	applyVote(&idea, user.ID, 1)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"idea": ic.toDTO(idea, user.ID)})
}

// Get returns one idea with everything the details panel shows.
func (ic *IdeaController) Get(c *fiber.Ctx) error {
	user := middleware.CurrentUser(c)

	idea, ok := ic.load(c)
	if !ok {
		return nil // the helper already wrote the error response
	}

	var history []models.IdeaEvent
	database.DB.Where("idea_id = ?", idea.ID).Order("created_at desc").Limit(50).Find(&history)
	events := make([]fiber.Map, 0, len(history))
	for _, e := range history {
		var actor models.User
		database.DB.First(&actor, e.ActorID)
		events = append(events, fiber.Map{
			"id": e.ID, "kind": e.Kind, "field": e.Field, "from": e.From, "to": e.To,
			"note": e.Note, "actor": toChatUser(actor),
			"at": e.CreatedAt.Format(time.RFC3339),
		})
	}

	// Initialised rather than declared: GORM leaves a slice nil when nothing
	// matches, which marshals to `null` and forces a guard on every client
	// access. An empty list is the honest representation of "no tasks".
	tasks := []models.IdeaTask{}
	database.DB.Where("idea_id = ?", idea.ID).Order("created_at asc").Find(&tasks)

	// Who's in the room — everyone who has posted in the thread.
	var authorIDs []uint
	database.DB.Model(&models.IdeaMessage{}).
		Where("idea_id = ? AND deleted_at IS NULL AND anonymous = ?", idea.ID, false).
		Distinct().Pluck("author_id", &authorIDs)
	participants := make([]chatUser, 0, len(authorIDs))
	for _, id := range authorIDs {
		var u models.User
		if database.DB.First(&u, id).Error == nil {
			participants = append(participants, toChatUser(u))
		}
	}

	var merged []models.Idea
	database.DB.Where("merged_into_id = ?", idea.ID).Find(&merged)
	mergedDTO := make([]ideaDTO, 0, len(merged))
	for _, m := range merged {
		mergedDTO = append(mergedDTO, ic.toDTO(m, user.ID))
	}

	return c.JSON(fiber.Map{
		"idea":         ic.toDTO(*idea, user.ID),
		"history":      events,
		"tasks":        tasks,
		"participants": participants,
		"mergedIn":     mergedDTO,
		"canEdit":      idea.OwnerID == user.ID,
		"statusFlow":   statusFlow,
		"similar":      ic.similarTo(*idea, 3),
	})
}

// Update edits an idea. Every change is recorded in the history, so the details
// panel can always answer "who changed this, and when".
func (ic *IdeaController) Update(c *fiber.Ctx) error {
	user := middleware.CurrentUser(c)

	idea, ok := ic.load(c)
	if !ok {
		return nil // the helper already wrote the error response
	}
	var in ideaInput
	if err := c.BodyParser(&in); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
	}

	if t := strings.TrimSpace(in.Title); t != "" && t != idea.Title {
		if !ic.mayEdit(idea, user) {
			return forbidden(c)
		}
		logIdeaEvent(idea.ID, user.ID, "edited", "title", idea.Title, t, "")
		idea.Title = t
	}
	if in.Description != idea.Description {
		if !ic.mayEdit(idea, user) {
			return forbidden(c)
		}
		logIdeaEvent(idea.ID, user.ID, "edited", "description",
			snippet(idea.Description), snippet(in.Description), "")
		idea.Description = strings.TrimSpace(in.Description)
	}
	if in.Status != "" && in.Status != idea.Status {
		if !validStatus(in.Status) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "unknown status"})
		}
		// Status is a team decision, not an ownership one — anyone can move it,
		// and the history records who did.
		logIdeaEvent(idea.ID, user.ID, "status", "status", idea.Status, in.Status, "")
		idea.Status = in.Status
		notifyThread(*idea, user.ID, "📌", "Status: "+statusLabel(in.Status),
			fmt.Sprintf("%s moved \"%s\" to %s.", displayName(*user), idea.Title, statusLabel(in.Status)))
	}
	if in.Tags != nil {
		if !ic.mayEdit(idea, user) {
			return forbidden(c)
		}
		before := strings.Join(tagsOf(idea.ID), ", ")
		setTags(idea.ID, in.Tags)
		after := strings.Join(tagsOf(idea.ID), ", ")
		if before != after {
			logIdeaEvent(idea.ID, user.ID, "tagged", "tags", before, after, "")
		}
	}

	idea.LastActivityAt = time.Now()
	database.DB.Save(idea)
	return c.JSON(fiber.Map{"idea": ic.toDTO(*idea, user.ID)})
}

// mayEdit: the owner edits content. Everyone can move status and vote — an
// idea board where only the author can act is just a suggestion box.
func (ic *IdeaController) mayEdit(idea *models.Idea, user *models.User) bool {
	return idea.OwnerID == user.ID
}

// Delete removes an idea and everything attached to it — permanently.
//
// The owner may delete at any point, including an idea people have discussed.
// Archiving is still the better move once others have contributed (it keeps the
// thread readable and records why), so the client warns about exactly what will
// be lost — but the decision belongs to whoever posted it, not to the API.
//
// The cascade is manual because these tables have no foreign keys: without it,
// deleting an idea would strand its messages, votes and attachment bytes in the
// database forever.
func (ic *IdeaController) Delete(c *fiber.Ctx) error {
	user := middleware.CurrentUser(c)

	idea, ok := ic.load(c)
	if !ok {
		return nil // the helper already wrote the error response
	}
	if idea.OwnerID != user.ID {
		return c.Status(fiber.StatusForbidden).
			JSON(fiber.Map{"error": "only the person who posted an idea can delete it"})
	}

	// Message ids first — reactions hang off them, not off the idea.
	var messageIDs []uint
	database.DB.Model(&models.IdeaMessage{}).Where("idea_id = ?", idea.ID).
		Pluck("id", &messageIDs)

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		if len(messageIDs) > 0 {
			if err := tx.Where("message_id IN ?", messageIDs).
				Delete(&models.IdeaReaction{}).Error; err != nil {
				return err
			}
		}
		for _, model := range []interface{}{
			&models.IdeaMessage{}, &models.IdeaTag{}, &models.IdeaVote{},
			&models.IdeaStar{}, &models.IdeaEvent{}, &models.IdeaTask{},
			&models.BrainstormSession{},
		} {
			if err := tx.Where("idea_id = ?", idea.ID).Delete(model).Error; err != nil {
				return err
			}
		}
		// Anything merged into this idea would otherwise point at a row that no
		// longer exists; release it back onto the board instead.
		if err := tx.Model(&models.Idea{}).Where("merged_into_id = ?", idea.ID).
			Update("merged_into_id", nil).Error; err != nil {
			return err
		}
		return tx.Delete(idea).Error
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"error": "could not delete that idea"})
	}

	// The notifications pointed at a thread that no longer exists.
	database.DB.Where("kind = ? AND link LIKE ?", "idea",
		fmt.Sprintf("%%idea=%d%%", idea.ID)).Delete(&models.Notification{})

	return c.JSON(fiber.Map{"ok": true, "deletedMessages": len(messageIDs)})
}

// --- voting ------------------------------------------------------------------

type voteInput struct {
	Value int `json:"value"` // 1, -1, or 0 to clear
}

// Vote records support or opposition in one tap. No comment required — asking
// people to justify a vote means most of them simply don't vote.
func (ic *IdeaController) Vote(c *fiber.Ctx) error {
	user := middleware.CurrentUser(c)

	idea, ok := ic.load(c)
	if !ok {
		return nil // the helper already wrote the error response
	}
	var in voteInput
	if err := c.BodyParser(&in); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
	}
	if in.Value != 1 && in.Value != -1 && in.Value != 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "vote must be 1, -1 or 0"})
	}

	applyVote(idea, user.ID, in.Value)
	return c.JSON(fiber.Map{"idea": ic.toDTO(*idea, user.ID)})
}

// applyVote writes the vote and recalculates the cached tallies from scratch.
// Recomputing rather than incrementing means a double-submit or a race can't
// leave the counters drifting away from the votes that actually exist.
func applyVote(idea *models.Idea, userID uint, value int) {
	database.DB.Where("idea_id = ? AND user_id = ?", idea.ID, userID).
		Delete(&models.IdeaVote{})
	if value != 0 {
		database.DB.Create(&models.IdeaVote{IdeaID: idea.ID, UserID: userID, Value: value})
	}

	var up, down int64
	database.DB.Model(&models.IdeaVote{}).Where("idea_id = ? AND value > 0", idea.ID).Count(&up)
	database.DB.Model(&models.IdeaVote{}).Where("idea_id = ? AND value < 0", idea.ID).Count(&down)

	wasBelow := idea.Score < reviewThreshold
	idea.Upvotes, idea.Downvotes = int(up), int(down)
	idea.Score = int(up - down)

	// Crossing the threshold escalates the idea on its own. Nobody has to
	// notice, which is the point — good ideas shouldn't need a champion with a
	// calendar reminder.
	if wasBelow && idea.Score >= reviewThreshold && idea.Status == models.IdeaDraft {
		idea.Status = models.IdeaUnderReview
		logIdeaEvent(idea.ID, userID, "vote_threshold", "status",
			models.IdeaDraft, models.IdeaUnderReview,
			fmt.Sprintf("reached %d votes", reviewThreshold))
		notifyThread(*idea, 0, "🚀", "Idea flagged for review",
			fmt.Sprintf("\"%s\" passed %d votes and moved to Under review.", idea.Title, reviewThreshold))
	}
	database.DB.Save(idea)
}

// Star toggles a private bookmark.
func (ic *IdeaController) Star(c *fiber.Ctx) error {
	user := middleware.CurrentUser(c)

	idea, ok := ic.load(c)
	if !ok {
		return nil // the helper already wrote the error response
	}
	var existing models.IdeaStar
	if database.DB.Where("idea_id = ? AND user_id = ?", idea.ID, user.ID).
		First(&existing).Error == nil {
		database.DB.Delete(&existing)
		return c.JSON(fiber.Map{"starred": false})
	}
	database.DB.Create(&models.IdeaStar{IdeaID: idea.ID, UserID: user.ID})
	return c.JSON(fiber.Map{"starred": true})
}

// --- archive / merge / tasks -------------------------------------------------

type archiveInput struct {
	Reason string `json:"reason"`
}

// Archive closes an idea, always with a reason. A board that only grows becomes
// a wall nobody reads; one that closes things silently teaches people their
// contributions vanish.
func (ic *IdeaController) Archive(c *fiber.Ctx) error {
	user := middleware.CurrentUser(c)

	idea, ok := ic.load(c)
	if !ok {
		return nil // the helper already wrote the error response
	}
	var in archiveInput
	_ = c.BodyParser(&in)
	reason := strings.TrimSpace(in.Reason)
	if reason == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "archiving needs a reason — it's what tells everyone else what happened",
		})
	}

	now := time.Now()
	idea.ArchivedAt = &now
	idea.ArchiveReason = reason
	idea.Status = models.IdeaArchived
	idea.LastActivityAt = now
	database.DB.Save(idea)

	logIdeaEvent(idea.ID, user.ID, "archived", "status", "", models.IdeaArchived, reason)
	notifyThread(*idea, user.ID, "🗄️", "Idea archived",
		fmt.Sprintf("\"%s\" was archived: %s", idea.Title, reason))

	return c.JSON(fiber.Map{"idea": ic.toDTO(*idea, user.ID)})
}

// Restore pulls an archived idea back onto the board.
func (ic *IdeaController) Restore(c *fiber.Ctx) error {
	user := middleware.CurrentUser(c)

	idea, ok := ic.load(c)
	if !ok {
		return nil // the helper already wrote the error response
	}
	idea.ArchivedAt = nil
	idea.ArchiveReason = ""
	idea.Status = models.IdeaUnderReview
	idea.LastActivityAt = time.Now()
	database.DB.Save(idea)

	logIdeaEvent(idea.ID, user.ID, "restored", "status", models.IdeaArchived, idea.Status, "")
	return c.JSON(fiber.Map{"idea": ic.toDTO(*idea, user.ID)})
}

type mergeInput struct {
	TargetID uint `json:"targetId"`
}

// Merge folds this idea into another: votes transfer (deduplicated per person),
// tags combine, and the thread stays readable under the surviving idea.
func (ic *IdeaController) Merge(c *fiber.Ctx) error {
	user := middleware.CurrentUser(c)

	idea, ok := ic.load(c)
	if !ok {
		return nil // the helper already wrote the error response
	}
	var in mergeInput
	if err := c.BodyParser(&in); err != nil || in.TargetID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "pick an idea to merge into"})
	}
	if in.TargetID == idea.ID {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "an idea can't merge into itself"})
	}

	var target models.Idea
	if database.DB.First(&target, in.TargetID).Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "target idea not found"})
	}
	if target.MergedIntoID != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "that idea has itself been merged — merge into the surviving one",
		})
	}

	// Move votes across, but only from people who haven't already voted on the
	// target: merging must not let one person count twice.
	var votes []models.IdeaVote
	database.DB.Where("idea_id = ?", idea.ID).Find(&votes)
	for _, v := range votes {
		var dup int64
		database.DB.Model(&models.IdeaVote{}).
			Where("idea_id = ? AND user_id = ?", target.ID, v.UserID).Count(&dup)
		if dup == 0 {
			database.DB.Create(&models.IdeaVote{
				IdeaID: target.ID, UserID: v.UserID, Value: v.Value,
			})
		}
	}
	// Carry the tags over.
	setTags(target.ID, append(tagsOf(target.ID), tagsOf(idea.ID)...))
	// And the conversation.
	database.DB.Model(&models.IdeaMessage{}).Where("idea_id = ?", idea.ID).
		Update("idea_id", target.ID)
	var msgCount int64
	database.DB.Model(&models.IdeaMessage{}).
		Where("idea_id = ? AND deleted_at IS NULL", target.ID).Count(&msgCount)
	target.MessageCount = int(msgCount)

	now := time.Now()
	idea.MergedIntoID = &target.ID
	idea.ArchivedAt = &now
	idea.ArchiveReason = "merged into #" + strconv.Itoa(int(target.ID))
	idea.Status = models.IdeaArchived
	database.DB.Save(idea)

	applyVote(&target, user.ID, voteValueOf(target.ID, user.ID)) // recount
	target.LastActivityAt = now
	database.DB.Save(&target)

	logIdeaEvent(idea.ID, user.ID, "merged", "", idea.Title, target.Title, "")
	logIdeaEvent(target.ID, user.ID, "merged", "", idea.Title, target.Title,
		"absorbed #"+strconv.Itoa(int(idea.ID)))

	return c.JSON(fiber.Map{"idea": ic.toDTO(*idea, user.ID), "target": ic.toDTO(target, user.ID)})
}

// voteValueOf reads back a user's current vote, used to force a recount without
// changing their position.
func voteValueOf(ideaID, userID uint) int {
	var v models.IdeaVote
	if database.DB.Where("idea_id = ? AND user_id = ?", ideaID, userID).First(&v).Error != nil {
		return 0
	}
	return v.Value
}

type taskInput struct {
	Title  string `json:"title"`
	Sprint string `json:"sprint"`
	Status string `json:"status"`
}

// CreateTask graduates an idea into work and moves it to In progress.
func (ic *IdeaController) CreateTask(c *fiber.Ctx) error {
	user := middleware.CurrentUser(c)

	idea, ok := ic.load(c)
	if !ok {
		return nil // the helper already wrote the error response
	}
	var in taskInput
	_ = c.BodyParser(&in)
	title := strings.TrimSpace(in.Title)
	if title == "" {
		title = idea.Title
	}

	task := models.IdeaTask{
		IdeaID: idea.ID, Title: title, Status: "todo",
		Sprint: strings.TrimSpace(in.Sprint), CreatedByID: user.ID,
	}
	database.DB.Create(&task)

	if idea.Status != models.IdeaInProgress && idea.Status != models.IdeaCompleted {
		logIdeaEvent(idea.ID, user.ID, "status", "status", idea.Status, models.IdeaInProgress, "converted to a task")
		idea.Status = models.IdeaInProgress
	}
	idea.LastActivityAt = time.Now()
	database.DB.Save(idea)

	logIdeaEvent(idea.ID, user.ID, "task", "", "", title, in.Sprint)
	notifyThread(*idea, user.ID, "✅", "Idea is now a task",
		fmt.Sprintf("\"%s\" was converted to a task by %s. Status: In progress.", idea.Title, displayName(*user)))

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"task": task, "idea": ic.toDTO(*idea, user.ID)})
}

// UpdateTask moves a task along, completing the idea when the last one is done.
func (ic *IdeaController) UpdateTask(c *fiber.Ctx) error {
	user := middleware.CurrentUser(c)

	id, err := strconv.Atoi(c.Params("taskId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid task"})
	}
	var task models.IdeaTask
	if database.DB.First(&task, id).Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "task not found"})
	}

	var in taskInput
	_ = c.BodyParser(&in)
	if in.Status != "" {
		task.Status = in.Status
		if in.Status == "done" {
			now := time.Now()
			task.CompletedAt = &now
		} else {
			task.CompletedAt = nil
		}
	}
	if in.Sprint != "" {
		task.Sprint = strings.TrimSpace(in.Sprint)
	}
	database.DB.Save(&task)

	var open int64
	database.DB.Model(&models.IdeaTask{}).
		Where("idea_id = ? AND status != ?", task.IdeaID, "done").Count(&open)
	var idea models.Idea
	if database.DB.First(&idea, task.IdeaID).Error == nil && open == 0 &&
		idea.Status == models.IdeaInProgress {
		logIdeaEvent(idea.ID, user.ID, "status", "status", idea.Status, models.IdeaCompleted,
			"all tasks completed")
		idea.Status = models.IdeaCompleted
		idea.LastActivityAt = time.Now()
		database.DB.Save(&idea)
	}
	return c.JSON(fiber.Map{"task": task})
}

// --- helpers -----------------------------------------------------------------

// load fetches the idea named in the path. It returns ok=false having ALREADY
// written the error response, so the caller must `return nil` rather than
// writing a second one.
//
// The bool is not decoration. Returning `error` here is a trap: Fiber's
// c.JSON() returns nil on success, so `return nil, c.Status(404).JSON(...)`
// hands the caller a nil error alongside a nil idea — the guard passes and the
// next line dereferences nil. A 404 becomes a 500 panic.
func (ic *IdeaController) load(c *fiber.Ctx) (*models.Idea, bool) {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		_ = c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid idea"})
		return nil, false
	}
	var idea models.Idea
	if database.DB.First(&idea, id).Error != nil {
		_ = c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "idea not found"})
		return nil, false
	}
	return &idea, true
}

func forbidden(c *fiber.Ctx) error {
	return c.Status(fiber.StatusForbidden).
		JSON(fiber.Map{"error": "only the owner can change this"})
}

func statusLabel(s string) string {
	switch s {
	case models.IdeaDraft:
		return "Draft"
	case models.IdeaUnderReview:
		return "Under review"
	case models.IdeaApproved:
		return "Approved"
	case models.IdeaInProgress:
		return "In progress"
	case models.IdeaCompleted:
		return "Completed"
	case models.IdeaArchived:
		return "Archived"
	}
	return s
}

// normaliseTag keeps "#AI", "ai" and " AI " from becoming three separate tags.
func normaliseTag(t string) string {
	t = strings.ToLower(strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(t), "#")))
	t = strings.ReplaceAll(t, " ", "-")
	if len(t) > 32 {
		t = t[:32]
	}
	return t
}

func setTags(ideaID uint, tags []string) {
	database.DB.Where("idea_id = ?", ideaID).Delete(&models.IdeaTag{})
	seen := map[string]bool{}
	for _, raw := range tags {
		t := normaliseTag(raw)
		if t == "" || seen[t] {
			continue
		}
		seen[t] = true
		database.DB.Create(&models.IdeaTag{IdeaID: ideaID, Tag: t})
	}
}

func tagsOf(ideaID uint) []string {
	var tags []models.IdeaTag
	database.DB.Where("idea_id = ?", ideaID).Order("tag asc").Find(&tags)
	out := make([]string, 0, len(tags))
	for _, t := range tags {
		out = append(out, t.Tag)
	}
	return out
}

func logIdeaEvent(ideaID, actorID uint, kind, field, from, to, note string) {
	database.DB.Create(&models.IdeaEvent{
		IdeaID: ideaID, ActorID: actorID, Kind: kind, Field: field,
		From: from, To: to, Note: note,
	})
}

// notifyThread tells everyone who has taken part in an idea that something
// happened to it — the "everyone in the thread gets notified" step. The actor
// is skipped: nobody needs telling about their own action.
func notifyThread(idea models.Idea, actorID uint, emoji, title, body string) {
	recipients := map[uint]bool{idea.OwnerID: true}

	var authorIDs []uint
	database.DB.Model(&models.IdeaMessage{}).
		Where("idea_id = ? AND deleted_at IS NULL", idea.ID).
		Distinct().Pluck("author_id", &authorIDs)
	for _, id := range authorIDs {
		recipients[id] = true
	}
	var voterIDs []uint
	database.DB.Model(&models.IdeaVote{}).Where("idea_id = ?", idea.ID).
		Distinct().Pluck("user_id", &voterIDs)
	for _, id := range voterIDs {
		recipients[id] = true
	}
	delete(recipients, actorID)

	link := "/ideas?idea=" + strconv.Itoa(int(idea.ID))
	for id := range recipients {
		database.DB.Create(&models.Notification{
			UserID: id,
			// Keyed per idea per kind so a burst of status changes collapses
			// rather than filling someone's feed.
			Key:  fmt.Sprintf("idea_%d_%s", idea.ID, strings.ToLower(title)),
			Kind: "idea", Emoji: emoji, Tint: "#6C3FC5",
			Title: title, Body: body, Link: link,
		})
	}
}
