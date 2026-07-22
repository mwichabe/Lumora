package controllers

import (
	"fmt"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"lumora/backend/database"
	"lumora/backend/models"
)

// newTestDB gives each test its own in-memory database with the league tables
// migrated. It replaces the package-level handle the engine reads, which is
// fine because tests in a package run sequentially by default.
func newTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared&_pragma=busy_timeout(5000)"),
		&gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatalf("open test db: %v", err)
	}
	if err := db.Migrator().DropTable(&models.User{}, &models.Notification{},
		&models.LeagueSeason{}, &models.LeagueMembership{},
		&models.LeagueDaily{}, &models.LeagueReport{}); err != nil {
		t.Fatalf("drop: %v", err)
	}
	if err := db.AutoMigrate(&models.User{}, &models.Notification{},
		&models.LeagueSeason{}, &models.LeagueMembership{},
		&models.LeagueDaily{}, &models.LeagueReport{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	database.DB = db
	return db
}

// --- scoring -----------------------------------------------------------------

func TestTaperPointsBands(t *testing.T) {
	cases := []struct {
		name          string
		already, gain int
		want          float64
	}{
		{"inside the full-value band", 0, 500, 500},
		{"straddles full and half", 400, 200, 100 + 50},
		{"entirely in the half band", 500, 500, 250},
		{"deep in the quarter band", 1000, 400, 100},
		{"one award spanning all three", 0, 1400, 500 + 250 + 100},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := taperPoints(c.already, c.gain); got != c.want {
				t.Errorf("taperPoints(%d, %d) = %v, want %v", c.already, c.gain, got, c.want)
			}
		})
	}
}

// Monday XP has to be worth measurably more than Sunday XP, or the Sunday
// sprint is still the optimal strategy.
func TestEarlinessFavoursTheStartOfTheWeek(t *testing.T) {
	start := time.Date(2026, 7, 20, 0, 0, 0, 0, time.UTC) // a Monday
	mon := earlinessWeight(start.Add(2*time.Hour), start)
	sun := earlinessWeight(start.AddDate(0, 0, 6), start)

	if mon <= sun {
		t.Fatalf("Monday (%v) should outweigh Sunday (%v)", mon, sun)
	}
	if mon < 1.19 || mon > 1.21 {
		t.Errorf("Monday weight = %v, want ~1.20", mon)
	}
	if sun < 0.84 || sun > 0.86 {
		t.Errorf("Sunday weight = %v, want ~0.85", sun)
	}
}

func TestAccuracyWeightRewardsFlawlessRuns(t *testing.T) {
	sloppy := accuracyWeight(50)
	clean := accuracyWeight(90)
	perfect := accuracyWeight(100)

	if !(sloppy < clean && clean < perfect) {
		t.Fatalf("expected 50%% < 90%% < 100%%, got %v %v %v", sloppy, clean, perfect)
	}
	// The deep-practice bonus should make a flawless run a visible step up, not
	// a rounding difference.
	if perfect/clean < 1.1 {
		t.Errorf("flawless bonus too small: perfect=%v clean=%v", perfect, clean)
	}
}

func TestDifficultyWeightScalesWithDepth(t *testing.T) {
	if difficultyWeight(0) != 1.0 {
		t.Errorf("beginner content should be 1.0x, got %v", difficultyWeight(0))
	}
	if difficultyWeight(5000) != 2.0 {
		t.Errorf("deepest content should be 2.0x, got %v", difficultyWeight(5000))
	}
	if difficultyWeight(700) <= difficultyWeight(150) {
		t.Error("later content should outweigh earlier content")
	}
}

// --- earning -----------------------------------------------------------------

func TestAwardCreatesMembershipAndScores(t *testing.T) {
	newTestDB(t)
	u := makeUser(t, "a@test.dev", 5)

	points := AwardLeaguePoints(u, LeagueAward{
		Source: "lesson", RawXP: 20, Accuracy: 100, Difficulty: 1.5,
	})
	if points <= 0 {
		t.Fatalf("expected points, got %d", points)
	}

	var m models.LeagueMembership
	if err := database.DB.Where("user_id = ?", u.ID).First(&m).Error; err != nil {
		t.Fatalf("membership was not created: %v", err)
	}
	if m.Points != points || m.RawXP != 20 || m.Activities != 1 || m.PerfectRuns != 1 {
		t.Errorf("membership mis-recorded: %+v", m)
	}
	if m.PodID == "" {
		t.Error("membership was not assigned a pod")
	}
}

