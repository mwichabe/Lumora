package controllers

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/image/draw"

	"lumora/backend/database"
	"lumora/backend/middleware"
	"lumora/backend/models"
)

// The centre panel: one thread per idea, branchable, with rich messages.

const (
	maxAttachment  = 8 << 20 // 8 MB upload ceiling
	maxImageEdge   = 1600    // longest side after downscaling
	maxMessageBody = 4000
	// A message can be edited freely, but the edit is always marked. Silent
	// edits let someone rewrite what they agreed to after the fact.
	editWindow = 24 * time.Hour
)

// The reaction set. A fixed palette rather than a full emoji picker: four
// meanings that cover agreement, insight, enthusiasm and confusion keep chat
// lightweight, which is the point of reacting instead of replying.
var allowedReactions = map[string]bool{
	"👍": true, "💡": true, "🔥": true, "❓": true, "🎉": true, "👀": true,
}

var (
	mentionUser = regexp.MustCompile(`@([A-Za-z][A-Za-z0-9_.-]{1,30})`)
	mentionIdea = regexp.MustCompile(`@idea#(\d+)`)
)

// --- payloads ----------------------------------------------------------------

type reactionDTO struct {
	Emoji string `json:"emoji"`
	Count int    `json:"count"`
	Mine  bool   `json:"mine"`
}

type messageDTO struct {
	ID        uint      `json:"id"`
	IdeaID    uint      `json:"ideaId"`
	ParentID  *uint     `json:"parentId"`
	Author    *chatUser `json:"author"` // nil for anonymous brainstorm posts
	Kind      string    `json:"kind"`
	Body      string    `json:"body"`
	FileName  string    `json:"fileName"`
	URL       string    `json:"url"` // attachment endpoint, when there is one
	Width     int       `json:"width"`
	Height    int       `json:"height"`
	Duration  int       `json:"duration"`
	Anonymous bool      `json:"anonymous"`
	// Present only when the message isn't English — see translation.go.
	Translation *translationDTO `json:"translation,omitempty"`
	Reactions   []reactionDTO   `json:"reactions"`
	Replies     []messageDTO    `json:"replies"`
	ReplyCount  int             `json:"replyCount"`
	Mine        bool            `json:"mine"`
	Edited      bool            `json:"edited"`
	Deleted     bool            `json:"deleted"`
	CanEdit     bool            `json:"canEdit"`
	CreatedAt   string          `json:"createdAt"`
}

func (ic *IdeaController) toMessageDTO(m models.IdeaMessage, viewerID uint) messageDTO {
	dto := messageDTO{
		ID: m.ID, IdeaID: m.IdeaID, ParentID: m.ParentID,
		Kind: m.Kind, Body: m.Body, FileName: m.FileName,
		Width: m.Width, Height: m.Height, Duration: m.Duration,
		Anonymous: m.Anonymous, Mine: m.AuthorID == viewerID,
		Edited: m.EditedAt != nil, Deleted: m.DeletedAt != nil,
		CreatedAt: m.CreatedAt.Format(time.RFC3339),
		// Empty slices, not nil: a nil slice marshals to `null`, and the client
		// reads these as arrays without a guard on every access.
		Reactions: []reactionDTO{},
		Replies:   []messageDTO{},
	}

	// A deleted message leaves a tombstone rather than vanishing: removing it
	// outright would silently rewrite a conversation other people replied to.
	if m.DeletedAt != nil {
		dto.Body = ""
		dto.Kind = models.MsgText
		return dto
	}

	// Anonymity is enforced here, at the boundary. The author is still stored
	// (moderation needs it) but never leaves the server.
	if !m.Anonymous {
		var u models.User
		if database.DB.First(&u, m.AuthorID).Error == nil {
			cu := toChatUser(u)
			dto.Author = &cu
		}
	}
	if len(m.Data) > 0 {
		dto.URL = fmt.Sprintf("/api/ideas/attachments/%d", m.ID)
	}
	dto.CanEdit = m.AuthorID == viewerID && time.Since(m.CreatedAt) < editWindow &&
		(m.Kind == models.MsgText || m.Kind == models.MsgCode)
	dto.Translation = translationFor(m.Body, m.DetectedLang, m.TranslatedBody, m.TranslatedAt)

	var reactions []models.IdeaReaction
	database.DB.Where("message_id = ?", m.ID).Find(&reactions)
	byEmoji := map[string]*reactionDTO{}
	order := []string{}
	for _, r := range reactions {
		e := byEmoji[r.Emoji]
		if e == nil {
			e = &reactionDTO{Emoji: r.Emoji}
			byEmoji[r.Emoji] = e
			order = append(order, r.Emoji)
		}
		e.Count++
		if r.UserID == viewerID {
			e.Mine = true
		}
	}
	for _, e := range order {
		dto.Reactions = append(dto.Reactions, *byEmoji[e])
	}
	return dto
}

