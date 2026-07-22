package models

import "time"

// The ideas workspace: a board of proposals, each owning its own discussion
// thread. Deliberately separate from direct messages (models/chat.go) — mixing
// a team's idea discussion into general chat is what turns both into noise.
//
// The shape is: an Idea has votes, stars, tags, a status that moves through a
// workflow, a threaded conversation, an audit trail of every change, and any
// tasks it was converted into.

// Idea statuses, in workflow order. Status is the unit of progress: an idea
// isn't "done" because someone said so in chat, it's done because it moved.
const (
	IdeaDraft       = "draft"
	IdeaUnderReview = "under_review"
	IdeaApproved    = "approved"
	IdeaInProgress  = "in_progress"
	IdeaCompleted   = "completed"
	IdeaArchived    = "archived"
)

// Idea is one proposal on the board.
type Idea struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	OwnerID     uint   `gorm:"index" json:"ownerId"`
	Title       string `gorm:"index" json:"title"`
	Description string `json:"description"`
	Status      string `gorm:"index" json:"status"`

	// Vote tallies are denormalised so the board can sort thousands of ideas
	// without a join per row. IdeaVote remains the source of truth; these are
	// recalculated from it on every vote.
	Upvotes   int `json:"upvotes"`
	Downvotes int `json:"downvotes"`
	Score     int `gorm:"index" json:"score"` // upvotes - downvotes

	MessageCount   int       `json:"messageCount"`
	LastActivityAt time.Time `gorm:"index" json:"lastActivityAt"`

	// Archiving always records why. An archive with no reason is indistinguishable
	// from an idea that was quietly lost, and the second kind poisons a board.
	ArchivedAt    *time.Time `json:"archivedAt"`
	ArchiveReason string     `json:"archiveReason"`

	// Set when this idea was folded into another. The thread stays readable;
	// the board stops showing it as a live competitor for votes.
	MergedIntoID *uint `gorm:"index" json:"mergedIntoId"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// IdeaVote is one person's position on one idea: +1 or -1. Storing downvotes
// rather than only counting support is what makes "controversial" sortable —
// an idea at 30 up / 28 down is a very different signal from one at 2 / 0.
type IdeaVote struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	IdeaID    uint      `gorm:"index:idx_idea_vote,unique" json:"ideaId"`
	UserID    uint      `gorm:"index:idx_idea_vote,unique" json:"userId"`
	Value     int       `json:"value"` // +1 or -1
	CreatedAt time.Time `json:"createdAt"`
}

// IdeaStar is a private bookmark — "come back to this" — with no effect on
// ranking, so saving something for later doesn't read as endorsing it.
type IdeaStar struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	IdeaID    uint      `gorm:"index:idx_idea_star,unique" json:"ideaId"`
	UserID    uint      `gorm:"index:idx_idea_star,unique" json:"userId"`
	CreatedAt time.Time `json:"createdAt"`
}

// IdeaTag is a free-form label. A join table rather than a comma-separated
// column so filtering by tag is an index hit instead of a LIKE scan.
type IdeaTag struct {
	ID     uint   `gorm:"primaryKey" json:"id"`
	IdeaID uint   `gorm:"index:idx_idea_tag,unique" json:"ideaId"`
	Tag    string `gorm:"index:idx_idea_tag,unique;index" json:"tag"`
}

// Message kinds. The composer picks one; the renderer switches on it.
const (
	MsgText  = "text"
	MsgImage = "image"
	MsgVoice = "voice"
	MsgCode  = "code"
)

// IdeaMessage is one message in an idea's thread. ParentID makes threads
// branchable: "what about privacy?" becomes its own sub-thread instead of
// derailing the main line.
type IdeaMessage struct {
	ID       uint  `gorm:"primaryKey" json:"id"`
	IdeaID   uint  `gorm:"index" json:"ideaId"`
	AuthorID uint  `gorm:"index" json:"authorId"`
	ParentID *uint `gorm:"index" json:"parentId"` // nil = top level

	Kind string `json:"kind"`
	Body string `json:"body"`

	// Attachments live in the database for the same reason avatars do: the
	// production filesystem is ephemeral and would drop them on every deploy.
	// Served by GET /api/ideas/attachments/:id.
	Data     []byte `json:"-"`
	Mime     string `json:"-"`
	FileName string `json:"fileName"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	Duration int    `json:"duration"` // voice memo length, seconds

	// Anonymous messages come from a silent-brainstorm session. The author is
	// still recorded (moderation needs it) but never exposed by the API.
	Anonymous bool `json:"anonymous"`

	// Detected language and its English translation — see the same fields on
	// models.Message. Computed once at post time, not per read.
	DetectedLang   string     `json:"detectedLang"`
	TranslatedBody string     `json:"-"`
	TranslatedAt   *time.Time `json:"-"`

	EditedAt  *time.Time `json:"editedAt"`
	DeletedAt *time.Time `gorm:"index" json:"deletedAt"`
	CreatedAt time.Time  `gorm:"index" json:"createdAt"`
}

// IdeaReaction is a one-tap emoji response, so agreement doesn't need a whole
// message. One reaction per emoji per person per message.
type IdeaReaction struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	MessageID uint      `gorm:"index:idx_idea_reaction,unique" json:"messageId"`
	UserID    uint      `gorm:"index:idx_idea_reaction,unique" json:"userId"`
	Emoji     string    `gorm:"index:idx_idea_reaction,unique" json:"emoji"`
	CreatedAt time.Time `json:"createdAt"`
}

// IdeaEvent is the version history: who changed what, when. Every mutation of
// an idea writes one, which is what lets the details panel answer "why is this
// approved?" months later.
type IdeaEvent struct {
	ID      uint   `gorm:"primaryKey" json:"id"`
	IdeaID  uint   `gorm:"index" json:"ideaId"`
	ActorID uint   `json:"actorId"`
	Kind    string `json:"kind"`  // created | status | edited | tagged | merged | archived | task | vote_threshold
	Field   string `json:"field"` // which attribute moved, when relevant
	From    string `json:"from"`
	To      string `json:"to"`
	Note    string `json:"note"`

	CreatedAt time.Time `gorm:"index" json:"createdAt"`
}

// IdeaTask is an idea that graduated into work. Deliberately a thin record
// rather than a task system: it carries the link back to the idea and a status,
// and expects a real tracker to own the rest.
type IdeaTask struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	IdeaID      uint       `gorm:"index" json:"ideaId"`
	Title       string     `json:"title"`
	Status      string     `json:"status"` // todo | doing | done
	Sprint      string     `json:"sprint"`
	AssigneeID  *uint      `json:"assigneeId"`
	CreatedByID uint       `json:"createdById"`
	CreatedAt   time.Time  `json:"createdAt"`
	CompletedAt *time.Time `json:"completedAt"`
}

// BrainstormSession is the "silent brainstorm": for a fixed window everyone
// contributes anonymously, so the loudest voice and the most senior one stop
// setting the direction before quieter ideas are on the board. When it ends,
// authorship stays hidden but voting opens.
type BrainstormSession struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	IdeaID    *uint     `gorm:"index" json:"ideaId"` // nil = board-wide
	StartedBy uint      `json:"startedBy"`
	Topic     string    `json:"topic"`
	StartsAt  time.Time `json:"startsAt"`
	EndsAt    time.Time `gorm:"index" json:"endsAt"`
	CreatedAt time.Time `json:"createdAt"`
}
