package controllers

import (
	"fmt"
	"time"

	"lumora/backend/database"
	"lumora/backend/models"
)

// ============================================================================
// League notifications
// ============================================================================
//
// Settlement notifications tell a learner how the week ended. These tell them
// what is happening *during* it — which is the half that actually drives
// participation, because a leaderboard nobody checks may as well not exist.
//
// The events, and what makes each one worth an interruption:
//
//	joined      you're racing, here's who against
//	promotion   you're in the promotion zone (or you just fell out of it)
//	demotion    you're in the drop zone (or you just climbed clear)
//	lead        you took #1
//	overtaken   someone passed you
//	pod goal    the collaborative target was hit
//	final day   24 hours left, here's exactly what you need
//
// Two rules keep this from becoming spam. Every notification fires on a
// *transition*, never on a state — LeagueMembership remembers the last rank and
// zone we notified about. And every one carries a cooldown, so a pod where two
// people trade places all evening produces one nudge, not forty.

const (
	// A zone flip can genuinely happen twice in a day; more than that is noise.
	zoneFlipCooldown = 6 * time.Hour
	// Being overtaken is the most frequent event in a busy pod, so it's the
	// quietest.
	overtakenCooldown = 12 * time.Hour
	leadCooldown      = 24 * time.Hour
)

// leagueNote creates a league notification, deep-linked to the league screen.
//
// Unlike ensureNotification (which dedups on the key forever, right for
// once-in-a-lifetime milestones) this suppresses only within `cooldown`, so a
// recurring event can fire again next week — or later the same week — without
// needing a new key each time. A zero cooldown means once, ever.
func leagueNote(userID uint, key, emoji, tint, title, body string, cooldown time.Duration) bool {
	q := database.DB.Model(&models.Notification{}).Where("user_id = ? AND key = ?", userID, key)
	if cooldown > 0 {
		q = q.Where("created_at > ?", time.Now().Add(-cooldown))
	}
	var recent int64
	q.Count(&recent)
	if recent > 0 {
		return false
	}

	database.DB.Create(&models.Notification{
		UserID: userID, Key: key, Kind: "league",
		Emoji: emoji, Tint: tint, Title: title, Body: body,
		Link: "/leaderboard",
	})
	return true
}

// --- the in-week watcher -----------------------------------------------------

// evaluateLeagueMoments runs after every scoring event and turns the change in
// the user's standing into notifications. It is the only place in-week league
// notifications originate.
//
// Cost is one indexed query for the pod (30 rows at most), which is why it's
// safe to call on every lesson completion.
func evaluateLeagueMoments(user *models.User, m *models.LeagueMembership, gained int) {
	if m == nil || user == nil {
		return
	}
	var pod []models.LeagueMembership
	database.DB.Where("pod_id = ?", m.PodID).Find(&pod)
	if len(pod) == 0 {
		return
	}
	sortPod(pod)

	def := tiers[clampTier(m.Tier)]
	promote, demote := slotsFor(def, len(pod))

	rank := 0
	for i, p := range pod {
		if p.UserID == user.ID {
			rank = i + 1
			break
		}
	}
	if rank == 0 {
		return
	}
	zone := zoneFor(rank, len(pod), promote, demote, m.Tier)

	prevRank, prevZone := m.LastRank, m.LastZone

	// 1. Welcome to the race — sent once, the first time they score in a season.
	if !m.JoinNotified {
		m.JoinNotified = true
		leagueNote(user.ID, "league_joined_"+m.SeasonID, "🏁", def.Tint,
			"You're in this week's "+def.Name+" race",
			fmt.Sprintf("You've been placed in a pod of learners matched to your level. Top %d move up to %s; the race closes Monday.",
				promote, tierName(minInt(m.Tier+1, diamondTier))),
			0)
	}

	// 2. Zone transitions — the events that actually change what you should do.
	//
	// On the very first evaluation there's no previous zone to have moved from.
	// Being in the promotion zone at that point isn't news (the join note
	// already explains the promotion band, and a pod that's still filling puts
	// almost everyone there), but landing straight in the drop zone is — that
	// happens when someone joins late into an established pod, and it's the one
	// case where they need to act immediately.
	if zone != prevZone && (prevZone != "" || zone == "demote") {
		notifyZoneChange(user, m, def, zone, prevZone, rank, len(pod), promote)
	}

	// 3. Taking the lead.
	if rank == 1 && prevRank != 1 && len(pod) > 1 {
		leagueNote(user.ID, "league_lead_"+m.SeasonID, "👑", def.Tint,
			"You're #1 in "+def.Name+" League",
			fmt.Sprintf("Top of your pod of %d with %d points. The gold chest is yours if you hold it to Monday.",
				len(pod), m.Points),
			leadCooldown)
	}

	// 4. Whoever this award pushed down a place.
	if gained > 0 && prevRank > 0 && rank < prevRank {
		notifyOvertaken(pod, rank, user, def)
	}

	// 5. The pod goal, announced to everyone who contributed.
	notifyGroupGoal(pod, def, m.SeasonID)

	m.LastRank, m.LastZone = rank, zone
	database.DB.Save(m)
}

