package controllers

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"

	"lumora/backend/database"
	"lumora/backend/middleware"
	"lumora/backend/models"
)

// Rich direct messages: photo sharing, editing and deleting.
//
// The rules mirror the ideas thread deliberately, because a user shouldn't have
// to learn two different sets: you can edit your own text for 24 hours and the
// result is always marked as edited; deleting leaves a tombstone rather than
// rewriting history the other person has already read.

// chatMessageDTO is what the client actually renders. The raw model is no
// longer safe to serialise directly — it carries attachment bytes.
type chatMessageDTO struct {
	ID          uint   `json:"id"`
	SenderID    uint   `json:"senderId"`
	RecipientID uint   `json:"recipientId"`
	Kind        string `json:"kind"`
	Body        string `json:"body"`
	URL         string `json:"url"`
	FileName    string `json:"fileName"`
	Width       int    `json:"width"`
	Height      int    `json:"height"`
	Read        bool   `json:"read"`
	Mine        bool   `json:"mine"`
	Edited      bool   `json:"edited"`
	Deleted     bool   `json:"deleted"`
	CanEdit     bool   `json:"canEdit"`
	CreatedAt   string `json:"createdAt"`

	// Present only when the message isn't English — see translation.go.
	Translation *translationDTO `json:"translation,omitempty"`
}

func toChatMessage(m models.Message, viewerID uint) chatMessageDTO {
	kind := m.Kind
	if kind == "" {
		kind = models.MsgText // rows written before attachments existed
	}
	dto := chatMessageDTO{
		ID: m.ID, SenderID: m.SenderID, RecipientID: m.RecipientID,
		Kind: kind, Body: m.Body, FileName: m.FileName,
		Width: m.Width, Height: m.Height, Read: m.Read,
		Mine:      m.SenderID == viewerID,
		Edited:    m.EditedAt != nil,
		Deleted:   m.DeletedAt != nil,
		CreatedAt: m.CreatedAt.Format(time.RFC3339),
	}
	if m.DeletedAt != nil {
		dto.Body = ""
		dto.Kind = models.MsgText
		return dto
	}
	if len(m.Data) > 0 {
		dto.URL = fmt.Sprintf("/api/chat/attachments/%d", m.ID)
	}
	dto.CanEdit = m.SenderID == viewerID && kind == models.MsgText &&
		time.Since(m.CreatedAt) < editWindow
	dto.Translation = translationFor(m.Body, m.DetectedLang, m.TranslatedBody, m.TranslatedAt)
	return dto
}

// SendImage posts a photo into a conversation (multipart "file"). The image is
// validated by decoding it, downscaled, and stored in the database.
func (cc *ChatController) SendImage(c *fiber.Ctx) error {
	user := middleware.CurrentUser(c)

	otherID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user"})
	}
	var other models.User
	if database.DB.First(&other, otherID).Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
	}

	fh, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "no image uploaded"})
	}
	if fh.Size > maxAttachment {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": "image must be under 8MB"})
	}
	f, err := fh.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"error": "could not read that image"})
	}
	defer f.Close()

	src, _, err := image.Decode(f)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": "could not read that image"})
	}
	scaled := downscale(src, maxImageEdge)
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, scaled, &jpeg.Options{Quality: 85}); err != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"error": "could not process that image"})
	}

	b := scaled.Bounds()
	msg := models.Message{
		SenderID: user.ID, RecipientID: uint(otherID),
		Kind: models.MsgImage, Body: strings.TrimSpace(c.FormValue("body")),
		Data: buf.Bytes(), Mime: "image/jpeg", FileName: fh.Filename,
		Width: b.Dx(), Height: b.Dy(),
	}
	if len(msg.Body) > 2000 {
		msg.Body = msg.Body[:2000]
	}
	if err := database.DB.Create(&msg).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"error": "could not send image"})
	}
	// A caption in another language gets translated like any other message.
	msg.DetectedLang = detectAndQueue("messages", msg.ID, msg.Body)

	body := "📷 Photo"
	if msg.Body != "" {
		body = "📷 " + snippet(msg.Body)
	}
	upsertChatNotification(user, uint(otherID), body)

	return c.Status(fiber.StatusCreated).
		JSON(fiber.Map{"message": toChatMessage(msg, user.ID)})
}

