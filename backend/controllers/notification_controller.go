package controllers

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"

	"lumora/backend/database"
	"lumora/backend/middleware"
	"lumora/backend/models"
)

// NotificationController serves a user's in-app notifications.
type NotificationController struct{}

// List returns the user's notifications (newest first) and the unread count.
func (nc *NotificationController) List(c *fiber.Ctx) error {
	user := middleware.CurrentUser(c)

	var items []models.Notification
	database.DB.Where("user_id = ?", user.ID).
		Order("created_at desc").Limit(50).Find(&items)

	var unread int64
	database.DB.Model(&models.Notification{}).
		Where("user_id = ? AND read = ?", user.ID, false).Count(&unread)

	return c.JSON(fiber.Map{"notifications": items, "unread": unread})
}

// MarkRead marks all of the user's notifications as read.
func (nc *NotificationController) MarkRead(c *fiber.Ctx) error {
	user := middleware.CurrentUser(c)
	database.DB.Model(&models.Notification{}).
		Where("user_id = ? AND read = ?", user.ID, false).
		Update("read", true)
	return c.JSON(fiber.Map{"ok": true})
}

// --- delivery ----------------------------------------------------------------

type campaignItem struct {
	Key, Kind, Emoji, Tint, Title, Body string
}

// The rotating pool of automated messages. Keys are stable so each is delivered
// to a given user at most once.
var campaigns = []campaignItem{
	{"feature_speaking", "feature", "🗣️", "#6C3FC5", "New: Speaking practice", "Say phrases out loud and get a live fluency score in the Practice tab."},
	{"tip_daily", "tip", "🎯", "#F5A623", "Tip: Little and often", "A few minutes every day beats one long cram session. Hit your daily goal!"},
	{"feature_listening", "feature", "🎧", "#00C2A8", "New: Listening sessions", "Hear short conversations voiced by your companions, then answer questions."},
	{"lang_de_fr", "language", "🌍", "#17A3DD", "German & French are live", "Tap the language switcher on Learn or Profile to add a new course anytime."},
	{"tip_speak", "tip", "💬", "#6C3FC5", "Tip: Speak it out loud", "Saying answers aloud — not just reading — cements vocabulary far faster."},
	{"feature_reading", "feature", "📖", "#17A3DD", "New: Reading passages", "Every unit now has a reading session to train your eyes on real text."},
	{"tip_review", "tip", "🔁", "#00C2A8", "Tip: Review your mistakes", "The Practice tab turns anything you miss into quick, targeted drills."},
	{"lang_soon", "language", "🇯🇵", "#FF5C5C", "Coming soon: Japanese & Italian", "More languages are on the way — keep your streak warm for launch!"},
	{"tip_streak", "tip", "🔥", "#FF5C5C", "Tip: Protect your streak", "One lesson a day keeps your flame alive. Don't let it go out!"},
}

// ensureNotification creates a notification for a user unless one with the same
// key already exists (dedup).
func ensureNotification(userID uint, item campaignItem) bool {
	if item.Key != "" {
		var existing models.Notification
		if database.DB.Where("user_id = ? AND key = ?", userID, item.Key).
			First(&existing).Error == nil {
			return false
		}
	}
	database.DB.Create(&models.Notification{
		UserID: userID, Key: item.Key, Kind: item.Kind,
		Emoji: item.Emoji, Tint: item.Tint, Title: item.Title, Body: item.Body,
	})
	return true
}

// DeliverWelcome sends the one-off welcome notification on sign-up.
func DeliverWelcome(userID uint) {
	ensureNotification(userID, campaignItem{
		Key: "welcome", Kind: "welcome", Emoji: "🦊", Tint: "#6C3FC5",
		Title: "Welcome to Lumora! 🎉",
		Body:  "I'm Lumora, your guide. Finish your first lesson to start a streak — you've got this!",
	})
}

// runCampaign delivers, to each user, the next campaign message they haven't
// received yet — but only if their most recent notification is older than `gap`,
// so messages trickle in casually instead of arriving all at once.
func runCampaign(gap time.Duration) {
	var users []models.User
	database.DB.Find(&users)

	for _, u := range users {
		var latest models.Notification
		if database.DB.Where("user_id = ?", u.ID).
			Order("created_at desc").First(&latest).Error == nil {
			if time.Since(latest.CreatedAt) < gap {
				continue // delivered something recently — let it breathe
			}
		}
		for _, item := range campaigns {
			if ensureNotification(u.ID, item) {
				break // one new message per user per run
			}
		}
	}
}

// StartNotificationScheduler runs the automated push loop in the background.
func StartNotificationScheduler() {
	interval := 120
	if v := os.Getenv("NOTIF_INTERVAL_SECONDS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			interval = n
		}
	}
	gap := time.Duration(interval) * time.Second

	go func() {
		// A short initial delay so the first batch lands soon after boot.
		time.Sleep(15 * time.Second)
		runSafely(gap)
		t := time.NewTicker(gap)
		defer t.Stop()
		for range t.C {
			runSafely(gap)
		}
	}()
	log.Printf("[notifications] scheduler started (every %ds)", interval)
}

func runSafely(gap time.Duration) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[notifications] campaign run recovered: %v", r)
		}
	}()
	runCampaign(gap)
}
