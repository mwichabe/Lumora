package controllers

import (
	"fmt"
	"log"
	"math"
	"sort"
	"time"

	"lumora/backend/database"
	"lumora/backend/models"
)

// ============================================================================
// The Lumora league engine
// ============================================================================
//
// A ten-tier weekly competition in pods of 30. The shape (tiers, promotion
// slots, weekly reset, a three-week Diamond bracket) will look familiar to
// anyone who has climbed a Duolingo leaderboard. What sits underneath is
// deliberately different, because the familiar shape has well-known exploits:
// grinding trivial content, sandbagging into a weak pod, tanking a week to farm
// an easier one, and cramming everything into the final hours.
//
// The countermeasures, and where each lives:
//
//	XP farming easy content   -> difficultyWeight(): league points scale with
//	                             lesson difficulty, so Unit 1 is worth half of
//	                             advanced content.
//	Low-effort accuracy       -> accuracyWeight(): 50% accuracy is worth 0.8x,
//	                             and a flawless run earns a deep-practice bonus.
//	Timed-challenge spam      -> sourceWeight() + dailyTimedCap: practice runs
//	                             count for less and stop counting after a cap.
//	All-night grinding/botting-> taperPoints() diminishing returns, plus
//	                             detectAnomaly() behavioural flags.
//	Join-late sandbagging     -> pods are seeded by LeagueMMR, not arrival time.
//	Intentional demotion      -> MMR persists across tiers, so a tanked week
//	                             lands you against the same calibre next week,
//	                             and costs Integrity.
//	Cross-language stacking   -> pods are segregated by target language.
//	Sunday sprint             -> earlinessWeight(): Monday XP is worth 1.20x,
//	                             Sunday 0.85x. Consistency beats cramming.
//	Zero-sum burnout          -> every pod has a collaborative group goal, and
//	                             LeagueCasual opts out of ranking entirely.

// --- tiers -------------------------------------------------------------------

type tierDef struct {
	Name string
	Tint string

	// PromoteTop learners move up; the bottom DemoteBottom move down. The
	// promotion band narrows as you climb, so higher tiers are genuinely harder
	// to hold.
	PromoteTop   int
	DemoteBottom int

	// GoldGems is the first-place chest; silver and bronze are derived from it.
	GoldGems int

	// GoalPerLearner x pod size is the collaborative target. Hitting it pays
	// GroupBonus gems to everyone who scored, win or lose — 30% of a good week's
	// reward comes from the pod pulling together rather than from beating it.
	GoalPerLearner int
	GroupBonus     int
}

// Bronze has no demotion: there is nothing below it. Diamond has no promotion —
// its top ten qualify for the tournament bracket instead.
var tiers = []tierDef{
	{"Bronze", "#CD7F32", 15, 0, 60, 120, 15},
	{"Silver", "#9CA3AF", 12, 5, 80, 150, 20},
	{"Gold", "#F5A623", 10, 5, 100, 180, 25},
	{"Sapphire", "#00C2A8", 8, 5, 130, 220, 30},
	{"Ruby", "#FF5C5C", 7, 5, 160, 260, 35},
	{"Emerald", "#10B981", 7, 5, 200, 300, 40},
	{"Amethyst", "#6C3FC5", 7, 5, 240, 350, 50},
	{"Pearl", "#E8D8F0", 7, 5, 290, 400, 60},
	{"Obsidian", "#1A1A2E", 5, 5, 350, 460, 70},
	{"Diamond", "#17A3DD", 10, 5, 450, 520, 90},
}

const (
	diamondTier = 9  // index of Diamond in tiers
	podSize     = 30 // learners per pod

	// Integrity below this loses the fair-play badge and halves chest rewards.
	fairPlayFloor = 70

	// Diminishing returns: raw XP in a single day past these thresholds is worth
	// less. Genuine study rarely passes the first band; farming always does.
	fullValueXP = 500  // 1.00x
	halfValueXP = 1000 // 0.50x, then 0.25x beyond

	// Practice/timed runs that still earn league points each day. Beyond this
	// they award XP as normal but stop moving the leaderboard.
	dailyTimedCap = 5

	// Behavioural anomaly thresholds — see detectAnomaly.
	anomalyDailyXP      = 2500
	anomalyDailyActs    = 80
	anomalyBurstActs    = 30
	anomalyBurstMinutes = 20
)

