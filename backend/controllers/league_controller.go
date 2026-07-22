package controllers

import (
	"sort"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"

	"lumora/backend/database"
	"lumora/backend/middleware"
	"lumora/backend/models"
)

// LeagueController serves the weekly competition: standings, results, and the
// two opt-outs (casual mode and reporting).
type LeagueController struct{}

// --- payloads ----------------------------------------------------------------

type leagueRow struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Points    int    `json:"points"`
	RawXP     int    `json:"rawXp"`
	Streak    int    `json:"streak"`
	Avatar    string `json:"avatarColor"`
	AvatarURL string `json:"avatarUrl"`
	Language  string `json:"language"`
	Rank      int    `json:"rank"`
	Zone      string `json:"zone"` // promote | hold | demote
	IsUser    bool   `json:"isUser"`
	Accuracy  int    `json:"accuracy"`
	Perfect   int    `json:"perfectRuns"`
	FairPlay  bool   `json:"fairPlay"`
	Flagged   bool   `json:"flagged"`
	Reported  bool   `json:"reported"` // has the viewer already reported them?
}

type tierInfo struct {
	Index        int    `json:"index"`
	Name         string `json:"name"`
	Tint         string `json:"tint"`
	PromoteTop   int    `json:"promoteTop"`
	DemoteBottom int    `json:"demoteBottom"`
	GoldGems     int    `json:"goldGems"`
	GroupBonus   int    `json:"groupBonus"`
}

func tierInfoAt(i int) tierInfo {
	i = clampTier(i)
	t := tiers[i]
	return tierInfo{i, t.Name, t.Tint, t.PromoteTop, t.DemoteBottom, t.GoldGems, t.GroupBonus}
}

// --- standings ---------------------------------------------------------------

// Standings returns the user's pod for the current season, plus everything the
// league screen renders around it: the tier ladder, the countdown, the group
// goal, and the user's own integrity state.
//
// It settles any finished season first, so the very act of opening the league is
// what closes last week out on a host that sleeps when idle.
func (lc *LeagueController) Standings(c *fiber.Ctx) error {
	SettleDueSeasons()
	// The ticker can't be relied on when the instance sleeps between visitors,
	// so opening the league also triggers the final-day nudge.
	DeliverFinalDayNudges()

	user := middleware.CurrentUser(c)
	normaliseLeagueState(user)

	now := time.Now().UTC()
	seasonID := seasonIDAt(now)
	season := ensureSeason(seasonID, now)

	base := fiber.Map{
		"seasonId":         seasonID,
		"endsAt":           season.EndsAt,
		"secondsRemaining": int(season.EndsAt.Sub(now).Seconds()),
		"tiers":            allTierInfo(),
		"tier":             tierInfoAt(user.LeagueTier),
		"podSize":          podSize,
		"you": fiber.Map{
			"integrity": user.Integrity,
			"fairPlay":  user.Integrity >= fairPlayFloor,
			"trophies":  user.Trophies,
			"best":      tierName(user.LeagueBest),
			"bestTier":  user.LeagueBest,
			"casual":    user.LeagueCasual,
		},
	}

	// Casual learners are never placed in a pod — that's the whole point.
	if user.LeagueCasual {
		base["casual"] = true
		base["joined"] = false
		base["rows"] = []leagueRow{}
		return c.JSON(base)
	}
	base["casual"] = false

	var me models.LeagueMembership
	if database.DB.Where("user_id = ? AND season_id = ?", user.ID, seasonID).
		First(&me).Error != nil {
		// The entry rule: no activity this week, no placement. Not last place —
		// simply not racing yet.
		base["joined"] = false
		base["rows"] = []leagueRow{}
		base["userRank"] = 0
		return c.JSON(base)
	}

	var pod []models.LeagueMembership
	database.DB.Where("pod_id = ?", me.PodID).Find(&pod)
	sortPod(pod)

	def := tiers[clampTier(me.Tier)]
	promote, demote := slotsFor(def, len(pod))

	// Which pod-mates has this user already reported? Drives the button state.
	reported := map[uint]bool{}
	var reports []models.LeagueReport
	database.DB.Where("season_id = ? AND reporter_id = ?", seasonID, user.ID).Find(&reports)
	for _, r := range reports {
		reported[r.SubjectID] = true
	}

	rows := make([]leagueRow, 0, len(pod))
	total, userRank := 0, 0
	for i, m := range pod {
		total += m.Points
		var u models.User
		if database.DB.First(&u, m.UserID).Error != nil {
			continue
		}
		rank := i + 1
		if m.UserID == user.ID {
			userRank = rank
		}
		acc := 0
		if m.Activities > 0 {
			acc = m.AccuracySum / m.Activities
		}
		rows = append(rows, leagueRow{
			ID: int(m.UserID), Name: displayName(u), Points: m.Points, RawXP: m.RawXP,
			Streak: u.Streak, Avatar: u.AvatarColor, AvatarURL: u.AvatarURL,
			Language: m.Language, Rank: rank, Zone: zoneFor(rank, len(pod), promote, demote, me.Tier),
			IsUser: m.UserID == user.ID, Accuracy: acc, Perfect: m.PerfectRuns,
			FairPlay: u.Integrity >= fairPlayFloor && !m.Flagged, Flagged: m.Flagged,
			Reported: reported[m.UserID],
		})
	}

	goal := def.GoalPerLearner * podSize
	base["joined"] = true
	base["rows"] = rows
	base["userRank"] = userRank
	base["promoteTop"] = promote
	base["demoteBottom"] = demote
	base["stage"] = me.Stage
	base["tier"] = tierInfoAt(me.Tier)
	base["groupGoal"] = fiber.Map{
		"target":  goal,
		"current": total,
		"hit":     total >= goal,
		"bonus":   def.GroupBonus,
	}
	base["me"] = fiber.Map{
		"points":      me.Points,
		"rawXp":       me.RawXP,
		"activities":  me.Activities,
		"perfectRuns": me.PerfectRuns,
		"flagged":     me.Flagged,
		"flagReason":  me.FlagReason,
	}
	return c.JSON(base)
}