// Casual mode is a real opt-out: no membership row, no placement, nothing to
// lose at the weekend.
func TestCasualLearnersAreNeverPlaced(t *testing.T) {
	newTestDB(t)
	u := makeUser(t, "casual@test.dev", 3)
	u.LeagueCasual = true
	database.DB.Save(u)

	if got := AwardLeaguePoints(u, LeagueAward{Source: "lesson", RawXP: 50, Accuracy: 100}); got != 0 {
		t.Errorf("casual learner scored %d points, want 0", got)
	}
	var n int64
	database.DB.Model(&models.LeagueMembership{}).Where("user_id = ?", u.ID).Count(&n)
	if n != 0 {
		t.Errorf("casual learner was placed in %d pods, want 0", n)
	}
}

// Drills pay XP but stop moving the leaderboard after the daily cap.
func TestTimedDrillsStopScoringAtTheDailyCap(t *testing.T) {
	newTestDB(t)
	u := makeUser(t, "grinder@test.dev", 1)

	scoredRuns := 0
	for i := 0; i < dailyTimedCap+4; i++ {
		if AwardLeaguePoints(u, LeagueAward{Source: "practice", RawXP: 10, Accuracy: 100}) > 0 {
			scoredRuns++
		}
	}
	if scoredRuns != dailyTimedCap {
		t.Errorf("%d drills scored, want the cap of %d", scoredRuns, dailyTimedCap)
	}

	// The XP still lands — only the league stops counting.
	var d models.LeagueDaily
	database.DB.Where("user_id = ?", u.ID).First(&d)
	if d.RawXP != (dailyTimedCap+4)*10 {
		t.Errorf("daily raw XP = %d, want %d", d.RawXP, (dailyTimedCap+4)*10)
	}
}

// A day of relentless grinding must be worth far less per-XP than a normal one.
func TestGrindingHitsDiminishingReturns(t *testing.T) {
	newTestDB(t)
	u := makeUser(t, "bot@test.dev", 0)

	first := AwardLeaguePoints(u, LeagueAward{Source: "lesson", RawXP: 400, Accuracy: 100, Difficulty: 1})
	// Now well past the taper thresholds.
	for i := 0; i < 3; i++ {
		AwardLeaguePoints(u, LeagueAward{Source: "lesson", RawXP: 400, Accuracy: 100, Difficulty: 1})
	}
	last := AwardLeaguePoints(u, LeagueAward{Source: "lesson", RawXP: 400, Accuracy: 100, Difficulty: 1})

	if last >= first/2 {
		t.Errorf("late-in-the-day award (%d) should be far below the first (%d)", last, first)
	}
}

func TestImplausibleVolumeIsFlagged(t *testing.T) {
	newTestDB(t)
	u := makeUser(t, "auto@test.dev", 0)

	for i := 0; i < 10; i++ {
		AwardLeaguePoints(u, LeagueAward{Source: "lesson", RawXP: 400, Accuracy: 100, Difficulty: 1})
	}

	var m models.LeagueMembership
	database.DB.Where("user_id = ?", u.ID).First(&m)
	if !m.Flagged {
		t.Fatal("4,000 XP in one day should have been flagged")
	}
	var got models.User
	database.DB.First(&got, u.ID)
	if got.Integrity >= 100 {
		t.Errorf("integrity should have been docked, still %d", got.Integrity)
	}
}

// --- placement ---------------------------------------------------------------

func TestPodsFillToThirtyThenSplit(t *testing.T) {
	newTestDB(t)
	season := currentSeasonID()

	for i := 0; i < podSize+5; i++ {
		u := makeUser(t, fmt.Sprintf("pod%d@test.dev", i), 0)
		AwardLeaguePoints(u, LeagueAward{Source: "lesson", RawXP: 10, Accuracy: 100})
	}

	var members []models.LeagueMembership
	database.DB.Where("season_id = ?", season).Find(&members)
	counts := map[string]int{}
	for _, m := range members {
		counts[m.PodID]++
	}
	if len(counts) != 2 {
		t.Fatalf("expected 2 pods for %d learners, got %d", podSize+5, len(counts))
	}
	for id, n := range counts {
		if n > podSize {
			t.Errorf("pod %s holds %d learners, over the cap of %d", id, n, podSize)
		}
	}
}