// --- reading -----------------------------------------------------------------

// Messages returns the idea's thread as a two-level tree: top-level messages,
// each with its branch of replies. Two levels, not n, on purpose — deeper
// nesting is where threaded discussions become unreadable.
func (ic *IdeaController) Messages(c *fiber.Ctx) error {
	user := middleware.CurrentUser(c)

	idea, ok := ic.load(c)
	if !ok {
		return nil // the helper already wrote the error response
	}

	var all []models.IdeaMessage
	database.DB.Where("idea_id = ?", idea.ID).Order("created_at asc").Find(&all)

	repliesByParent := map[uint][]messageDTO{}
	var roots []models.IdeaMessage
	for _, m := range all {
		if m.ParentID == nil {
			roots = append(roots, m)
			continue
		}
		repliesByParent[*m.ParentID] = append(repliesByParent[*m.ParentID],
			ic.toMessageDTO(m, user.ID))
	}

	out := make([]messageDTO, 0, len(roots))
	for _, m := range roots {
		dto := ic.toMessageDTO(m, user.ID)
		if replies := repliesByParent[m.ID]; replies != nil {
			dto.Replies = replies
		}
		dto.ReplyCount = len(dto.Replies)
		out = append(out, dto)
	}

	return c.JSON(fiber.Map{
		"messages":   out,
		"idea":       ic.toDTO(*idea, user.ID),
		"brainstorm": activeBrainstorm(idea.ID),
		"reactions":  reactionPalette(),
	})
}

func reactionPalette() []string {
	return []string{"👍", "💡", "🔥", "❓", "🎉", "👀"}
}

// Attachment serves an image or voice memo. Unauthenticated for the same reason
// avatars are: the browser renders these in a plain <img>/<audio>, which can't
// attach a bearer token. Only the bytes are exposed, keyed by an opaque id.
func (ic *IdeaController) Attachment(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("bad id")
	}
	var m models.IdeaMessage
	if database.DB.First(&m, id).Error != nil || len(m.Data) == 0 || m.DeletedAt != nil {
		return c.Status(fiber.StatusNotFound).SendString("not found")
	}
	c.Set("Content-Type", m.Mime)
	c.Set("Cache-Control", "public, max-age=31536000, immutable")
	return c.Send(m.Data)
}

// --- writing -----------------------------------------------------------------

type postInput struct {
	Body     string `json:"body"`
	ParentID *uint  `json:"parentId"`
	Kind     string `json:"kind"`
}