func tierName(i int) string {
	if i < 0 || i >= len(tiers) {
		return tiers[0].Name
	}
	return tiers[i].Name
}

func clampTier(i int) int {
	if i < 0 {
		return 0
	}
	if i > diamondTier {
		return diamondTier
	}
	return i
}

// tierIndexByName maps a stored league name back to its index. It exists for the
// migration path: accounts created before the ten-tier system have a League
// string but no LeagueTier.
func tierIndexByName(name string) int {
	for i, t := range tiers {
		if t.Name == name {
			return i
		}
	}
	return 0
}

// --- seasons -----------------------------------------------------------------

// Seasons are ISO weeks in UTC: they start Monday 00:00 and run seven days. UTC
// rather than local time so that every learner in a pod races the same clock —
// a timezone-relative reset would hand an advantage to whoever resets last.
func seasonIDAt(t time.Time) string {
	y, w := t.UTC().ISOWeek()
	return fmt.Sprintf("%d-W%02d", y, w)
}

func currentSeasonID() string { return seasonIDAt(time.Now()) }

// seasonBounds returns the [start, end) of the ISO week containing t.
func seasonBounds(t time.Time) (time.Time, time.Time) {
	u := t.UTC()
	// Go's Weekday is Sunday=0; shift so Monday=0.
	offset := (int(u.Weekday()) + 6) % 7
	start := time.Date(u.Year(), u.Month(), u.Day(), 0, 0, 0, 0, time.UTC).
		AddDate(0, 0, -offset)
	return start, start.AddDate(0, 0, 7)
}

// ensureSeason creates the season row on first use and returns it.
func ensureSeason(id string, at time.Time) models.LeagueSeason {
	var s models.LeagueSeason
	if database.DB.Where("id = ?", id).First(&s).Error == nil {
		return s
	}
	start, end := seasonBounds(at)
	s = models.LeagueSeason{ID: id, StartsAt: start, EndsAt: end}
	database.DB.Create(&s)
	return s
}

// --- scoring -----------------------------------------------------------------

// difficultyWeight turns a skill's XP gate into a difficulty multiplier. The
// gate is how far into the course the content sits, which is the best proxy the
// content model gives us. Beginner content is worth 1.0x and the hardest 2.0x,
// so pushing your actual boundary outscores replaying Unit 1.
func difficultyWeight(requiredXP int) float64 {
	switch {
	case requiredXP >= 2000:
		return 2.0
	case requiredXP >= 1200:
		return 1.8
	case requiredXP >= 600:
		return 1.6
	case requiredXP >= 300:
		return 1.4
	case requiredXP >= 100:
		return 1.2
	default:
		return 1.0
	}
}

// difficultyForLesson looks up the skill behind a lesson and weights it.
func difficultyForLesson(lesson models.Lesson) float64 {
	var skill models.Skill
	if database.DB.First(&skill, lesson.SkillID).Error != nil {
		return 1.0
	}
	return difficultyWeight(skill.RequiredXP)
}

// accuracyWeight rewards getting it right. A 50% run is worth 0.8x, a 100% run
// 1.0x — and flawless completions earn a further deep-practice bonus, so
// carefully clearing one hard lesson beats sloppily clearing three easy ones.
func accuracyWeight(accuracy int) float64 {
	if accuracy < 0 {
		accuracy = 0
	}
	if accuracy > 100 {
		accuracy = 100
	}
	w := 0.6 + 0.4*(float64(accuracy)/100)
	if accuracy == 100 {
		w *= 1.15 // deep-practice bonus
	}
	return w
}

// consistencyWeight pays a small, capped bonus for showing up daily. Capped at
// 1.25x so a long streak is an edge, not an insurmountable head start for
// newcomers.
func consistencyWeight(streak int) float64 {
	if streak > 25 {
		streak = 25
	}
	if streak < 0 {
		streak = 0
	}
	return 1 + float64(streak)*0.01
}

// earlinessWeight is the answer to the Sunday sprint. XP earned on Monday is
// worth 1.20x and decays to 0.85x by Sunday, so a learner who studies a little
// every day beats one who hoards boosts for the final hours.
func earlinessWeight(at time.Time, seasonStart time.Time) float64 {
	day := int(at.UTC().Sub(seasonStart).Hours() / 24)
	if day < 0 {
		day = 0
	}
	if day > 6 {
		day = 6
	}
	return 1.20 - 0.0583*float64(day)
}

