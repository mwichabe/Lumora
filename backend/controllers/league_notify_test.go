package controllers

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"lumora/backend/database"
	"lumora/backend/models"
)

// notes returns a user's league notifications, newest first.
func notes(t *testing.T, userID uint) []models.Notification {
	t.Helper()
	var out []models.Notification
	database.DB.Where("user_id = ? AND kind = ?", userID, "league").
		Order("id desc").Find(&out)
	return out
}

func hasNote(t *testing.T, userID uint, keyPrefix string) *models.Notification {
	t.Helper()
	for _, n := range notes(t, userID) {
		if strings.HasPrefix(n.Key, keyPrefix) {
			return &n
		}
	}
	return nil
}

// --- joining -----------------------------------------------------------------

func TestJoiningTheRaceNotifiesOnce(t *testing.T) {
	newTestDB(t)
	u := makeUser(t, "join@test.dev", 2)

	for i := 0; i < 3; i++ {
		AwardLeaguePoints(u, LeagueAward{Source: "lesson", RawXP: 20, Accuracy: 100})
	}

	var count int
	for _, n := range notes(t, u.ID) {
		if strings.HasPrefix(n.Key, "league_joined_") {
			count++
		}
	}
	if count != 1 {
		t.Fatalf("got %d join notifications after 3 lessons, want exactly 1", count)
	}

	n := hasNote(t, u.ID, "league_joined_")
	if n.Link != "/leaderboard" {
		t.Errorf("league notification link = %q, want /leaderboard", n.Link)
	}
}

// --- zone transitions --------------------------------------------------------

// The core rule: notifications fire when the user *crosses* a boundary, not
// every time they score while on one side of it.
func TestZoneChangeNotifiesOnTransitionOnly(t *testing.T) {
	newTestDB(t)

	// A Gold pod of 12: promotion is the top 6 (capped at half the pod).
	rivals := make([]*models.User, 0, 11)
	for i := 0; i < 11; i++ {
		r := makeUser(t, fmt.Sprintf("rival%d@test.dev", i), 3)
		r.LeagueTier = 2
		database.DB.Save(r)
		AwardLeaguePoints(r, LeagueAward{Source: "lesson", RawXP: 100, Accuracy: 100, Difficulty: 2})
		rivals = append(rivals, r)
	}

	// Our learner joins last with a single small lesson: bottom of the pod.
	u := makeUser(t, "climber@test.dev", 3)
	u.LeagueTier = 2
	database.DB.Save(u)
	AwardLeaguePoints(u, LeagueAward{Source: "lesson", RawXP: 5, Accuracy: 60, Difficulty: 1})

	if hasNote(t, u.ID, "league_zone_up_") != nil {
		t.Fatal("bottom of the pod should not be told it's in the promotion zone")
	}

	// Now out-work the pod and cross into the promotion zone.
	for i := 0; i < 6; i++ {
		AwardLeaguePoints(u, LeagueAward{Source: "lesson", RawXP: 80, Accuracy: 100, Difficulty: 2})
	}

	if hasNote(t, u.ID, "league_zone_up_") == nil {
		t.Fatalf("crossing into the promotion zone was not notified; notes=%v", keysOf(notes(t, u.ID)))
	}

	// Scoring again while already inside the zone must not re-notify.
	before := len(notes(t, u.ID))
	AwardLeaguePoints(u, LeagueAward{Source: "lesson", RawXP: 40, Accuracy: 100, Difficulty: 2})
	if after := len(notes(t, u.ID)); after != before {
		t.Errorf("scoring inside the promotion zone added %d notifications, want 0", after-before)
	}
}

func TestFallingOutOfThePromotionZoneNotifies(t *testing.T) {
	newTestDB(t)
	leader := makeUser(t, "leader@test.dev", 3)
	chaser := makeUser(t, "chaser@test.dev", 3)

	// A two-person pod: promotion is the top 1.
	AwardLeaguePoints(leader, LeagueAward{Source: "lesson", RawXP: 40, Accuracy: 100})
	AwardLeaguePoints(chaser, LeagueAward{Source: "lesson", RawXP: 10, Accuracy: 100})

	// Being top of a pod that's still filling isn't news — the join note already
	// covered the promotion band, and a second notification in the same breath
	// would just be noise.
	if hasNote(t, leader.ID, "league_zone_up_") != nil {
		t.Error("leading a one-person pod triggered a promotion-zone notification")
	}

	// The chaser overtakes.
	AwardLeaguePoints(chaser, LeagueAward{Source: "lesson", RawXP: 200, Accuracy: 100, Difficulty: 2})
	// The leader scores again, and discovers they've slipped.
	AwardLeaguePoints(leader, LeagueAward{Source: "lesson", RawXP: 5, Accuracy: 100})

	if hasNote(t, leader.ID, "league_zone_lost_") == nil {
		t.Errorf("slipping out of the promotion zone was not notified; notes=%v",
			keysOf(notes(t, leader.ID)))
	}
}