// Post adds a message to an idea's thread. Accepts either JSON (text and code)
// or multipart (image and voice), so the composer has one endpoint.
func (ic *IdeaController) Post(c *fiber.Ctx) error {
	user := middleware.CurrentUser(c)

	idea, ok := ic.load(c)
	if !ok {
		return nil // the helper already wrote the error response
	}

	msg := models.IdeaMessage{IdeaID: idea.ID, AuthorID: user.ID, Kind: models.MsgText}

	if strings.HasPrefix(c.Get("Content-Type"), "multipart/form-data") {
		if status, problem := ic.attach(c, &msg); problem != "" {
			return c.Status(status).JSON(fiber.Map{"error": problem})
		}
		msg.Body = strings.TrimSpace(c.FormValue("body"))
		if pid := c.FormValue("parentId"); pid != "" {
			if n, err := strconv.Atoi(pid); err == nil && n > 0 {
				id := uint(n)
				msg.ParentID = &id
			}
		}
	} else {
		var in postInput
		if err := c.BodyParser(&in); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
		}
		msg.Body = strings.TrimSpace(in.Body)
		msg.ParentID = in.ParentID
		if in.Kind == models.MsgCode {
			msg.Kind = models.MsgCode
		}
		if msg.Body == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "message is empty"})
		}
	}

	if len(msg.Body) > maxMessageBody {
		msg.Body = msg.Body[:maxMessageBody]
	}
	// A reply may only hang off a top-level message in this thread — that keeps
	// the tree two deep and stops a reply being re-parented into another idea.
	if msg.ParentID != nil {
		var parent models.IdeaMessage
		if database.DB.First(&parent, *msg.ParentID).Error != nil ||
			parent.IdeaID != idea.ID || parent.ParentID != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid reply target"})
		}
	}

	// During a silent brainstorm every contribution is anonymous, so the
	// discussion is judged on content rather than on who said it.
	if b := activeBrainstorm(idea.ID); b != nil {
		msg.Anonymous = true
	}

	if err := database.DB.Create(&msg).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not post message"})
	}

	// Detect the language now and translate in the background if needed. Code
	// blocks are skipped — a snippet is not prose, and "detecting" a language
	// in it would produce a nonsense label.
	if msg.Kind != models.MsgCode {
		msg.DetectedLang = detectAndQueue("idea_messages", msg.ID, msg.Body)
	}

	touchIdea(idea)
	if !msg.Anonymous {
		notifyMentions(*idea, msg, *user)
	}

	return c.Status(fiber.StatusCreated).
		JSON(fiber.Map{"message": ic.toMessageDTO(msg, user.ID)})
}

// attach reads the uploaded file, validates it by decoding, and downscales
// images. Voice memos are stored as-is — they're already compressed by the
// browser's recorder.
// It returns (status, problem); an empty problem means the attachment was
// accepted. Deliberately NOT (error): writing the error response here and
// returning c.JSON's result hands back a nil error, so the caller carries on
// and posts an empty message on top of the 400 it already sent.
func (ic *IdeaController) attach(c *fiber.Ctx, msg *models.IdeaMessage) (int, string) {
	fh, err := c.FormFile("file")
	if err != nil {
		return fiber.StatusBadRequest, "no file uploaded"
	}
	if fh.Size > maxAttachment {
		return fiber.StatusBadRequest, "attachment must be under 8MB"
	}
	f, err := fh.Open()
	if err != nil {
		return fiber.StatusInternalServerError, "could not read that file"
	}
	defer f.Close()

	kind := c.FormValue("kind")
	msg.FileName = fh.Filename

	if kind == models.MsgVoice {
		buf := new(bytes.Buffer)
		if _, err := buf.ReadFrom(f); err != nil {
			return fiber.StatusInternalServerError, "could not read that recording"
		}
		msg.Kind = models.MsgVoice
		msg.Data = buf.Bytes()
		msg.Mime = fh.Header.Get("Content-Type")
		if msg.Mime == "" {
			msg.Mime = "audio/webm"
		}
		if d, err := strconv.Atoi(c.FormValue("duration")); err == nil {
			msg.Duration = d
		}
		return 0, ""
	}

	// Decoding doubles as validation: a file that only claims to be an image by
	// its name fails here rather than being stored and served back.
	src, _, err := image.Decode(f)
	if err != nil {
		return fiber.StatusBadRequest, "could not read that image"
	}
	scaled := downscale(src, maxImageEdge)
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, scaled, &jpeg.Options{Quality: 85}); err != nil {
		return fiber.StatusInternalServerError, "could not process that image"
	}

	b := scaled.Bounds()
	msg.Kind = models.MsgImage
	msg.Data = buf.Bytes()
	msg.Mime = "image/jpeg"
	msg.Width, msg.Height = b.Dx(), b.Dy()
	return 0, ""
}