// Pods are language-segregated, so XP from an easy second language can't be
// stacked against learners of a different one.
func TestPodsAreSegregatedByLanguage(t *testing.T) {
	newTestDB(t)

	es := makeUser(t, "es@test.dev", 0)
	de := makeUser(t, "de@test.dev", 0)
	de.TargetLanguage = "de"
	database.DB.Save(de)

	AwardLeaguePoints(es, LeagueAward{Source: "lesson", RawXP: 10, Accuracy: 100})
	AwardLeaguePoints(de, LeagueAward{Source: "lesson", RawXP: 10, Accuracy: 100})

	var a, b models.LeagueMembership
	database.DB.Where("user_id = ?", es.ID).First(&a)
	database.DB.Where("user_id = ?", de.ID).First(&b)
	if a.PodID == b.PodID {
		t.Errorf("learners of different languages shared pod %s", a.PodID)
	}
}

// --- settlement --------------------------------------------------------------

func TestSettlementPromotesDemotesAndPays(t *testing.T) {
	newTestDB(t)
	season := seedFinishedSeason(t)
	pod := seedPod(t, season.ID, 2 /* Gold */, podSize, "")

	SettleDueSeasons()

	def := tiers[2]
	for i, id := range pod {
		var m models.LeagueMembership
		database.DB.Where("user_id = ? AND season_id = ?", id, season.ID).First(&m)
		rank := i + 1

		if !m.Settled {
			t.Fatalf("rank %d was never settled", rank)
		}
		if m.FinalRank != rank {
			t.Errorf("member seeded at rank %d finished %d", rank, m.FinalRank)
		}

		var u models.User
		database.DB.First(&u, id)
		switch {
		case rank <= def.PromoteTop:
			if m.Result != "promoted" || u.LeagueTier != 3 {
				t.Errorf("rank %d: got %q tier %d, want promoted to 3", rank, m.Result, u.LeagueTier)
			}
		case rank > podSize-def.DemoteBottom:
			if m.Result != "demoted" || u.LeagueTier != 1 {
				t.Errorf("rank %d: got %q tier %d, want demoted to 1", rank, m.Result, u.LeagueTier)
			}
		default:
			if m.Result != "held" {
				t.Errorf("rank %d: got %q, want held", rank, m.Result)
			}
		}
	}

	// Podium chests, and nothing for fourth.
	gems := func(rank int) int {
		var m models.LeagueMembership
		database.DB.Where("user_id = ? AND season_id = ?", pod[rank-1], season.ID).First(&m)
		return m.GemsAwarded
	}
	if !(gems(1) > gems(2) && gems(2) > gems(3)) {
		t.Errorf("chests should descend: %d %d %d", gems(1), gems(2), gems(3))
	}
	if gems(4) > def.GroupBonus {
		t.Errorf("fourth place earned %d gems beyond the group bonus", gems(4))
	}
}

// Bronze is the floor: there is nowhere below it, so nobody is demoted out.
func TestBronzeNeverDemotes(t *testing.T) {
	newTestDB(t)
	season := seedFinishedSeason(t)
	pod := seedPod(t, season.ID, 0, podSize, "")

	SettleDueSeasons()

	for _, id := range pod {
		var m models.LeagueMembership
		database.DB.Where("user_id = ? AND season_id = ?", id, season.ID).First(&m)
		if m.Result == "demoted" {
			t.Fatalf("user %d was demoted out of Bronze", id)
		}
	}
}

// Equal totals are ranked by who reached them first — the reward for not
// waiting until the deadline.
func TestTiesGoToWhoeverGotThereFirst(t *testing.T) {
	newTestDB(t)
	season := seedFinishedSeason(t)

	early := makeUser(t, "early@test.dev", 0)
	late := makeUser(t, "late@test.dev", 0)
	base := season.StartsAt

	database.DB.Create(&models.LeagueMembership{
		UserID: early.ID, SeasonID: season.ID, Tier: 1, PodID: "tie-pod",
		Language: "es", Points: 500, RawXP: 500, Activities: 5,
		JoinedAt: base, LastPointAt: base.Add(24 * time.Hour),
	})
	database.DB.Create(&models.LeagueMembership{
		UserID: late.ID, SeasonID: season.ID, Tier: 1, PodID: "tie-pod",
		Language: "es", Points: 500, RawXP: 500, Activities: 5,
		JoinedAt: base, LastPointAt: base.Add(6 * 24 * time.Hour),
	})

	SettleDueSeasons()

	var a, b models.LeagueMembership
	database.DB.Where("user_id = ?", early.ID).First(&a)
	database.DB.Where("user_id = ?", late.ID).First(&b)
	if a.FinalRank != 1 || b.FinalRank != 2 {
		t.Errorf("tie broken wrongly: early=%d late=%d", a.FinalRank, b.FinalRank)
	}
}