// Source weights: structured study moves the leaderboard, quick drills mostly
// don't. Practice and timed challenges still pay full XP — they just aren't a
// route to the top of a pod.
func sourceWeight(source string) float64 {
	switch source {
	case "lesson":
		return 1.0
	case "listening", "reading":
		return 1.0
	case "practice":
		return 0.6
	case "timed":
		return 0.5
	default:
		return 0.8
	}
}

// taperPoints applies diminishing returns across the day's raw XP. The first
// 500 XP of a day counts in full, the next 500 at half, everything beyond at a
// quarter. A human session never notices this; an overnight bot earns a
// fraction of what the raw numbers suggest.
func taperPoints(alreadyToday, gain int) float64 {
	remaining := float64(gain)
	cursor := float64(alreadyToday)
	var out float64
	for remaining > 0 {
		var bandEnd, rate float64
		switch {
		case cursor < fullValueXP:
			bandEnd, rate = fullValueXP, 1.0
		case cursor < halfValueXP:
			bandEnd, rate = halfValueXP, 0.5
		default:
			bandEnd, rate = math.Inf(1), 0.25
		}
		chunk := math.Min(remaining, bandEnd-cursor)
		out += chunk * rate
		cursor += chunk
		remaining -= chunk
	}
	return out
}

// --- earning -----------------------------------------------------------------

// LeagueAward describes one scoring event. Callers fill in what they know; the
// engine supplies the rest.
type LeagueAward struct {
	Source     string  // lesson | listening | reading | practice | timed
	RawXP      int     // XP the activity paid the user
	Accuracy   int     // 0..100 (pass 100 when the activity has no accuracy)
	Difficulty float64 // 1.0..2.0; 0 means "use 1.0"
}

// AwardLeaguePoints converts one completed activity into weighted league points
// and records them against the user's current-season membership. It is the only
// way points enter the system.
//
// Safe to call for anyone: casual-mode learners are skipped, and a membership
// row is created on demand — which is exactly the entry rule. Earn at least one
// point and you're placed; stay idle all week and you're simply not ranked,
// rather than being parked at the bottom of a pod you never joined.
func AwardLeaguePoints(user *models.User, a LeagueAward) int {
	if user == nil || a.RawXP <= 0 || user.LeagueCasual {
		return 0
	}
	now := time.Now().UTC()
	seasonID := seasonIDAt(now)
	season := ensureSeason(seasonID, now)

	daily := loadDaily(user.ID, seasonID, now)

	// Timed/practice drills stop counting toward the league past the daily cap.
	isTimed := a.Source == "practice" || a.Source == "timed"
	if isTimed && daily.TimedRuns >= dailyTimedCap {
		bumpDaily(&daily, a.RawXP, 0, isTimed, now)
		return 0
	}

	if a.Difficulty <= 0 {
		a.Difficulty = 1.0
	}
	weight := a.Difficulty *
		accuracyWeight(a.Accuracy) *
		consistencyWeight(user.Streak) *
		earlinessWeight(now, season.StartsAt) *
		sourceWeight(a.Source)

	// Diminishing returns are applied to the raw XP, then the quality weights
	// scale what survives.
	points := int(math.Round(taperPoints(daily.RawXP, a.RawXP) * weight))
	if points < 0 {
		points = 0
	}

	m := ensureMembership(user, seasonID)
	if m == nil {
		return 0
	}
	m.Points += points
	m.RawXP += a.RawXP
	m.Activities++
	m.AccuracySum += a.Accuracy
	if a.Accuracy >= 100 {
		m.PerfectRuns++
	}
	// The tie-breaker: equal totals are ranked by who reached them first, so
	// this only advances when points actually move.
	if points > 0 {
		m.LastPointAt = now
	}
	database.DB.Save(m)

	bumpDaily(&daily, a.RawXP, points, isTimed, now)
	detectAnomaly(user, m, daily)

	// Turn the change in standing into in-week notifications — entering the
	// promotion zone, falling into the drop zone, being overtaken. See
	// league_notify.go.
	evaluateLeagueMoments(user, m, points)

	return points
}