// downscale shrinks an image so its longest edge is at most max, preserving
// aspect ratio. Images smaller than that are re-encoded but not upscaled.
func downscale(src image.Image, max int) image.Image {
	b := src.Bounds()
	w, h := b.Dx(), b.Dy()
	if w <= max && h <= max {
		dst := image.NewRGBA(image.Rect(0, 0, w, h))
		draw.Draw(dst, dst.Bounds(), src, b.Min, draw.Src)
		return dst
	}
	if w > h {
		h = h * max / w
		w = max
	} else {
		w = w * max / h
		h = max
	}
	dst := image.NewRGBA(image.Rect(0, 0, w, h))
	draw.CatmullRom.Scale(dst, dst.Bounds(), src, b, draw.Over, nil)
	return dst
}

type editInput struct {
	Body string `json:"body"`
}

// Edit rewrites a message's text. Only the author, only within the edit window,
// and the result is always marked as edited — a silent edit lets someone
// rewrite what they agreed to after other people replied to it.
func (ic *IdeaController) EditMessage(c *fiber.Ctx) error {
	user := middleware.CurrentUser(c)

	msg, ok := loadMessage(c)
	if !ok {
		return nil // the helper already wrote the error response
	}
	if msg.AuthorID != user.ID {
		return c.Status(fiber.StatusForbidden).
			JSON(fiber.Map{"error": "you can only edit your own messages"})
	}
	if msg.DeletedAt != nil {
		return c.Status(fiber.StatusConflict).
			JSON(fiber.Map{"error": "that message was deleted"})
	}
	if time.Since(msg.CreatedAt) > editWindow {
		return c.Status(fiber.StatusConflict).
			JSON(fiber.Map{"error": "messages can only be edited within 24 hours"})
	}
	if msg.Kind != models.MsgText && msg.Kind != models.MsgCode {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": "attachments can't be edited — delete and repost"})
	}

	var in editInput
	if err := c.BodyParser(&in); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
	}
	body := strings.TrimSpace(in.Body)
	if body == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "message is empty"})
	}
	if len(body) > maxMessageBody {
		body = body[:maxMessageBody]
	}

	now := time.Now()
	msg.Body = body
	msg.EditedAt = &now
	database.DB.Save(msg)

	return c.JSON(fiber.Map{"message": ic.toMessageDTO(*msg, user.ID)})
}

// DeleteMessage soft-deletes, leaving a tombstone. Replies hanging off it stay
// readable, and the thread's shape is preserved.
func (ic *IdeaController) DeleteMessage(c *fiber.Ctx) error {
	user := middleware.CurrentUser(c)

	msg, ok := loadMessage(c)
	if !ok {
		return nil // the helper already wrote the error response
	}

	// The author deletes their own; an idea's owner can also remove anything
	// from their thread, because they're the one accountable for it.
	var idea models.Idea
	database.DB.First(&idea, msg.IdeaID)
	if msg.AuthorID != user.ID && idea.OwnerID != user.ID {
		return c.Status(fiber.StatusForbidden).
			JSON(fiber.Map{"error": "you can only delete your own messages"})
	}
	if msg.DeletedAt != nil {
		return c.JSON(fiber.Map{"ok": true})
	}

	now := time.Now()
	msg.DeletedAt = &now
	msg.Data = nil // reclaim the attachment bytes immediately
	msg.Mime = ""
	database.DB.Save(msg)

	touchIdea(&idea)
	return c.JSON(fiber.Map{"ok": true, "message": ic.toMessageDTO(*msg, user.ID)})
}