// ChatAttachment serves a message's image. Unauthenticated for the same reason
// avatars are — a plain <img> can't attach a bearer token — and keyed by an
// opaque row id.
func (cc *ChatController) ChatAttachment(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("bad id")
	}
	var m models.Message
	if database.DB.First(&m, id).Error != nil || len(m.Data) == 0 || m.DeletedAt != nil {
		return c.Status(fiber.StatusNotFound).SendString("not found")
	}
	c.Set("Content-Type", m.Mime)
	c.Set("Cache-Control", "public, max-age=31536000, immutable")
	return c.Send(m.Data)
}

// EditMessage rewrites a sent message. Sender only, text only, 24-hour window,
// and always marked as edited.
func (cc *ChatController) EditMessage(c *fiber.Ctx) error {
	user := middleware.CurrentUser(c)

	msg, ok := loadChatMessage(c)
	if !ok {
		return nil // the helper already wrote the error response
	}
	if msg.SenderID != user.ID {
		return c.Status(fiber.StatusForbidden).
			JSON(fiber.Map{"error": "you can only edit your own messages"})
	}
	if msg.DeletedAt != nil {
		return c.Status(fiber.StatusConflict).
			JSON(fiber.Map{"error": "that message was deleted"})
	}
	if msg.Kind == models.MsgImage {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": "photos can't be edited — delete and send again"})
	}
	if time.Since(msg.CreatedAt) > editWindow {
		return c.Status(fiber.StatusConflict).
			JSON(fiber.Map{"error": "messages can only be edited within 24 hours"})
	}

	var in editInput
	if err := c.BodyParser(&in); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
	}
	body := strings.TrimSpace(in.Body)
	if body == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "message is empty"})
	}
	if len(body) > 2000 {
		body = body[:2000]
	}

	now := time.Now()
	msg.Body = body
	msg.EditedAt = &now
	database.DB.Save(msg)

	return c.JSON(fiber.Map{"message": toChatMessage(*msg, user.ID)})
}

// DeleteMessage soft-deletes so the conversation keeps its shape — the other
// person still sees that something was said and removed, rather than a reply
// that suddenly answers nothing.
func (cc *ChatController) DeleteMessage(c *fiber.Ctx) error {
	user := middleware.CurrentUser(c)

	msg, ok := loadChatMessage(c)
	if !ok {
		return nil // the helper already wrote the error response
	}
	if msg.SenderID != user.ID {
		return c.Status(fiber.StatusForbidden).
			JSON(fiber.Map{"error": "you can only delete your own messages"})
	}
	if msg.DeletedAt != nil {
		return c.JSON(fiber.Map{"ok": true})
	}

	now := time.Now()
	msg.DeletedAt = &now
	msg.Data = nil
	msg.Mime = ""
	database.DB.Save(msg)

	return c.JSON(fiber.Map{"ok": true, "message": toChatMessage(*msg, user.ID)})
}

// previewOf renders a message for the conversation list, where a photo or a
// removed message has no body text to show.
func previewOf(m models.Message) string {
	switch {
	case m.DeletedAt != nil:
		return "Message deleted"
	case m.Kind == models.MsgImage && strings.TrimSpace(m.Body) == "":
		return "📷 Photo"
	case m.Kind == models.MsgImage:
		return "📷 " + snippet(m.Body)
	default:
		return snippet(m.Body)
	}
}

// loadChatMessage follows the same ok-not-error contract as IdeaController.load
// — see the note there for why.
func loadChatMessage(c *fiber.Ctx) (*models.Message, bool) {
	id, err := strconv.Atoi(c.Params("messageId"))
	if err != nil {
		_ = c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid message"})
		return nil, false
	}
	var m models.Message
	if database.DB.First(&m, id).Error != nil {
		_ = c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "message not found"})
		return nil, false
	}
	return &m, true
}

// upsertChatNotification keeps one notification per conversation rather than
// one per message, so a burst of photos doesn't bury everything else.
func upsertChatNotification(sender *models.User, recipientID uint, body string) {
	key := "chat_" + strconv.FormatUint(uint64(sender.ID), 10)
	title := "New message from " + displayName(*sender)

	var existing models.Notification
	if database.DB.Where("user_id = ? AND key = ?", recipientID, key).
		First(&existing).Error == nil {
		existing.Title = title
		existing.Body = body
		existing.Read = false
		database.DB.Save(&existing)
		return
	}
	database.DB.Create(&models.Notification{
		UserID: recipientID, Key: key, Kind: "chat",
		Emoji: "💬", Tint: "#6C3FC5", Title: title, Body: body,
		Link: "/chat/" + strconv.FormatUint(uint64(sender.ID), 10),
	})
}