func loadDaily(userID uint, seasonID string, now time.Time) models.LeagueDaily {
	day := now.Format("2006-01-02")
	var d models.LeagueDaily
	if database.DB.Where("user_id = ? AND day = ?", userID, day).First(&d).Error != nil {
		d = models.LeagueDaily{UserID: userID, Day: day, SeasonID: seasonID, FirstAt: now}
	}
	return d
}

func bumpDaily(d *models.LeagueDaily, rawXP, points int, timed bool, now time.Time) {
	d.RawXP += rawXP
	d.Points += points
	d.Activities++
	if timed {
		d.TimedRuns++
	}
	if d.FirstAt.IsZero() {
		d.FirstAt = now
	}
	d.LastAt = now
	database.DB.Save(d)
}

// detectAnomaly flags accounts whose activity pattern isn't physically
// plausible: an impossible daily volume, an impossible number of completions, or
// a burst of lessons far faster than a person can read them. Flagged members
// keep competing and keep their XP — they're excluded from promotion and chest
// rewards for that season, which makes botting pointless without punishing a
// false positive too harshly.
func detectAnomaly(user *models.User, m *models.LeagueMembership, d models.LeagueDaily) {
	reason := ""
	switch {
	case d.RawXP > anomalyDailyXP:
		reason = "implausible daily XP volume"
	case d.Activities > anomalyDailyActs:
		reason = "implausible daily completion count"
	case d.Activities >= anomalyBurstActs &&
		d.LastAt.Sub(d.FirstAt) < anomalyBurstMinutes*time.Minute:
		reason = "completion velocity faster than humanly possible"
	}
	if reason == "" || m.Flagged {
		return
	}
	m.Flagged = true
	m.FlagReason = reason
	database.DB.Save(m)

	user.Integrity = clamp(user.Integrity-25, 0, 100)
	database.DB.Model(user).Update("integrity", user.Integrity)
	log.Printf("[league] flagged user %d in %s: %s", user.ID, m.SeasonID, reason)
}

// --- placement ---------------------------------------------------------------

// ensureMembership returns the user's row for the season, creating and seeding
// it on first activity.
func ensureMembership(user *models.User, seasonID string) *models.LeagueMembership {
	var m models.LeagueMembership
	err := database.DB.Where("user_id = ? AND season_id = ?", user.ID, seasonID).
		First(&m).Error
	if err == nil {
		return &m
	}

	normaliseLeagueState(user)

	lang := user.TargetLanguage
	if lang == "" {
		lang = "es"
	}
	stage := user.TournamentStage
	tier := user.LeagueTier
	if stage != "" {
		tier = diamondTier // the bracket is played out from Diamond
	}

	m = models.LeagueMembership{
		UserID:      user.ID,
		SeasonID:    seasonID,
		Tier:        tier,
		Language:    lang,
		Stage:       stage,
		SeedMMR:     user.LeagueMMR,
		JoinedAt:    time.Now().UTC(),
		LastPointAt: time.Now().UTC(),
	}
	m.PodID = assignPod(seasonID, tier, lang, stage, user.LeagueMMR)
	if err := database.DB.Create(&m).Error; err != nil {
		// Lost a race with a concurrent request — re-read the winner.
		if database.DB.Where("user_id = ? AND season_id = ?", user.ID, seasonID).
			First(&m).Error != nil {
			return nil
		}
	}
	return &m
}

// normaliseLeagueState backfills league fields for accounts that predate the
// ten-tier system (or that have never competed), so seeding always has sane
// inputs. Old rows carry a League name but no tier, no MMR and no integrity.
func normaliseLeagueState(user *models.User) {
	dirty := false
	if user.LeagueTier == 0 && user.League != "" && user.League != tiers[0].Name {
		user.LeagueTier = tierIndexByName(user.League)
		dirty = true
	}
	user.LeagueTier = clampTier(user.LeagueTier)
	if user.Integrity == 0 {
		user.Integrity = 100
		dirty = true
	}
	if user.LeagueMMR == 0 {
		// Cold start: estimate from lifetime XP so a long-standing account isn't
		// dropped into a beginners' pod on its first competitive week.
		user.LeagueMMR = clamp(user.XP/4, 0, 4000)
		dirty = true
	}
	if user.LeagueBest < user.LeagueTier {
		user.LeagueBest = user.LeagueTier
		dirty = true
	}
	if name := tierName(user.LeagueTier); user.League != name {
		user.League = name
		dirty = true
	}
	if dirty {
		database.DB.Save(user)
	}
}