// sortPod applies the ranking rule: points first, then whoever reached the
// total earliest.
func sortPod(pod []models.LeagueMembership) {
	sort.SliceStable(pod, func(i, j int) bool {
		if pod[i].Points != pod[j].Points {
			return pod[i].Points > pod[j].Points
		}
		return pod[i].LastPointAt.Before(pod[j].LastPointAt)
	})
}

func zoneFor(rank, size, promote, demote, tier int) string {
	if rank <= promote {
		return "promote"
	}
	if demote > 0 && rank > size-demote && tier > 0 {
		return "demote"
	}
	return "hold"
}

func allTierInfo() []tierInfo {
	out := make([]tierInfo, 0, len(tiers))
	for i := range tiers {
		out = append(out, tierInfoAt(i))
	}
	return out
}

// --- results (the end-of-race ceremony) --------------------------------------

// Result returns the user's most recent settled season if they haven't watched
// the ceremony for it yet. The frontend plays the animation off this and then
// POSTs to ResultSeen, so a result is celebrated exactly once.
func (lc *LeagueController) Result(c *fiber.Ctx) error {
	SettleDueSeasons()
	user := middleware.CurrentUser(c)

	var m models.LeagueMembership
	if database.DB.Where("user_id = ? AND settled = ? AND ceremony_seen = ?",
		user.ID, true, false).
		Order("settled_at desc").First(&m).Error != nil {
		return c.JSON(fiber.Map{"result": nil})
	}

	var pod []models.LeagueMembership
	database.DB.Where("pod_id = ?", m.PodID).Find(&pod)
	sortPod(pod)

	// The three names above the fold in the ceremony's podium.
	podium := make([]fiber.Map, 0, 3)
	for i, p := range pod {
		if i >= 3 {
			break
		}
		var u models.User
		if database.DB.First(&u, p.UserID).Error != nil {
			continue
		}
		podium = append(podium, fiber.Map{
			"rank": i + 1, "name": displayName(u), "points": p.Points,
			"avatarColor": u.AvatarColor, "avatarUrl": u.AvatarURL,
			"isUser": p.UserID == user.ID,
		})
	}

	acc := 0
	if m.Activities > 0 {
		acc = m.AccuracySum / m.Activities
	}

	return c.JSON(fiber.Map{"result": fiber.Map{
		"seasonId":     m.SeasonID,
		"result":       m.Result,
		"rank":         m.FinalRank,
		"podSize":      len(pod),
		"points":       m.Points,
		"rawXp":        m.RawXP,
		"activities":   m.Activities,
		"accuracy":     acc,
		"perfectRuns":  m.PerfectRuns,
		"gems":         m.GemsAwarded,
		"groupGoalHit": m.GroupGoalHit,
		"flagged":      m.Flagged,
		"flagReason":   m.FlagReason,
		"stage":        m.Stage,
		"from":         tierInfoAt(m.Tier),
		"to":           tierInfoAt(m.NextTier),
		"podium":       podium,
		"trophies":     user.Trophies,
	}})
}