// Settling twice must not pay twice — the ticker and a page load can collide.
func TestSettlementIsIdempotent(t *testing.T) {
	newTestDB(t)
	season := seedFinishedSeason(t)
	pod := seedPod(t, season.ID, 2, 12, "")

	SettleDueSeasons()
	var first models.User
	database.DB.First(&first, pod[0])

	SettleDueSeasons()
	SettleDueSeasons()

	var again models.User
	database.DB.First(&again, pod[0])
	if again.Gems != first.Gems {
		t.Errorf("gems changed on re-settlement: %d -> %d", first.Gems, again.Gems)
	}
	if again.LeagueTier != first.LeagueTier {
		t.Errorf("tier changed on re-settlement: %d -> %d", first.LeagueTier, again.LeagueTier)
	}
}

// A flagged account keeps its XP and its place but forfeits promotion and the
// chest — botting buys nothing.
func TestFlaggedMembersDoNotPromote(t *testing.T) {
	newTestDB(t)
	season := seedFinishedSeason(t)
	pod := seedPod(t, season.ID, 2, podSize, "")

	database.DB.Model(&models.LeagueMembership{}).
		Where("user_id = ? AND season_id = ?", pod[0], season.ID).
		Updates(map[string]interface{}{"flagged": true, "flag_reason": "test"})

	SettleDueSeasons()

	var m models.LeagueMembership
	database.DB.Where("user_id = ? AND season_id = ?", pod[0], season.ID).First(&m)
	if m.Result != "held" {
		t.Errorf("flagged leader result = %q, want held", m.Result)
	}
	if m.GemsAwarded != 0 && m.GemsAwarded != tiers[2].GroupBonus {
		t.Errorf("flagged leader took a chest of %d gems", m.GemsAwarded)
	}
}

// A small pod shouldn't promote most of itself, nor demote anyone.
func TestSmallPodsScaleTheirSlots(t *testing.T) {
	newTestDB(t)
	season := seedFinishedSeason(t)
	pod := seedPod(t, season.ID, 2, 6, "")

	SettleDueSeasons()

	promoted, demoted := 0, 0
	for _, id := range pod {
		var m models.LeagueMembership
		database.DB.Where("user_id = ? AND season_id = ?", id, season.ID).First(&m)
		switch m.Result {
		case "promoted":
			promoted++
		case "demoted":
			demoted++
		}
	}
	if promoted > 3 {
		t.Errorf("%d of 6 promoted, want at most half", promoted)
	}
	if demoted != 0 {
		t.Errorf("%d demoted from a 6-person pod, want 0", demoted)
	}
}

// Diamond has no tier above it: its top ten enter the tournament bracket.
func TestDiamondTopTenQualifyForTheTournament(t *testing.T) {
	newTestDB(t)
	season := seedFinishedSeason(t)
	pod := seedPod(t, season.ID, diamondTier, podSize, "")

	SettleDueSeasons()

	for i, id := range pod[:10] {
		var m models.LeagueMembership
		database.DB.Where("user_id = ? AND season_id = ?", id, season.ID).First(&m)
		var u models.User
		database.DB.First(&u, id)
		if m.Result != "qualified" || u.TournamentStage != "quarterfinal" {
			t.Errorf("rank %d: %q / stage %q, want qualified / quarterfinal",
				i+1, m.Result, u.TournamentStage)
		}
		if u.LeagueTier != diamondTier {
			t.Errorf("rank %d left Diamond (tier %d)", i+1, u.LeagueTier)
		}
	}
}