// Joining late into an established pod can drop you straight into the demotion
// band. There's no previous zone to have moved from, but it's still the one
// thing that learner most needs to know.
func TestJoiningStraightIntoTheDropZoneNotifies(t *testing.T) {
	newTestDB(t)

	// A Silver pod of 24 with real scores on the board. The size matters: a
	// demotion band only exists once a pod has enough people for the bottom to
	// mean something (see slotsFor), which for Silver is 20.
	for i := 0; i < 24; i++ {
		r := makeUser(t, fmt.Sprintf("est%d@test.dev", i), 5)
		r.LeagueTier = 1
		database.DB.Save(r)
		AwardLeaguePoints(r, LeagueAward{Source: "lesson", RawXP: 200, Accuracy: 100, Difficulty: 2})
	}

	latecomer := makeUser(t, "late-joiner@test.dev", 1)
	latecomer.LeagueTier = 1
	database.DB.Save(latecomer)
	AwardLeaguePoints(latecomer, LeagueAward{Source: "lesson", RawXP: 5, Accuracy: 50, Difficulty: 1})

	if hasNote(t, latecomer.ID, "league_zone_down_") == nil {
		t.Errorf("landing in the drop zone on entry was not notified; notes=%v",
			keysOf(notes(t, latecomer.ID)))
	}
}

// --- overtaking --------------------------------------------------------------

// The learner who *loses* a place is the one who needs to hear about it.
func TestBeingOvertakenNotifiesTheLearnerWhoLostThePlace(t *testing.T) {
	newTestDB(t)
	held := makeUser(t, "held@test.dev", 3)
	mover := makeUser(t, "mover@test.dev", 3)

	AwardLeaguePoints(held, LeagueAward{Source: "lesson", RawXP: 50, Accuracy: 100})
	AwardLeaguePoints(mover, LeagueAward{Source: "lesson", RawXP: 10, Accuracy: 100})
	AwardLeaguePoints(mover, LeagueAward{Source: "lesson", RawXP: 300, Accuracy: 100, Difficulty: 2})

	n := hasNote(t, held.ID, "league_passed_")
	if n == nil {
		t.Fatalf("overtaken learner was not notified; notes=%v", keysOf(notes(t, held.ID)))
	}
	if !strings.Contains(n.Title, displayName(*mover)) {
		t.Errorf("notification %q doesn't name who passed them", n.Title)
	}
	// The learner doing the overtaking shouldn't be told they were overtaken.
	if hasNote(t, mover.ID, "league_passed_") != nil {
		t.Error("the mover was told they were passed")
	}
}

// A pod trading places all evening must produce one nudge, not forty.
func TestOvertakenNotificationIsRateLimited(t *testing.T) {
	newTestDB(t)
	held := makeUser(t, "held2@test.dev", 3)
	mover := makeUser(t, "mover2@test.dev", 3)

	AwardLeaguePoints(held, LeagueAward{Source: "lesson", RawXP: 50, Accuracy: 100})
	AwardLeaguePoints(mover, LeagueAward{Source: "lesson", RawXP: 10, Accuracy: 100})

	// Trade the lead back and forth several times.
	for i := 0; i < 4; i++ {
		AwardLeaguePoints(mover, LeagueAward{Source: "lesson", RawXP: 200, Accuracy: 100, Difficulty: 2})
		AwardLeaguePoints(held, LeagueAward{Source: "lesson", RawXP: 200, Accuracy: 100, Difficulty: 2})
	}

	var passed int
	for _, n := range notes(t, held.ID) {
		if strings.HasPrefix(n.Key, "league_passed_") {
			passed++
		}
	}
	if passed > 1 {
		t.Errorf("got %d 'you were passed' notifications inside the cooldown, want 1", passed)
	}
}

func TestTakingTheLeadNotifies(t *testing.T) {
	newTestDB(t)
	a := makeUser(t, "a1@test.dev", 3)
	b := makeUser(t, "b1@test.dev", 3)

	AwardLeaguePoints(a, LeagueAward{Source: "lesson", RawXP: 100, Accuracy: 100})
	AwardLeaguePoints(b, LeagueAward{Source: "lesson", RawXP: 10, Accuracy: 100})
	AwardLeaguePoints(b, LeagueAward{Source: "lesson", RawXP: 400, Accuracy: 100, Difficulty: 2})

	if hasNote(t, b.ID, "league_lead_") == nil {
		t.Errorf("taking #1 was not notified; notes=%v", keysOf(notes(t, b.ID)))
	}
}

// --- collaborative goal ------------------------------------------------------

func TestGroupGoalNotifiesEveryContributor(t *testing.T) {
	newTestDB(t)

	// Bronze needs 120 x 30 = 3,600 points. Three learners grinding hard clear
	// it between them.
	members := make([]*models.User, 0, 3)
	for i := 0; i < 3; i++ {
		u := makeUser(t, fmt.Sprintf("goal%d@test.dev", i), 20)
		members = append(members, u)
	}
	for round := 0; round < 6; round++ {
		for _, u := range members {
			AwardLeaguePoints(u, LeagueAward{
				Source: "lesson", RawXP: 400, Accuracy: 100, Difficulty: 2,
			})
		}
	}

	for i, u := range members {
		if hasNote(t, u.ID, "league_goal_") == nil {
			t.Errorf("member %d was not told the pod goal was hit; notes=%v",
				i, keysOf(notes(t, u.ID)))
		}
	}

	// And only once, however much more the pod scores.
	AwardLeaguePoints(members[0], LeagueAward{Source: "lesson", RawXP: 200, Accuracy: 100})
	var goals int
	for _, n := range notes(t, members[0].ID) {
		if strings.HasPrefix(n.Key, "league_goal_") {
			goals++
		}
	}
	if goals != 1 {
		t.Errorf("got %d pod-goal notifications, want 1", goals)
	}
}

