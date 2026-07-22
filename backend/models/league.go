package models

import "time"

// The league system is a weekly, pod-based competition. Three tables carry it:
//
//	LeagueSeason     — one row per ISO week, tracks whether settlement has run
//	LeagueMembership — one row per user per season: their pod, points and result
//	LeagueDaily      — per-user, per-day ledger used for diminishing returns and
//	                   bot detection (the membership row alone can't tell how the
//	                   points were spread across the week)
//
// Points are deliberately NOT raw XP. See controllers/league_engine.go for the
// weighting: difficulty x accuracy x consistency x earliness x source, then
// diminishing returns. RawXP is kept alongside purely so the UI can show "you
// earned N XP this week" and so anti-cheat has an unweighted signal.

// LeagueSeason is one weekly cycle, keyed by ISO week ("2026-W30").
type LeagueSeason struct {
	ID        string    `gorm:"primaryKey" json:"id"` // "2026-W30"
	StartsAt  time.Time `json:"startsAt"`
	EndsAt    time.Time `json:"endsAt"`
	Settled   bool      `gorm:"index" json:"settled"`
	SettledAt time.Time `json:"settledAt"`
	CreatedAt time.Time `json:"createdAt"`
}

// LeagueMembership is a user's entry in one season. A row only exists once the
// user has earned at least 1 point that week — that is the "entry rule": no
// activity means no placement (a "league pause") rather than a last-place
// finish.
type LeagueMembership struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	UserID   uint   `gorm:"index:idx_league_user_season,unique" json:"userId"`
	SeasonID string `gorm:"index:idx_league_user_season,unique;index" json:"seasonId"`

	// Placement
	Tier     int    `gorm:"index" json:"tier"`  // 0 Bronze .. 9 Diamond
	PodID    string `gorm:"index" json:"podId"` // 30 learners max
	Language string `json:"language"`           // pods are language-segregated
	SeedMMR  int    `json:"seedMmr"`            // hidden rating at seeding time
	Stage    string `json:"stage"`              // "" | quarterfinal | semifinal | final

	// Score
	Points      int       `gorm:"index" json:"points"` // weighted league points
	RawXP       int       `json:"rawXp"`               // unweighted, for display
	Activities  int       `json:"activities"`          // lessons/sessions completed
	AccuracySum int       `json:"-"`                   // running sum, for the average
	PerfectRuns int       `json:"perfectRuns"`         // 100%-accuracy completions
	JoinedAt    time.Time `json:"joinedAt"`
	LastPointAt time.Time `json:"lastPointAt"` // tie-breaker: earliest to the total

	// Result — written once at settlement.
	Settled      bool      `gorm:"index" json:"settled"`
	FinalRank    int       `json:"finalRank"`
	Result       string    `json:"result"` // promoted | held | demoted | champion | eliminated
	NextTier     int       `json:"nextTier"`
	GemsAwarded  int       `json:"gemsAwarded"`
	GroupGoalHit bool      `json:"groupGoalHit"`
	CeremonySeen bool      `gorm:"index" json:"ceremonySeen"` // has the user watched the animation?
	SettledAt    time.Time `json:"settledAt"`

	// Live-notification state. The league tells the user what's happening to
	// them during the week, not only at settlement, and that needs a memory of
	// what was true last time we looked — a notification should fire on the
	// *transition* (dropped out of the promotion zone) rather than every time
	// they finish a lesson while outside it.
	LastRank     int    `json:"-"`
	LastZone     string `json:"-"` // promote | hold | demote
	JoinNotified bool   `json:"-"`
	EndingNotice bool   `json:"-"` // final-day nudge sent

	// Integrity
	Flagged     bool   `json:"flagged"` // behavioural anomaly detected
	FlagReason  string `json:"flagReason"`
	ReportCount int    `json:"reportCount"`
}

// LeagueDaily is the per-day ledger behind two anti-gaming rules: diminishing
// returns (the first N raw XP of a day are worth full value, the rest taper) and
// behavioural anomaly detection (impossible activity volume or velocity).
type LeagueDaily struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	UserID   uint   `gorm:"index:idx_league_daily,unique" json:"userId"`
	Day      string `gorm:"index:idx_league_daily,unique" json:"day"` // YYYY-MM-DD (UTC)
	SeasonID string `gorm:"index" json:"seasonId"`

	RawXP      int       `json:"rawXp"`
	Points     int       `json:"points"`
	Activities int       `json:"activities"`
	TimedRuns  int       `json:"timedRuns"` // practice/timed challenges — separately capped
	FirstAt    time.Time `json:"firstAt"`
	LastAt     time.Time `json:"lastAt"`
}

// LeagueReport is a one-tap "this looks like cheating" report from one learner
// about another. Reports are advisory: they raise a counter that, combined with
// behavioural signals, flags an account rather than acting on their own.
type LeagueReport struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	SeasonID   string    `gorm:"index:idx_league_report,unique" json:"seasonId"`
	ReporterID uint      `gorm:"index:idx_league_report,unique" json:"reporterId"`
	SubjectID  uint      `gorm:"index:idx_league_report,unique" json:"subjectId"`
	Reason     string    `json:"reason"`
	CreatedAt  time.Time `json:"createdAt"`
}