// The bracket runs quarterfinal -> semifinal -> final, and the final pays
// trophies to the podium before returning everyone to Diamond.
func TestTournamentBracketAdvancesThenCrowns(t *testing.T) {
	newTestDB(t)

	// Quarterfinal.
	qf := seedFinishedSeason(t)
	pod := seedPod(t, qf.ID, diamondTier, podSize, "quarterfinal")
	SettleDueSeasons()

	var advanced models.User
	database.DB.First(&advanced, pod[0])
	if advanced.TournamentStage != "semifinal" {
		t.Fatalf("winner stage = %q, want semifinal", advanced.TournamentStage)
	}
	var knocked models.User
	database.DB.First(&knocked, pod[podSize-1])
	if knocked.TournamentStage != "" {
		t.Errorf("eliminated player still staged as %q", knocked.TournamentStage)
	}

	// Final.
	fin := seedFinishedSeasonAt(t, time.Now().UTC().AddDate(0, 0, -14))
	finalists := seedPod(t, fin.ID, diamondTier, 30, "final")
	SettleDueSeasons()

	for i, id := range finalists[:3] {
		var u models.User
		database.DB.First(&u, id)
		if u.Trophies != 1 {
			t.Errorf("finalist %d has %d trophies, want 1", i+1, u.Trophies)
		}
		if u.TournamentStage != "" {
			t.Errorf("finalist %d still in stage %q after the final", i+1, u.TournamentStage)
		}
	}
	var fourth models.User
	database.DB.First(&fourth, finalists[3])
	if fourth.Trophies != 0 {
		t.Errorf("fourth place took a trophy")
	}
}

// The hidden rating is what makes tanking pointless: it survives a demotion.
func TestRatingSurvivesDemotion(t *testing.T) {
	newTestDB(t)
	season := seedFinishedSeason(t)
	pod := seedPod(t, season.ID, 3, podSize, "")

	// Give the bottom finisher a strong prior rating — the sandbagger's profile.
	last := pod[podSize-1]
	database.DB.Model(&models.User{}).Where("id = ?", last).
		Updates(map[string]interface{}{"league_mmr": 1200, "integrity": 100})

	SettleDueSeasons()

	var u models.User
	database.DB.First(&u, last)
	if u.LeagueTier != 2 {
		t.Errorf("expected demotion to tier 2, got %d", u.LeagueTier)
	}
	if u.LeagueMMR < 700 {
		t.Errorf("rating collapsed to %d after one tanked week — sandbagging would pay", u.LeagueMMR)
	}
	if u.Integrity >= 100 {
		t.Errorf("integrity untouched (%d) despite a tank", u.Integrity)
	}
}

// --- fixtures ----------------------------------------------------------------

func makeUser(t *testing.T, email string, streak int) *models.User {
	t.Helper()
	u := &models.User{
		Email: email, Name: email, TargetLanguage: "es",
		Streak: streak, Integrity: 100, Hearts: 5,
	}
	if err := database.DB.Create(u).Error; err != nil {
		t.Fatalf("create user: %v", err)
	}
	return u
}

// seedFinishedSeason creates a season that ended yesterday, so SettleDueSeasons
// picks it up.
func seedFinishedSeason(t *testing.T) models.LeagueSeason {
	t.Helper()
	return seedFinishedSeasonAt(t, time.Now().UTC().AddDate(0, 0, -7))
}

func seedFinishedSeasonAt(t *testing.T, within time.Time) models.LeagueSeason {
	t.Helper()
	start, end := seasonBounds(within)
	s := models.LeagueSeason{ID: seasonIDAt(within), StartsAt: start, EndsAt: end}
	if err := database.DB.Create(&s).Error; err != nil {
		t.Fatalf("create season: %v", err)
	}
	return s
}

// seedPod fills one pod with `size` learners whose points descend by rank, and
// returns their IDs in finishing order.
func seedPod(t *testing.T, seasonID string, tier, size int, stage string) []uint {
	t.Helper()
	podID := fmt.Sprintf("%s-pod-%d-%s", seasonID, tier, stage)
	ids := make([]uint, 0, size)

	for i := 0; i < size; i++ {
		u := makeUser(t, fmt.Sprintf("%s-%d-%d@test.dev", seasonID, tier, i), 3)
		u.LeagueTier = tier
		u.TournamentStage = stage
		database.DB.Save(u)

		start, _ := seasonBounds(time.Now().UTC().AddDate(0, 0, -7))
		if err := database.DB.Create(&models.LeagueMembership{
			UserID: u.ID, SeasonID: seasonID, Tier: tier, PodID: podID,
			Language: "es", Stage: stage,
			Points: (size - i) * 100,
			RawXP:  (size - i) * 100,
			// Descending accuracy sum keeps Activities > 0 so the row counts as
			// active rather than a no-show.
			Activities: 10, AccuracySum: 900,
			JoinedAt: start, LastPointAt: start.Add(time.Duration(i) * time.Hour),
		}).Error; err != nil {
			t.Fatalf("create membership: %v", err)
		}
		ids = append(ids, u.ID)
	}
	return ids
}