type reactInput struct {
	Emoji string `json:"emoji"`
}

// React toggles an emoji reaction — tap once to add, again to remove.
func (ic *IdeaController) React(c *fiber.Ctx) error {
	user := middleware.CurrentUser(c)

	msg, ok := loadMessage(c)
	if !ok {
		return nil // the helper already wrote the error response
	}
	var in reactInput
	if err := c.BodyParser(&in); err != nil || !allowedReactions[in.Emoji] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "unsupported reaction"})
	}

	var existing models.IdeaReaction
	if database.DB.Where("message_id = ? AND user_id = ? AND emoji = ?",
		msg.ID, user.ID, in.Emoji).First(&existing).Error == nil {
		database.DB.Delete(&existing)
	} else {
		database.DB.Create(&models.IdeaReaction{
			MessageID: msg.ID, UserID: user.ID, Emoji: in.Emoji,
		})
	}
	return c.JSON(fiber.Map{"message": ic.toMessageDTO(*msg, user.ID)})
}

// --- silent brainstorm -------------------------------------------------------

type brainstormInput struct {
	Minutes int    `json:"minutes"`
	Topic   string `json:"topic"`
}

// StartBrainstorm opens a timed anonymous window on an idea's thread. While it
// runs, every message posts without a name attached, so the discussion is
// judged on content rather than on who said it. When it closes, authorship
// stays hidden — revealing it afterwards would defeat the point.
func (ic *IdeaController) StartBrainstorm(c *fiber.Ctx) error {
	user := middleware.CurrentUser(c)

	idea, ok := ic.load(c)
	if !ok {
		return nil // the helper already wrote the error response
	}
	var in brainstormInput
	_ = c.BodyParser(&in)
	minutes := in.Minutes
	if minutes <= 0 || minutes > 60 {
		minutes = 10
	}

	if b := activeBrainstorm(idea.ID); b != nil {
		return c.Status(fiber.StatusConflict).
			JSON(fiber.Map{"error": "a silent brainstorm is already running"})
	}

	now := time.Now()
	session := models.BrainstormSession{
		IdeaID: &idea.ID, StartedBy: user.ID, Topic: strings.TrimSpace(in.Topic),
		StartsAt: now, EndsAt: now.Add(time.Duration(minutes) * time.Minute),
	}
	database.DB.Create(&session)

	logIdeaEvent(idea.ID, user.ID, "brainstorm", "", "",
		fmt.Sprintf("%d minutes", minutes), session.Topic)
	notifyThread(*idea, user.ID, "🤫", "Silent brainstorm started",
		fmt.Sprintf("%s opened a %d-minute anonymous brainstorm on \"%s\". Every post is unattributed.",
			displayName(*user), minutes, idea.Title))

	return c.JSON(fiber.Map{"brainstorm": brainstormDTO(&session)})
}

// StopBrainstorm ends the window early.
func (ic *IdeaController) StopBrainstorm(c *fiber.Ctx) error {
	idea, ok := ic.load(c)
	if !ok {
		return nil // the helper already wrote the error response
	}
	b := activeBrainstorm(idea.ID)
	if b == nil {
		return c.JSON(fiber.Map{"brainstorm": nil})
	}
	b.EndsAt = time.Now()
	database.DB.Save(b)
	return c.JSON(fiber.Map{"brainstorm": nil})
}

func activeBrainstorm(ideaID uint) *models.BrainstormSession {
	var b models.BrainstormSession
	if database.DB.Where("idea_id = ? AND ends_at > ?", ideaID, time.Now()).
		Order("ends_at desc").First(&b).Error != nil {
		return nil
	}
	return &b
}