func notifyZoneChange(
	user *models.User, m *models.LeagueMembership, def tierDef,
	zone, prevZone string, rank, size, promote int,
) {
	up := tierName(minInt(m.Tier+1, diamondTier))
	down := tierName(maxInt(m.Tier-1, 0))

	switch {
	case zone == "promote":
		leagueNote(user.ID, "league_zone_up_"+m.SeasonID, "🚀", "#00C2A8",
			"You're in the promotion zone",
			fmt.Sprintf("#%d of %d — inside the top %d. Hold this until Monday and you move up to %s.",
				rank, size, promote, up),
			zoneFlipCooldown)

	case prevZone == "promote":
		leagueNote(user.ID, "league_zone_lost_"+m.SeasonID, "⚡", "#F5A623",
			"You've slipped out of the promotion zone",
			fmt.Sprintf("You're #%d and the top %d go up. One good lesson usually takes it back — and points earned earlier in the week count for more.",
				rank, promote),
			zoneFlipCooldown)

	case zone == "demote":
		leagueNote(user.ID, "league_zone_down_"+m.SeasonID, "🌧️", "#FF5C5C",
			"You're in the drop zone",
			fmt.Sprintf("#%d of %d in %s League. The bottom few move down to %s when the week closes — there's still time to climb out.",
				rank, size, def.Name, down),
			zoneFlipCooldown)

	case prevZone == "demote":
		leagueNote(user.ID, "league_zone_safe_"+m.SeasonID, "🛟", "#00C2A8",
			"You're out of the drop zone",
			fmt.Sprintf("Back up to #%d of %d. Keep it steady and you'll hold your place in %s.",
				rank, size, def.Name),
			zoneFlipCooldown)
	}
}

// notifyOvertaken tells the learner who just lost a place. Only the person
// directly displaced hears about it: telling everyone below would turn one
// lesson into a dozen push notifications.
func notifyOvertaken(pod []models.LeagueMembership, newRank int, mover *models.User, def tierDef) {
	if newRank >= len(pod) {
		return
	}
	passed := pod[newRank] // the row now sitting one place below the mover
	if passed.UserID == mover.ID {
		return
	}

	var victim models.User
	if database.DB.First(&victim, passed.UserID).Error != nil || victim.LeagueCasual {
		return
	}
	leagueNote(victim.ID, "league_passed_"+passed.SeasonID, "👀", def.Tint,
		displayName(*mover)+" just passed you",
		fmt.Sprintf("You've dropped to #%d in %s League. %s is on %d points — close enough to take back today.",
			newRank+1, def.Name, displayName(*mover), pod[newRank-1].Points),
		overtakenCooldown)
}

// notifyGroupGoal announces the collaborative target the moment it's reached,
// to everyone in the pod who put points on the board.
func notifyGroupGoal(pod []models.LeagueMembership, def tierDef, seasonID string) {
	total := 0
	for _, p := range pod {
		total += p.Points
	}
	goal := def.GoalPerLearner * podSize
	if total < goal {
		return
	}
	for _, p := range pod {
		if p.Points <= 0 {
			continue
		}
		leagueNote(p.UserID, "league_goal_"+p.PodID, "🎯", "#00C2A8",
			"Your pod smashed its goal",
			fmt.Sprintf("Between you, your pod cleared %d points. Everyone who scored collects a %d-gem bonus at Monday's reset — win or lose.",
				goal, def.GroupBonus),
			0)
	}
}

