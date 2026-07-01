package controllers

import (
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"

	"lumora/backend/database"
	"lumora/backend/middleware"
	"lumora/backend/models"
)

// HeartsController manages the hearts economy: hearts are spent on wrong
// answers, regenerate over time, and can be topped up with a purchase.
type HeartsController struct{}

const maxHearts = 5

// heartRegen is how long one heart takes to regenerate. Configurable via
// HEART_REGEN_MINUTES (default 30 minutes).
func heartRegen() time.Duration {
	m := 30
	if v := os.Getenv("HEART_REGEN_MINUTES"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			m = n
		}
	}
	return time.Duration(m) * time.Minute
}

// refreshHearts applies any hearts earned since the last checkpoint. Returns
// true if the hearts just became full again (crossed to max), so callers can
// notify. It mutates the user but does NOT save.
func refreshHearts(user *models.User) (becameFull bool) {
	if user.Hearts >= maxHearts {
		return false
	}
	if user.HeartsUpdatedAt.IsZero() {
		user.HeartsUpdatedAt = time.Now()
		return false
	}
	interval := heartRegen()
	elapsed := time.Since(user.HeartsUpdatedAt)
	gained := int(elapsed / interval)
	if gained <= 0 {
		return false
	}
	user.Hearts += gained
	if user.Hearts >= maxHearts {
		user.Hearts = maxHearts
		user.HeartsUpdatedAt = time.Time{} // full — clock not needed
		return true
	}
	// Advance the anchor by the whole hearts consumed so partial progress is kept.
	user.HeartsUpdatedAt = user.HeartsUpdatedAt.Add(time.Duration(gained) * interval)
	return false
}

// secondsToNextHeart reports how long until the next heart regenerates (0 when
// full).
func secondsToNextHeart(user models.User) int {
	if user.Hearts >= maxHearts || user.HeartsUpdatedAt.IsZero() {
		return 0
	}
	remaining := time.Until(user.HeartsUpdatedAt.Add(heartRegen()))
	if remaining < 0 {
		return 0
	}
	return int(remaining.Seconds())
}

func heartsPayload(user models.User) fiber.Map {
	rate := 130
	if v := os.Getenv("KES_PER_USD"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			rate = n
		}
	}
	usd := float64(heartsRefillPriceKES) / float64(rate)
	usd = float64(int(usd*100+0.5)) / 100

	return fiber.Map{
		"hearts":          user.Hearts,
		"max":             maxHearts,
		"full":            user.Hearts >= maxHearts,
		"secondsToNext":   secondsToNextHeart(user),
		"regenMinutes":    int(heartRegen().Minutes()),
		"paymentsEnabled": PaymentsEnabled(),
		"refillPriceKes":  heartsRefillPriceKES,
		"refillPriceUsd":  usd,
	}
}

// Status returns the current hearts, applying any pending regeneration first.
func (hc *HeartsController) Status(c *fiber.Ctx) error {
	user := middleware.CurrentUser(c)
	if refreshHearts(user) {
		DeliverHeartsFull(user.ID)
	}
	database.DB.Save(user)
	return c.JSON(heartsPayload(*user))
}

// Lose spends one heart (on a wrong answer). Regenerates first so the count is
// current, then decrements and starts the regen clock if hearts were full.
func (hc *HeartsController) Lose(c *fiber.Ctx) error {
	user := middleware.CurrentUser(c)
	refreshHearts(user)

	if user.Hearts > 0 {
		wasFull := user.Hearts >= maxHearts
		user.Hearts--
		if wasFull {
			user.HeartsUpdatedAt = time.Now() // start the regen clock now
		}
		if user.Hearts == 0 {
			DeliverHeartsEmpty(user.ID, secondsToNextHeart(*user), int(heartRegen().Minutes()))
		}
	}
	database.DB.Save(user)
	return c.JSON(heartsPayload(*user))
}

// grantFullHearts refills to max (used after a hearts purchase).
func grantFullHearts(userID uint) {
	var user models.User
	if database.DB.First(&user, userID).Error != nil {
		return
	}
	user.Hearts = maxHearts
	user.HeartsUpdatedAt = time.Time{}
	database.DB.Save(&user)
}
