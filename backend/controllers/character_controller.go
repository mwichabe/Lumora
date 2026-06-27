package controllers

import (
	"github.com/gofiber/fiber/v2"

	"lumora/backend/database"
	"lumora/backend/middleware"
	"lumora/backend/models"
)

// CharacterController exposes the Lumora cast and friendship levels.
type CharacterController struct{}

// List returns every character plus the user's friendship level with each.
func (cc *CharacterController) List(c *fiber.Ctx) error {
	user := middleware.CurrentUser(c)

	var characters []models.Character
	database.DB.Find(&characters)

	var friendships []models.Friendship
	database.DB.Where("user_id = ?", user.ID).Find(&friendships)
	levelByChar := map[uint]int{}
	for _, f := range friendships {
		levelByChar[f.CharacterID] = f.Level
	}

	type charWithFriendship struct {
		models.Character
		FriendshipLevel int `json:"friendshipLevel"`
	}
	out := make([]charWithFriendship, 0, len(characters))
	for _, ch := range characters {
		lvl := levelByChar[ch.ID]
		if lvl == 0 {
			lvl = 1
		}
		out = append(out, charWithFriendship{Character: ch, FriendshipLevel: lvl})
	}

	return c.JSON(fiber.Map{"characters": out})
}

// LeaderboardController serves league standings.
type LeaderboardController struct{}

type leaderRow struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	XP       int    `json:"xp"`
	Streak   int    `json:"streak"`
	Avatar   string `json:"avatarColor"`
	IsUser   bool   `json:"isUser"`
	Rank     int    `json:"rank"`
}

// League returns the top users by total XP, marking the current user's row.
// To make a fresh database feel alive, a few synthetic rivals are blended in.
func (lc *LeaderboardController) League(c *fiber.Ctx) error {
	user := middleware.CurrentUser(c)

	var users []models.User
	database.DB.Order("xp desc").Limit(20).Find(&users)

	rows := make([]leaderRow, 0, len(users))
	for _, u := range users {
		rows = append(rows, leaderRow{
			ID: u.ID, Name: displayName(u), XP: u.XP, Streak: u.Streak,
			Avatar: u.AvatarColor, IsUser: u.ID == user.ID,
		})
	}

	// Blend in lovable rivals so an early leaderboard isn't lonely.
	rivals := []leaderRow{
		{Name: "Riko 🐼", XP: 240, Streak: 12, Avatar: "#F5A623"},
		{Name: "Cora 🐙", XP: 180, Streak: 7, Avatar: "#00C2A8"},
		{Name: "Zephyr 🦅", XP: 95, Streak: 4, Avatar: "#17A3DD"},
	}
	rows = append(rows, rivals...)

	// Sort by XP desc (simple insertion since the list is tiny).
	for i := 1; i < len(rows); i++ {
		for j := i; j > 0 && rows[j].XP > rows[j-1].XP; j-- {
			rows[j], rows[j-1] = rows[j-1], rows[j]
		}
	}
	for i := range rows {
		rows[i].Rank = i + 1
	}

	return c.JSON(fiber.Map{
		"league": user.League,
		"rows":   rows,
	})
}

func displayName(u models.User) string {
	if u.Name != "" {
		return u.Name
	}
	return "Learner"
}