// --- the final-day nudge -----------------------------------------------------

func TestFinalDayNudgeTellsYouWhatYouNeed(t *testing.T) {
	newTestDB(t)

	// Move the current season's end to two hours from now.
	seasonID := currentSeasonID()
	now := time.Now().UTC()
	start, _ := seasonBounds(now)
	database.DB.Create(&models.LeagueSeason{
		ID: seasonID, StartsAt: start, EndsAt: now.Add(2 * time.Hour),
	})

	leader := makeUser(t, "fd-leader@test.dev", 3)
	trailer := makeUser(t, "fd-trailer@test.dev", 3)
	AwardLeaguePoints(leader, LeagueAward{Source: "lesson", RawXP: 300, Accuracy: 100, Difficulty: 2})
	AwardLeaguePoints(trailer, LeagueAward{Source: "lesson", RawXP: 10, Accuracy: 100})

	DeliverFinalDayNudges()

	up := hasNote(t, leader.ID, "league_final_")
	if up == nil {
		t.Fatal("the leader got no final-day nudge")
	}
	if !strings.Contains(up.Title, "promoting") {
		t.Errorf("leader's nudge = %q, expected it to say they're promoting", up.Title)
	}

	down := hasNote(t, trailer.ID, "league_final_")
	if down == nil {
		t.Fatal("the trailing learner got no final-day nudge")
	}
	// A two-person pod has no demotion slots, so the trailer gets the "here's
	// the gap" variant — which must actually quantify the gap.
	if !strings.Contains(down.Body, "points away") {
		t.Errorf("trailer's nudge doesn't quantify the gap: %q", down.Body)
	}

	// Sending twice must not duplicate it.
	DeliverFinalDayNudges()
	var count int
	for _, n := range notes(t, leader.ID) {
		if strings.HasPrefix(n.Key, "league_final_") {
			count++
		}
	}
	if count != 1 {
		t.Errorf("got %d final-day nudges, want 1", count)
	}
}

// The nudge is for the last day only — it must stay silent mid-week.
func TestNoFinalDayNudgeEarlyInTheWeek(t *testing.T) {
	newTestDB(t)

	seasonID := currentSeasonID()
	now := time.Now().UTC()
	start, _ := seasonBounds(now)
	database.DB.Create(&models.LeagueSeason{
		ID: seasonID, StartsAt: start, EndsAt: now.Add(72 * time.Hour),
	})

	u := makeUser(t, "midweek@test.dev", 3)
	AwardLeaguePoints(u, LeagueAward{Source: "lesson", RawXP: 50, Accuracy: 100})

	DeliverFinalDayNudges()

	if hasNote(t, u.ID, "league_final_") != nil {
		t.Error("final-day nudge fired with three days left")
	}
}

// --- settlement --------------------------------------------------------------

// Promotion and demotion each produce a notification naming the destination
// tier, deep-linked to the league screen.
func TestSettlementNotifiesPromotionAndDemotion(t *testing.T) {
	newTestDB(t)
	season := seedFinishedSeason(t)
	pod := seedPod(t, season.ID, 2 /* Gold */, podSize, "")

	SettleDueSeasons()

	winner := hasNote(t, pod[0], "league_"+season.ID)
	if winner == nil {
		t.Fatal("the pod winner got no settlement notification")
	}
	if !strings.Contains(winner.Title, "Promoted") || !strings.Contains(winner.Title, "Sapphire") {
		t.Errorf("winner's notification = %q, want a promotion to Sapphire", winner.Title)
	}
	if winner.Link != "/leaderboard" {
		t.Errorf("settlement notification link = %q, want /leaderboard", winner.Link)
	}
	if !strings.Contains(winner.Body, "gems") {
		t.Errorf("winner's notification doesn't mention the chest: %q", winner.Body)
	}

	loser := hasNote(t, pod[podSize-1], "league_"+season.ID)
	if loser == nil {
		t.Fatal("the bottom finisher got no settlement notification")
	}
	if !strings.Contains(loser.Title, "Silver") {
		t.Errorf("bottom finisher's notification = %q, want a move down to Silver", loser.Title)
	}
	// Demotion copy has to be encouraging, not punitive — it's the message most
	// likely to make someone quit.
	if !strings.Contains(loser.Body, "climb straight back") {
		t.Errorf("demotion copy isn't encouraging: %q", loser.Body)
	}
}

func keysOf(ns []models.Notification) []string {
	out := make([]string, 0, len(ns))
	for _, n := range ns {
		out = append(out, n.Key)
	}
	return out
}
