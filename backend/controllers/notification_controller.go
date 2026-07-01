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

// MarkOneRead marks a single notification (the one the user opened) as read.
func (nc *NotificationController) MarkOneRead(c *fiber.Ctx) error {
	user := middleware.CurrentUser(c)
	id := c.Params("id")

	var n models.Notification
	if err := database.DB.Where("id = ? AND user_id = ?", id, user.ID).
		First(&n).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "notification not found"})
	}
	if !n.Read {
		database.DB.Model(&n).Update("read", true)
	}

	var unread int64
	database.DB.Model(&models.Notification{}).
		Where("user_id = ? AND read = ?", user.ID, false).Count(&unread)
	return c.JSON(fiber.Map{"ok": true, "unread": unread})
}

// Delete removes a single notification the user owns.
func (nc *NotificationController) Delete(c *fiber.Ctx) error {
	user := middleware.CurrentUser(c)
	id := c.Params("id")
	res := database.DB.Where("id = ? AND user_id = ?", id, user.ID).
		Delete(&models.Notification{})
	if res.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "notification not found"})
	}
	var unread int64
	database.DB.Model(&models.Notification{}).
		Where("user_id = ? AND read = ?", user.ID, false).Count(&unread)
	return c.JSON(fiber.Map{"ok": true, "unread": unread})
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
		Title: "Welcome to Lumora!",
		Body:  "I'm Lumora, your guide. Finish your first lesson to start a streak — you've got this!",
	})
}

// DeliverLoginWelcome greets a returning user by name. To avoid spamming on
// frequent logins/session-resumes, at most one is created every few hours.
// Returns true when it actually created one (so callers can also email).
func DeliverLoginWelcome(user models.User) bool {
	var recent int64
	database.DB.Model(&models.Notification{}).
		Where("user_id = ? AND kind = ? AND created_at > ?",
			user.ID, "welcome_back", time.Now().Add(-3*time.Hour)).
		Count(&recent)
	if recent > 0 {
		return false
	}
	name := displayName(user)
	database.DB.Create(&models.Notification{
		UserID: user.ID, Kind: "welcome_back", Emoji: "👋", Tint: "#6C3FC5",
		Title: "Welcome back, " + name + "!",
		Body:  "Great to see you again. Keep your streak alive with a quick lesson today.",
	})
	return true
}

// DeliverUnitComplete congratulates a user for finishing a whole unit (once).
func DeliverUnitComplete(userID uint, unit string) {
	ensureNotification(userID, campaignItem{
		Key: "unit_done_" + unit, Kind: "milestone", Emoji: "🏆", Tint: "#00C2A8",
		Title: "Unit complete!",
		Body:  "You finished " + unit + ". Brilliant work — the next unit is unlocked. Keep the momentum going!",
	})
}

// DeliverLevelUp celebrates reaching a new CEFR level (once per level).
func DeliverLevelUp(userID uint, cefr, levelName string) {
	ensureNotification(userID, campaignItem{
		Key: "levelup_" + cefr, Kind: "milestone", Emoji: "⭐", Tint: "#F5A623",
		Title: "Level up! You reached " + cefr,
		Body:  "You've advanced to " + levelName + " (" + cefr + "). Your hard work is paying off!",
	})
}

// DeliverExamStarted notifies the user that they've begun an exam attempt. Not
// deduped — every attempt gets its own notification (empty Key).
func DeliverExamStarted(userID uint, level, langName string) {
	ensureNotification(userID, campaignItem{
		Kind: "exam", Emoji: "📝", Tint: "#6C3FC5",
		Title: "Exam started — " + level,
		Body:  "You've started the " + level + " " + langName + " exam. Stay focused and good luck!",
	})
}

// DeliverExamPassed congratulates the user and links to the certificate.
func DeliverExamPassed(userID uint, level string, score int, certLink string) {
	database.DB.Create(&models.Notification{
		UserID: userID, Kind: "exam", Emoji: "🎓", Tint: "#00C2A8",
		Title: "You passed the " + level + " exam!",
		Body: "Congratulations! You scored " + strconv.Itoa(score) +
			"%. Your certificate is ready to view and download.",
		Link: certLink,
	})
}

// DeliverExamFailed encourages the user and points them back to the exam.
func DeliverExamFailed(userID uint, level string, score, passMark int) {
	database.DB.Create(&models.Notification{
		UserID: userID, Kind: "exam", Emoji: "📚", Tint: "#F5A623",
		Title: level + " exam — not passed this time",
		Body: "You scored " + strconv.Itoa(score) + "% (pass mark " +
			strconv.Itoa(passMark) + "%). Don't give up — review and retake when you're ready.",
		Link: "/exam",
	})
}

// DeliverHeartsEmpty tells the user they've run out of hearts.
func DeliverHeartsEmpty(userID uint, secondsToNext, regenMinutes int) {
	database.DB.Create(&models.Notification{
		UserID: userID, Kind: "hearts", Emoji: "💔", Tint: "#FF5C5C",
		Title: "You're out of hearts",
		Body: "A new heart regenerates about every " + strconv.Itoa(regenMinutes) +
			" minutes. Wait for one, or refill instantly to keep learning.",
		Link: "/learn",
	})
}

// DeliverHeartsFull tells the user their hearts have fully regenerated.
func DeliverHeartsFull(userID uint) {
	database.DB.Create(&models.Notification{
		UserID: userID, Kind: "hearts", Emoji: "❤️", Tint: "#00C2A8",
		Title: "Your hearts are full!",
		Body:  "All five hearts have regenerated. Jump back in and keep your streak alive!",
		Link:  "/learn",
	})
}

// DeliverHeartsPurchased confirms a hearts refill purchase.
func DeliverHeartsPurchased(userID uint) {
	database.DB.Create(&models.Notification{
		UserID: userID, Kind: "hearts", Emoji: "❤️", Tint: "#00C2A8",
		Title: "Hearts refilled",
		Body:  "Thanks! Your hearts are full again. Happy learning!",
		Link:  "/learn",
	})
}

// DeliverPaymentSuccess tells the user their payment went through and the exam
// is unlocked (deduped per transaction reference).
func DeliverPaymentSuccess(userID uint, reference string) {
	ensureNotification(userID, campaignItem{
		Key: "pay_ok_" + reference, Kind: "payment", Emoji: "✅", Tint: "#00C2A8",
		Title: "Payment successful",
		Body:  "Your payment went through — your exam attempt is ready. Head to the exam whenever you like. Good luck!",
	})
}

// DeliverPaymentFailed tells the user a payment didn't complete (deduped per
// transaction reference).
func DeliverPaymentFailed(userID uint, reference string) {
	ensureNotification(userID, campaignItem{
		Key: "pay_fail_" + reference, Kind: "payment", Emoji: "⚠️", Tint: "#FF5C5C",
		Title: "Payment not completed",
		Body:  "We couldn't confirm your payment. If you were charged it will be applied automatically — otherwise please try again.",
	})
}

// DeliverStreakMilestone celebrates streak milestones (7, 30, 100… once each).
func DeliverStreakMilestone(userID, days int) {
	if days != 3 && days != 7 && days != 14 && days != 30 && days != 100 {
		return
	}
	ensureNotification(uint(userID), campaignItem{
		Key: "streak_" + strconv.Itoa(days), Kind: "streak", Emoji: "🔥", Tint: "#FF5C5C",
		Title: strconv.Itoa(days) + "-day streak!",
		Body:  "You've practised " + strconv.Itoa(days) + " days in a row. Don't break the chain!",
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
