package controllers

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"

	"lumora/backend/database"
	"lumora/backend/middleware"
	"lumora/backend/models"
)

// ChatController powers 1:1 direct messaging between users.
type ChatController struct{}

type chatUser struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	AvatarColor string `json:"avatarColor"`
	AvatarURL   string `json:"avatarUrl"`
	LevelName   string `json:"levelName"`
}

func toChatUser(u models.User) chatUser {
	return chatUser{
		ID: u.ID, Name: displayName(u), AvatarColor: u.AvatarColor,
		AvatarURL: u.AvatarURL, LevelName: u.LevelName,
	}
}

func snippet(s string) string {
	s = strings.TrimSpace(s)
	if len(s) > 80 {
		return s[:80] + "…"
	}
	return s
}

// Contacts lists everyone the current user can start a chat with.
func (cc *ChatController) Contacts(c *fiber.Ctx) error {
	user := middleware.CurrentUser(c)
	var users []models.User
	database.DB.Where("id != ?", user.ID).Order("name asc").Limit(100).Find(&users)
	out := make([]chatUser, 0, len(users))
	for _, u := range users {
		out = append(out, toChatUser(u))
	}
	return c.JSON(fiber.Map{"contacts": out})
}

type threadDTO struct {
	User        chatUser `json:"user"`
	LastMessage string   `json:"lastMessage"`
	LastAt      string   `json:"lastAt"`
	Unread      int      `json:"unread"`
}

// Threads returns the user's conversations, newest first.
func (cc *ChatController) Threads(c *fiber.Ctx) error {
	user := middleware.CurrentUser(c)
	uid := user.ID

	var msgs []models.Message
	database.DB.Where("sender_id = ? OR recipient_id = ?", uid, uid).
		Order("created_at desc").Find(&msgs)

	order := []uint{}
	last := map[uint]models.Message{}
	unread := map[uint]int{}
	for _, m := range msgs {
		other := m.SenderID
		if m.SenderID == uid {
			other = m.RecipientID
		}
		if _, seen := last[other]; !seen {
			last[other] = m // first seen = most recent (desc order)
			order = append(order, other)
		}
		if m.RecipientID == uid && !m.Read {
			unread[other]++
		}
	}

	threads := make([]threadDTO, 0, len(order))
	for _, other := range order {
		var u models.User
		if database.DB.First(&u, other).Error != nil {
			continue
		}
		m := last[other]
		prefix := ""
		if m.SenderID == uid {
			prefix = "You: "
		}
		threads = append(threads, threadDTO{
			User:        toChatUser(u),
			LastMessage: prefix + previewOf(m),
			LastAt:      m.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			Unread:      unread[other],
		})
	}
	return c.JSON(fiber.Map{"threads": threads})
}

// Unread returns the total number of unread messages.
func (cc *ChatController) Unread(c *fiber.Ctx) error {
	user := middleware.CurrentUser(c)
	var n int64
	database.DB.Model(&models.Message{}).
		Where("recipient_id = ? AND read = ?", user.ID, false).Count(&n)
	return c.JSON(fiber.Map{"count": n})
}

// Messages returns the conversation with another user and marks it read.
func (cc *ChatController) Messages(c *fiber.Ctx) error {
	user := middleware.CurrentUser(c)
	otherID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user"})
	}

	var other models.User
	if database.DB.First(&other, otherID).Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
	}

	var msgs []models.Message
	database.DB.Where(
		"(sender_id = ? AND recipient_id = ?) OR (sender_id = ? AND recipient_id = ?)",
		user.ID, otherID, otherID, user.ID,
	).Order("created_at asc").Find(&msgs)

	out := make([]chatMessageDTO, 0, len(msgs))
	for _, m := range msgs {
		out = append(out, toChatMessage(m, user.ID))
	}

	// Mark messages from the other person as read.
	database.DB.Model(&models.Message{}).
		Where("sender_id = ? AND recipient_id = ? AND read = ?", otherID, user.ID, false).
		Update("read", true)
	// Clear the chat notification for this conversation.
	database.DB.Model(&models.Notification{}).
		Where("user_id = ? AND key = ?", user.ID, "chat_"+strconv.Itoa(otherID)).
		Update("read", true)

	return c.JSON(fiber.Map{"messages": out, "user": toChatUser(other)})
}

type sendInput struct {
	Body string `json:"body"`
}

// Send delivers a message and notifies the recipient.
func (cc *ChatController) Send(c *fiber.Ctx) error {
	user := middleware.CurrentUser(c)
	otherID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user"})
	}
	var in sendInput
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

	var other models.User
	if database.DB.First(&other, otherID).Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
	}

	msg := models.Message{
		SenderID: user.ID, RecipientID: uint(otherID),
		Kind: models.MsgText, Body: body,
	}
	database.DB.Create(&msg)

	// Label the language now (free, offline) and translate in the background
	// if it isn't English — the sender shouldn't wait on a model round-trip
	// to see their own message posted.
	msg.DetectedLang = detectAndQueue("messages", msg.ID, msg.Body)

	upsertChatNotification(user, uint(otherID), snippet(body))

	return c.JSON(fiber.Map{"message": toChatMessage(msg, user.ID)})
}
