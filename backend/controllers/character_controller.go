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
	ID        int    `json:"id"`
	Name      string `json:"name"`
	XP        int    `json:"xp"`
	Streak    int    `json:"streak"`
	Avatar    string `json:"avatarColor"`
	AvatarURL string `json:"avatarUrl"` // uploaded profile photo, if any
	Language  string `json:"language"`  // language they're learning (badge)
	IsUser    bool   `json:"isUser"`
	Rank      int    `json:"rank"`
}

// League returns a global ranking of all learners by total XP, each tagged with
// the language they're studying. The current user is always included.
func (lc *LeaderboardController) League(c *fiber.Ctx) error {
	user := middleware.CurrentUser(c)

	// Keep the user's league badge in sync with their XP.
	if l := leagueForXP(user.XP); user.League != l {
		user.League = l
		database.DB.Save(user)
	}

	var users []models.User
	database.DB.Order("xp desc, id asc").Limit(100).Find(&users)

	rows := make([]leaderRow, 0, len(users)+1)
	included := false
	for _, u := range users {
		if u.ID == user.ID {
			included = true
		}
		rows = append(rows, leaderRow{
			ID: int(u.ID), Name: displayName(u), XP: u.XP, Streak: u.Streak,
			Avatar: u.AvatarColor, AvatarURL: u.AvatarURL,
			Language: u.TargetLanguage, IsUser: u.ID == user.ID,
		})
	}

	// Always show the current user, even if they're beyond the top 100.
	if !included {
		rows = append(rows, leaderRow{
			ID: int(user.ID), Name: displayName(*user), XP: user.XP, Streak: user.Streak,
			Avatar: user.AvatarColor, AvatarURL: user.AvatarURL,
			Language: user.TargetLanguage, IsUser: true,
		})
	}

	// Sort by XP desc and assign ranks.
	for i := 1; i < len(rows); i++ {
		for j := i; j > 0 && rows[j].XP > rows[j-1].XP; j-- {
			rows[j], rows[j-1] = rows[j-1], rows[j]
		}
	}
	userRank := 0
	for i := range rows {
		rows[i].Rank = i + 1
		if rows[i].IsUser {
			userRank = rows[i].Rank
		}
	}

	return c.JSON(fiber.Map{
		"league":   user.League,
		"rows":     rows,
		"userRank": userRank,
	})
}

func displayName(u models.User) string {
	if u.Name != "" {
		return u.Name
	}
	return "Learner"
}