// assignPod places a learner in a pod of at most 30. Candidates are the open
// pods for this season/tier/language; the one whose average hidden rating is
// closest to the learner's wins, so pods are competitively even.
//
// This is the fix for join-late sandbagging. Waiting until Friday no longer
// drops you into a pod of stragglers — it drops you into a pod of people with
// your rating, who happen to also have joined late.
func assignPod(seasonID string, tier int, lang, stage string, mmr int) string {
	kind := "L"
	if stage != "" {
		kind = "T" // tournament bracket pods are kept separate
	}

	var members []models.LeagueMembership
	database.DB.Where("season_id = ? AND tier = ? AND language = ? AND stage = ?",
		seasonID, tier, lang, stage).Find(&members)

	type podStat struct {
		count  int
		mmrSum int
	}
	stats := map[string]*podStat{}
	for _, m := range members {
		s := stats[m.PodID]
		if s == nil {
			s = &podStat{}
			stats[m.PodID] = s
		}
		s.count++
		s.mmrSum += m.SeedMMR
	}

	best, bestDelta := "", math.MaxFloat64
	for id, s := range stats {
		if s.count >= podSize {
			continue
		}
		delta := math.Abs(float64(s.mmrSum)/float64(s.count) - float64(mmr))
		if delta < bestDelta {
			best, bestDelta = id, delta
		}
	}
	if best != "" {
		return best
	}
	return fmt.Sprintf("%s-%s%d-%s-%03d", seasonID, kind, tier, lang, len(stats)+1)
}

// --- settlement --------------------------------------------------------------

// SettleDueSeasons closes out every season whose week has ended. It's called
// from the league endpoints and from a background ticker, and it is idempotent:
// the season row's Settled flag means a race between the two can't double-award.
func SettleDueSeasons() {
	var due []models.LeagueSeason
	database.DB.Where("settled = ? AND ends_at <= ?", false, time.Now().UTC()).
		Order("ends_at asc").Find(&due)
	for _, s := range due {
		settleSeason(s)
	}
}

func settleSeason(season models.LeagueSeason) {
	// Claim the season first. An UPDATE ... WHERE settled = false is atomic, so
	// only one caller proceeds even if the ticker and a request collide.
	res := database.DB.Model(&models.LeagueSeason{}).
		Where("id = ? AND settled = ?", season.ID, false).
		Updates(map[string]interface{}{"settled": true, "settled_at": time.Now().UTC()})
	if res.RowsAffected == 0 {
		return
	}
	log.Printf("[league] settling season %s", season.ID)

	var members []models.LeagueMembership
	database.DB.Where("season_id = ?", season.ID).Find(&members)

	pods := map[string][]models.LeagueMembership{}
	for _, m := range members {
		pods[m.PodID] = append(pods[m.PodID], m)
	}
	for _, pod := range pods {
		settlePod(pod)
	}
	log.Printf("[league] season %s settled: %d pods, %d members",
		season.ID, len(pods), len(members))
}