// --- the final-day nudge -----------------------------------------------------

// DeliverFinalDayNudges tells everyone still racing exactly where they stand
// with a day to go, and what it would take to change it. This is the single
// highest-value league notification: it arrives at the one moment when the
// information is still actionable.
//
// Called from the scheduler and, on a host that sleeps, from the league screen —
// so a learner who opens the app on Sunday still gets it.
func DeliverFinalDayNudges() {
	now := time.Now().UTC()
	seasonID := seasonIDAt(now)

	var season models.LeagueSeason
	if database.DB.Where("id = ?", seasonID).First(&season).Error != nil {
		return
	}
	left := season.EndsAt.Sub(now)
	if left <= 0 || left > 24*time.Hour {
		return
	}

	var pending []models.LeagueMembership
	database.DB.Where("season_id = ? AND ending_notice = ? AND points > 0",
		seasonID, false).Find(&pending)
	if len(pending) == 0 {
		return
	}

	// Group by pod so each pod is ranked once rather than once per member.
	byPod := map[string][]models.LeagueMembership{}
	for _, m := range pending {
		byPod[m.PodID] = append(byPod[m.PodID], m)
	}

	for podID := range byPod {
		var pod []models.LeagueMembership
		database.DB.Where("pod_id = ?", podID).Find(&pod)
		if len(pod) == 0 {
			continue
		}
		sortPod(pod)
		def := tiers[clampTier(pod[0].Tier)]
		promote, demote := slotsFor(def, len(pod))

		for i := range pod {
			m := &pod[i]
			if m.EndingNotice || m.Points <= 0 {
				continue
			}
			rank := i + 1
			hours := int(left.Hours())
			if hours < 1 {
				hours = 1
			}

			title, body := finalDayCopy(pod, m, def, rank, promote, demote, hours)
			leagueNote(m.UserID, "league_final_"+m.SeasonID, "⏳", def.Tint, title, body, 0)

			m.EndingNotice = true
			database.DB.Model(m).Update("ending_notice", true)
		}
	}
}

// finalDayCopy writes the nudge for one learner: their position, and the gap in
// points to the place that would change their week.
func finalDayCopy(
	pod []models.LeagueMembership, m *models.LeagueMembership, def tierDef,
	rank, promote, demote, hours int,
) (string, string) {
	size := len(pod)
	up := tierName(minInt(m.Tier+1, diamondTier))
	down := tierName(maxInt(m.Tier-1, 0))
	inDrop := demote > 0 && rank > size-demote && m.Tier > 0

	switch {
	case rank <= promote:
		// Holding a promotion spot: name the threat behind them.
		gap := 0
		if rank < size {
			gap = m.Points - pod[rank].Points
		}
		return fmt.Sprintf("%dh left — you're promoting to %s", hours, up),
			fmt.Sprintf("You're #%d of %d in %s League with %d points%s. Finish the week here and you move up.",
				rank, size, def.Name, m.Points,
				gapPhrase(gap, "and only %d points ahead of the learner behind you"))

	case inDrop:
		safe := pod[maxInt(size-demote-1, 0)].Points
		return fmt.Sprintf("%dh left — you're in the drop zone", hours),
			fmt.Sprintf("You're #%d of %d and the bottom %d move down to %s. You need %d points to climb clear. There's still time.",
				rank, size, demote, down, maxInt(safe-m.Points+1, 1))

	default:
		need := pod[maxInt(promote-1, 0)].Points - m.Points + 1
		return fmt.Sprintf("%dh left in the %s race", hours, def.Name),
			fmt.Sprintf("You're #%d of %d. Promotion is %d points away — roughly %s. Every point still counts.",
				rank, size, maxInt(need, 1), lessonsPhrase(maxInt(need, 1)))
	}
}

func gapPhrase(gap int, format string) string {
	if gap <= 0 {
		return ""
	}
	return " " + fmt.Sprintf(format, gap)
}

// lessonsPhrase turns a point gap into something a learner can act on. A clean
// lesson is worth roughly 20 points once weighting is applied.
func lessonsPhrase(points int) string {
	n := (points + 19) / 20
	if n <= 1 {
		return "one good lesson"
	}
	if n > 12 {
		return "a serious push"
	}
	return fmt.Sprintf("%d lessons", n)
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
