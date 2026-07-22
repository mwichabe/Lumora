package models

import "time"

// User is the core account model. It holds learning progress, gamification
// state (XP, gems, streaks) and the gameplay metadata Lumora needs.
type User struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	Email        string `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash string `gorm:"not null" json:"-"`
	Name         string `json:"name"`
	AvatarColor  string `json:"avatarColor"` // hex used for the placeholder avatar ring
	AvatarURL    string `json:"avatarUrl"`   // points at /api/avatars/:id when a photo is set

	// The profile photo lives in the database rather than on disk: hosts with an
	// ephemeral filesystem (Render's free tier) would otherwise drop it on every
	// deploy. Uploads are downscaled before being stored — see UploadAvatar.
	// (GORM maps []byte per driver: bytea on Postgres, BLOB on SQLite.)
	AvatarData []byte `json:"-"`
	AvatarMime string `json:"-"`

	// Learning setup (chosen during onboarding)
	TargetLanguage string `json:"targetLanguage"` // e.g. "es"
	NativeLanguage string `json:"nativeLanguage"` // e.g. "en"
	CEFRLevel      string `json:"cefrLevel"`      // A1..C2
	LevelName      string `json:"levelName"`      // Spark, Glow, Flame...
	DailyGoalXP    int    `json:"dailyGoalXp"`    // 10 / 20 / 30 / 50

	// Gamification state
	XP           int    `json:"xp"`
	XPToday      int    `json:"xpToday"`
	Gems         int    `json:"gems"`
	Hearts       int    `json:"hearts"`
	Streak       int    `json:"streak"`
	FluencyScore int    `json:"fluencyScore"` // 0..1000
	League       string `json:"league"`       // display name of LeagueTier

	// --- League state (see models/league.go and controllers/league_engine.go) ---
	//
	// LeagueTier is the persistent one. Membership rows are per-season and are
	// created lazily on the week's first activity, so the user's standing between
	// seasons lives here: settlement writes the promoted/demoted tier back and the
	// next season seeds from it.
	LeagueTier int `json:"leagueTier"` // 0 Bronze .. 9 Diamond
	LeagueBest int `json:"leagueBest"` // highest tier ever reached (a keepsake)

	// LeagueMMR is a hidden skill rating that persists across leagues. Pods are
	// seeded from it, which is what makes sandbagging pointless: dropping a tier
	// on purpose doesn't drop the rating, so you land against the same calibre of
	// opponent anyway.
	LeagueMMR int `json:"-"`

	// Integrity is 0..100 and decays on detected sandbagging or bot-like
	// behaviour. Below fairPlayFloor the "fair play" badge disappears and chest
	// rewards are halved. It recovers slowly with clean weeks.
	Integrity int `json:"integrity"`

	// LeagueCasual opts the user out of competition entirely: they keep every
	// feature and all XP, they're just never placed in a pod. An escape hatch for
	// people who find the weekly demotion threat stressful.
	LeagueCasual bool `json:"leagueCasual"`

	// TournamentStage is set when a Diamond finisher qualifies for the next
	// three-week bracket: quarterfinal -> semifinal -> final, then back to "".
	TournamentStage string `json:"tournamentStage"`
	Trophies        int    `json:"trophies"` // finals podium finishes

	LastActiveDate string `json:"lastActiveDate"` // YYYY-MM-DD, drives streak logic

	// Monetization entitlement — set once the exam+certificate is paid for.
	ExamUnlocked bool `json:"examUnlocked"`

	// Hearts regenerate over time; this anchors the regen clock (the moment the
	// current partial refill began). Zero when hearts are full.
	HeartsUpdatedAt time.Time `json:"-"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