func brainstormDTO(b *models.BrainstormSession) fiber.Map {
	if b == nil {
		return nil
	}
	return fiber.Map{
		"id": b.ID, "topic": b.Topic,
		"endsAt":           b.EndsAt.Format(time.RFC3339),
		"secondsRemaining": int(time.Until(b.EndsAt).Seconds()),
	}
}

// --- mentions ----------------------------------------------------------------

// notifyMentions handles both mention forms: @name pings a person, and
// @idea#12 cross-links two ideas (and tells that idea's owner they were
// referenced, which is how related work finds each other).
func notifyMentions(idea models.Idea, msg models.IdeaMessage, author models.User) {
	link := "/ideas?idea=" + strconv.Itoa(int(idea.ID)) + "&message=" + strconv.Itoa(int(msg.ID))

	// @idea#N
	for _, m := range mentionIdea.FindAllStringSubmatch(msg.Body, -1) {
		refID, err := strconv.Atoi(m[1])
		if err != nil || uint(refID) == idea.ID {
			continue
		}
		var ref models.Idea
		if database.DB.First(&ref, refID).Error != nil || ref.OwnerID == author.ID {
			continue
		}
		database.DB.Create(&models.Notification{
			UserID: ref.OwnerID,
			Key:    fmt.Sprintf("idea_ref_%d_%d", ref.ID, msg.ID),
			Kind:   "idea", Emoji: "🔗", Tint: "#17A3DD",
			Title: "Your idea was referenced",
			Body: fmt.Sprintf("%s linked \"%s\" from the discussion on \"%s\".",
				displayName(author), ref.Title, idea.Title),
			Link: link,
		})
		logIdeaEvent(idea.ID, author.ID, "linked", "", "", ref.Title, "")
	}

	// @name — matched against display names, longest first so "@ana-maria"
	// isn't swallowed by a user called "ana".
	names := mentionUser.FindAllStringSubmatch(msg.Body, -1)
	if len(names) == 0 {
		return
	}
	var users []models.User
	database.DB.Where("id != ?", author.ID).Limit(500).Find(&users)

	notified := map[uint]bool{}
	for _, m := range names {
		handle := strings.ToLower(m[1])
		if handle == "idea" {
			continue // that's the @idea#N form, already handled
		}
		for _, u := range users {
			if notified[u.ID] {
				continue
			}
			name := strings.ToLower(displayName(u))
			if name == handle || strings.HasPrefix(name, handle+" ") ||
				strings.ReplaceAll(name, " ", "-") == handle {
				notified[u.ID] = true
				database.DB.Create(&models.Notification{
					UserID: u.ID,
					Key:    fmt.Sprintf("idea_mention_%d", msg.ID),
					Kind:   "idea", Emoji: "💬", Tint: "#6C3FC5",
					Title: displayName(author) + " mentioned you",
					Body:  fmt.Sprintf("In \"%s\": %s", idea.Title, snippet(msg.Body)),
					Link:  link,
				})
			}
		}
	}
}

// --- shared ------------------------------------------------------------------

// loadMessage follows the same ok-not-error contract as IdeaController.load —
// see the note there for why.
func loadMessage(c *fiber.Ctx) (*models.IdeaMessage, bool) {
	id, err := strconv.Atoi(c.Params("messageId"))
	if err != nil {
		_ = c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid message"})
		return nil, false
	}
	var m models.IdeaMessage
	if database.DB.First(&m, id).Error != nil {
		_ = c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "message not found"})
		return nil, false
	}
	return &m, true
}

// touchIdea keeps the denormalised message count and activity stamp in step —
// the board's Hot sort reads both.
func touchIdea(idea *models.Idea) {
	var n int64
	database.DB.Model(&models.IdeaMessage{}).
		Where("idea_id = ? AND deleted_at IS NULL", idea.ID).Count(&n)
	idea.MessageCount = int(n)
	idea.LastActivityAt = time.Now()
	database.DB.Save(idea)
}