// settlePod ranks one pod and writes every member's result.
func settlePod(pod []models.LeagueMembership) {
	if len(pod) == 0 {
		return
	}
	// Rank by points, then by who reached the total first — so a tie is broken
	// in favour of the learner who got there without waiting for the deadline.
	sort.SliceStable(pod, func(i, j int) bool {
		if pod[i].Points != pod[j].Points {
			return pod[i].Points > pod[j].Points
		}
		return pod[i].LastPointAt.Before(pod[j].LastPointAt)
	})

	tier := clampTier(pod[0].Tier)
	def := tiers[tier]
	stage := pod[0].Stage

	// Collaborative goal: 30% of a good week comes from the pod pulling
	// together, which takes some heat out of a purely zero-sum race.
	goal := def.GoalPerLearner * podSize
	total := 0
	for _, m := range pod {
		total += m.Points
	}
	goalHit := total >= goal

	promote, demote := slotsFor(def, len(pod))

	for i := range pod {
		m := &pod[i]
		rank := i + 1

		var user models.User
		if database.DB.First(&user, m.UserID).Error != nil {
			continue
		}

		result, nextTier := "held", tier
		switch {
		case stage != "":
			result, nextTier = settleTournamentPlace(&user, m, rank, stage)
		case m.Points <= 0:
			// Placed but never scored: hold, don't demote. Nobody is punished for
			// a row that only exists because of a stray point.
			result = "held"
		case m.Flagged:
			// Flagged accounts hold regardless of rank: no promotion, no chest.
			result = "held"
		case rank <= promote && tier < diamondTier:
			result, nextTier = "promoted", tier+1
		case rank <= promote && tier == diamondTier:
			result = "qualified" // into the Diamond tournament
			user.TournamentStage = "quarterfinal"
		case demote > 0 && rank > len(pod)-demote && tier > 0:
			result, nextTier = "demoted", tier-1
		}

		gems := chestFor(def, rank, m, user)
		if goalHit && m.Points > 0 {
			gems += def.GroupBonus
			m.GroupGoalHit = true
		}

		m.FinalRank = rank
		m.Result = result
		m.NextTier = nextTier
		m.GemsAwarded = gems
		m.Settled = true
		m.SettledAt = time.Now().UTC()
		database.DB.Save(m)

		applyResultToUser(&user, m, nextTier, gems, len(pod))
		notifyLeagueResult(user, *m, def, len(pod))
	}
}

// slotsFor scales the promotion and demotion bands to the pod's actual size. A
// half-empty pod shouldn't promote two thirds of its members, and it shouldn't
// demote anyone unless there are enough people to make the bottom meaningful.
func slotsFor(def tierDef, size int) (promote, demote int) {
	promote = def.PromoteTop
	if promote > size/2 {
		promote = size / 2
	}
	if promote < 1 {
		promote = 1
	}
	demote = def.DemoteBottom
	if size < promote+demote+3 {
		demote = 0
	}
	return promote, demote
}

// chestFor pays the podium. Gold, silver and bronze chests only — matching a
// pod's top three, with everyone else earning from the group goal instead.
// Flagged accounts and accounts below the fair-play floor forfeit or halve it.
func chestFor(def tierDef, rank int, m *models.LeagueMembership, user models.User) int {
	if m.Flagged || m.Points <= 0 {
		return 0
	}
	var gems int
	switch rank {
	case 1:
		gems = def.GoldGems
	case 2:
		gems = def.GoldGems * 6 / 10
	case 3:
		gems = def.GoldGems * 35 / 100
	default:
		return 0
	}
	if user.Integrity < fairPlayFloor {
		gems /= 2
	}
	return gems
}

// settleTournamentPlace resolves a Diamond bracket pod. Ten of thirty advance
// through the quarterfinal and semifinal; the final pays trophies to the top
// three. Everyone drops back into Diamond afterwards either way.
func settleTournamentPlace(user *models.User, m *models.LeagueMembership, rank int, stage string) (string, int) {
	advance := rank <= 10 && !m.Flagged && m.Points > 0
	switch stage {
	case "quarterfinal":
		if advance {
			user.TournamentStage = "semifinal"
			return "advanced", diamondTier
		}
	case "semifinal":
		if advance {
			user.TournamentStage = "final"
			return "advanced", diamondTier
		}
	case "final":
		user.TournamentStage = ""
		if rank <= 3 && !m.Flagged && m.Points > 0 {
			user.Trophies++
			return "champion", diamondTier
		}
		return "eliminated", diamondTier
	}
	user.TournamentStage = ""
	return "eliminated", diamondTier
}

// applyResultToUser writes the season back onto the account: the new tier, the
// gems, and the hidden rating that will seed next week's pod.
func applyResultToUser(user *models.User, m *models.LeagueMembership, nextTier, gems, podSizeActual int) {
	user.LeagueTier = clampTier(nextTier)
	user.League = tierName(user.LeagueTier)
	if user.LeagueBest < user.LeagueTier {
		user.LeagueBest = user.LeagueTier
	}
	user.Gems += gems

	// The hidden rating is a rolling blend, so one strong or weak week nudges it
	// rather than redefining it. This is what makes intentional demotion a dead
	// end: tanking costs you the tier but not the rating, and the rating is what
	// picks your next opponents.
	user.LeagueMMR = int(math.Round(0.7*float64(user.LeagueMMR) + 0.3*float64(m.Points)))

	switch {
	case detectSandbagging(*user, *m, podSizeActual):
		user.Integrity = clamp(user.Integrity-15, 0, 100)
	case m.Flagged:
		// Already docked at detection time; no further decay here.
	case m.Points > 0:
		// A clean, active week repairs integrity slowly.
		user.Integrity = clamp(user.Integrity+5, 0, 100)
	}
	database.DB.Save(user)
}