// ResultSeen marks the ceremony watched so it doesn't replay on every visit.
func (lc *LeagueController) ResultSeen(c *fiber.Ctx) error {
	user := middleware.CurrentUser(c)
	seasonID := c.Query("season")

	q := database.DB.Model(&models.LeagueMembership{}).
		Where("user_id = ? AND settled = ?", user.ID, true)
	if seasonID != "" {
		q = q.Where("season_id = ?", seasonID)
	}
	q.Update("ceremony_seen", true)
	return c.JSON(fiber.Map{"ok": true})
}

// History returns the user's recent settled seasons for the profile/league
// timeline: what tier, what rank, what happened.
func (lc *LeagueController) History(c *fiber.Ctx) error {
	user := middleware.CurrentUser(c)

	var rows []models.LeagueMembership
	database.DB.Where("user_id = ? AND settled = ?", user.ID, true).
		Order("season_id desc").Limit(12).Find(&rows)

	out := make([]fiber.Map, 0, len(rows))
	for _, m := range rows {
		out = append(out, fiber.Map{
			"seasonId": m.SeasonID, "tier": tierInfoAt(m.Tier), "to": tierInfoAt(m.NextTier),
			"rank": m.FinalRank, "points": m.Points, "result": m.Result,
			"gems": m.GemsAwarded, "stage": m.Stage,
		})
	}
	return c.JSON(fiber.Map{"history": out})
}

// --- opt-out and reporting ---------------------------------------------------

type casualInput struct {
	Enabled bool `json:"enabled"`
}

// Casual toggles the non-competitive track. Turning it on removes the learner
// from this week's pod; nothing else about their account changes — same XP, same
// streak, same features, no leaderboard.
func (lc *LeagueController) Casual(c *fiber.Ctx) error {
	user := middleware.CurrentUser(c)

	var in casualInput
	if err := c.BodyParser(&in); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
	}
	user.LeagueCasual = in.Enabled
	database.DB.Save(user)

	if in.Enabled {
		// Withdraw from the running season so the pod isn't left with a ghost.
		database.DB.Where("user_id = ? AND season_id = ? AND settled = ?",
			user.ID, currentSeasonID(), false).
			Delete(&models.LeagueMembership{})
	}
	return c.JSON(fiber.Map{"ok": true, "casual": user.LeagueCasual})
}

type reportInput struct {
	Reason string `json:"reason"`
}

// Report is the one-tap "this doesn't look real" action. Reports are advisory:
// they raise a counter that, combined with the behavioural signals in
// detectAnomaly, can flag an account. Reports alone never punish anyone — that
// would just hand the bullies a new weapon.
func (lc *LeagueController) Report(c *fiber.Ctx) error {
	user := middleware.CurrentUser(c)

	subjectID, err := strconv.Atoi(c.Params("id"))
	if err != nil || subjectID <= 0 || uint(subjectID) == user.ID {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid target"})
	}

	var in reportInput
	_ = c.BodyParser(&in)
	seasonID := currentSeasonID()

	var existing models.LeagueReport
	if database.DB.Where("season_id = ? AND reporter_id = ? AND subject_id = ?",
		seasonID, user.ID, subjectID).First(&existing).Error == nil {
		return c.JSON(fiber.Map{"ok": true, "alreadyReported": true})
	}

	database.DB.Create(&models.LeagueReport{
		SeasonID: seasonID, ReporterID: user.ID,
		SubjectID: uint(subjectID), Reason: in.Reason,
	})

	// Three independent reports in one season is enough corroboration to flag
	// the membership for review; the account keeps its XP either way.
	var m models.LeagueMembership
	if database.DB.Where("user_id = ? AND season_id = ?", subjectID, seasonID).
		First(&m).Error == nil {
		m.ReportCount++
		if m.ReportCount >= 3 && !m.Flagged {
			m.Flagged = true
			m.FlagReason = "flagged by multiple learners for review"
		}
		database.DB.Save(&m)
	}
	return c.JSON(fiber.Map{"ok": true, "alreadyReported": false})
}