// detectSandbagging spots a deliberate tank: a highly-rated account that scored
// far below what its rating predicts and landed in the demotion zone. That's the
// signature of someone farming an easier pod next week.
func detectSandbagging(user models.User, m models.LeagueMembership, size int) bool {
	if m.Result != "demoted" || m.Points <= 0 {
		return false
	}
	return user.LeagueMMR > 400 && m.Points < user.LeagueMMR/4 && size >= 10
}

// --- notifications -----------------------------------------------------------

func notifyLeagueResult(user models.User, m models.LeagueMembership, def tierDef, size int) {
	key := "league_" + m.SeasonID
	title, body, emoji, tint := "", "", "🏅", def.Tint

	place := fmt.Sprintf("#%d of %d", m.FinalRank, size)
	switch m.Result {
	case "promoted":
		emoji = "🚀"
		title = "Promoted to " + tierName(m.NextTier) + " League!"
		body = fmt.Sprintf("You finished %s in %s League with %d points. Next week you race in %s — open the League tab to collect your reward.",
			place, def.Name, m.Points, tierName(m.NextTier))
	case "demoted":
		emoji = "🌧️"
		title = "Moved down to " + tierName(m.NextTier) + " League"
		body = fmt.Sprintf("You finished %s in %s League this week. A fresh pod is waiting — one lesson a day is usually enough to climb straight back.",
			place, def.Name)
	case "qualified":
		emoji = "💠"
		title = "You qualified for the Diamond Tournament!"
		body = fmt.Sprintf("Top ten in Diamond with %d points. The three-week bracket starts Monday: quarterfinal, semifinal, final.", m.Points)
	case "advanced":
		emoji = "⚔️"
		title = "Through to the " + tierName(diamondTier) + " " + nextStageLabel(user.TournamentStage)
		body = fmt.Sprintf("You finished %s and advanced. One more round to survive.", place)
	case "champion":
		emoji = "🏆"
		title = "Diamond Tournament champion!"
		body = fmt.Sprintf("You finished %s in the final. The trophy is yours — it's on your profile for good.", place)
	case "eliminated":
		emoji = "🛡️"
		title = "Tournament run over"
		body = fmt.Sprintf("You finished %s. Back to Diamond League — and another shot at the bracket.", place)
	default:
		emoji = "🏅"
		title = "You held your place in " + def.Name + " League"
		body = fmt.Sprintf("You finished %s with %d points. A new week has started — the race is open again.",
			place, m.Points)
	}
	if m.GemsAwarded > 0 {
		body += fmt.Sprintf(" You earned %d gems.", m.GemsAwarded)
	}
	if m.Flagged {
		emoji = "⚠️"
		title = "Your league result was withheld"
		body = "This week's activity was flagged as automated (" + m.FlagReason +
			"), so promotion and rewards were held back. Your XP and progress are untouched."
	}

	// Deep-linked to the league screen, where the result ceremony is waiting.
	leagueNote(user.ID, key, emoji, tint, title, body, 0)
}

func nextStageLabel(stage string) string {
	switch stage {
	case "semifinal":
		return "semifinal"
	case "final":
		return "final"
	}
	return "next round"
}

// StartLeagueScheduler runs settlement in the background. Seasons are also
// settled lazily whenever anyone opens the league, which is what actually
// carries a host that sleeps when idle — the ticker just means an always-on
// instance closes the week promptly instead of waiting for the first visitor.
func StartLeagueScheduler() {
	go func() {
		for {
			func() {
				defer func() {
					if r := recover(); r != nil {
						log.Printf("[league] settlement panic: %v", r)
					}
				}()
				SettleDueSeasons()
				DeliverFinalDayNudges()
			}()
			time.Sleep(10 * time.Minute)
		}
	}()
}
