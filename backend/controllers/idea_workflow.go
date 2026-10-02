package controllers

import (
	"fmt"

	"lumora/backend/database"
	"lumora/backend/models"
)

// The idea lifecycle. Status only moves along these edges, and each edge says
// who may take it — so "how does a draft get approved?" has one answer, the
// same on web and mobile, enforced here rather than by whichever buttons a
// client happens to show.
//
//	draft ──submit (owner, or auto at reviewThreshold votes)──▶ under_review
//	under_review ──approve (anyone but the owner, net support > 0,
//	                        or auto at approveThreshold votes)──▶ approved
//	under_review ──withdraw (owner) / send back (others)──▶ draft
//	approved ──start work / convert to task──▶ in_progress
//	approved ──reopen review──▶ under_review
//	in_progress ──complete (no open tasks; automatic when the last task is done)──▶ completed
//	in_progress ──pause──▶ approved
//	completed ──reopen──▶ in_progress
//
// Archiving leaves the ladder from any step through its own endpoint (it needs
// a reason) and Restore returns the idea to the step it left from.

const (
	// Net support at which a draft is put into review without anyone asking.
	reviewThreshold = 5
	// Net support at which an idea under review is approved by the community.
	approveThreshold = 10
)

type ideaTransition struct {
	To      string `json:"to"`
	Label   string `json:"label"`
	Hint    string `json:"hint"`
	Allowed bool   `json:"allowed"`
	Reason  string `json:"reason,omitempty"` // why it's not allowed right now
	Primary bool   `json:"primary"`          // the natural next step
}

func openTaskCount(ideaID uint) int64 {
	var n int64
	database.DB.Model(&models.IdeaTask{}).
		Where("idea_id = ? AND status != ?", ideaID, "done").Count(&n)
	return n
}

// transitionsFor lists the moves available from the idea's current status, as
// seen by this viewer. Moves the viewer can't make are still listed (with the
// reason) so the UI can explain what has to happen next.
func transitionsFor(idea models.Idea, viewerID uint) []ideaTransition {
	if idea.ArchivedAt != nil || idea.MergedIntoID != nil {
		return []ideaTransition{}
	}
	owner := idea.OwnerID == viewerID
	t := func(to, label, hint string, primary bool, reason string) ideaTransition {
		return ideaTransition{To: to, Label: label, Hint: hint, Primary: primary,
			Allowed: reason == "", Reason: reason}
	}

	switch idea.Status {
	case models.IdeaDraft:
		reason := ""
		if !owner {
			reason = fmt.Sprintf("Only the author can submit it — or it goes to review on its own at %d votes.", reviewThreshold)
		}
		return []ideaTransition{
			t(models.IdeaUnderReview, "Submit for review", "Ask the community to review this idea.", true, reason),
		}

	case models.IdeaUnderReview:
		approveReason := ""
		switch {
		case owner:
			approveReason = fmt.Sprintf("You can't approve your own idea — another member has to, or it's approved automatically at %d votes.", approveThreshold)
		case idea.Score <= 0:
			approveReason = "It needs more support than opposition before it can be approved."
		}
		back := t(models.IdeaDraft, "Send back to draft", "Ask the author to rework it.", false, "")
		if owner {
			back = t(models.IdeaDraft, "Withdraw to draft", "Take it out of review to keep working on it.", false, "")
		}
		return []ideaTransition{
			t(models.IdeaApproved, "Approve", "Agree it's worth doing.", true, approveReason),
			back,
		}

	case models.IdeaApproved:
		return []ideaTransition{
			t(models.IdeaInProgress, "Start work", "Mark that someone is building it.", true, ""),
			t(models.IdeaUnderReview, "Reopen review", "Put the decision back up for discussion.", false, ""),
		}

	case models.IdeaInProgress:
		completeReason := ""
		if n := openTaskCount(idea.ID); n > 0 {
			completeReason = fmt.Sprintf("%d linked task(s) still open — finish them first.", n)
		}
		return []ideaTransition{
			t(models.IdeaCompleted, "Mark completed", "The work has shipped.", true, completeReason),
			t(models.IdeaApproved, "Pause work", "Move it back to approved.", false, ""),
		}

	case models.IdeaCompleted:
		return []ideaTransition{
			t(models.IdeaInProgress, "Reopen", "More work turned out to be needed.", false, ""),
		}
	}
	return []ideaTransition{}
}

// checkTransition reports whether viewer may move the idea to `to`, and if
// not, why. An empty string means yes.
func checkTransition(idea models.Idea, viewerID uint, to string) string {
	if idea.MergedIntoID != nil {
		return "this idea was merged into another — change the surviving one"
	}
	if idea.ArchivedAt != nil || idea.Status == models.IdeaArchived {
		return "restore the idea before changing its status"
	}
	if to == models.IdeaArchived {
		return "archive it with a reason using the Archive action"
	}
	for _, tr := range transitionsFor(idea, viewerID) {
		if tr.To == to {
			if !tr.Allowed {
				return tr.Reason
			}
			return ""
		}
	}
	return fmt.Sprintf("an idea can't move from %s to %s", statusLabel(idea.Status), statusLabel(to))
}

// nextStepHint is the one line under the status that says what happens next.
func nextStepHint(idea models.Idea) string {
	if idea.MergedIntoID != nil {
		return "Merged into another idea — discussion and votes continue there."
	}
	if idea.ArchivedAt != nil {
		return "Archived. Restore it to put it back on the board."
	}
	switch idea.Status {
	case models.IdeaDraft:
		need := reviewThreshold - idea.Score
		if need < 1 {
			need = 1
		}
		return fmt.Sprintf("Draft — the author can submit it for review, or it moves there on its own with %d more vote(s).", need)
	case models.IdeaUnderReview:
		need := approveThreshold - idea.Score
		if need < 1 {
			need = 1
		}
		return fmt.Sprintf("Under review — any member other than the author can approve it, or it's approved on its own with %d more vote(s).", need)
	case models.IdeaApproved:
		return "Approved — start work or convert it to a task to move it into progress."
	case models.IdeaInProgress:
		return "In progress — it completes automatically when every linked task is done."
	case models.IdeaCompleted:
		return "Completed. Reopen it if more work is needed."
	}
	return ""
}

// autoAdvance applies the vote-driven escalations when a tally rises across a
// threshold. Crossing, not sitting above: an idea the author withdrew to draft
// at 7 votes must not bounce straight back into review on the next vote.
// Returns true when the status moved.
func autoAdvance(idea *models.Idea, actorID uint, prevScore int) bool {
	crossed := func(threshold int) bool { return prevScore < threshold && idea.Score >= threshold }
	switch {
	case idea.Status == models.IdeaDraft && crossed(reviewThreshold):
		idea.Status = models.IdeaUnderReview
		logIdeaEvent(idea.ID, actorID, "vote_threshold", "status",
			models.IdeaDraft, models.IdeaUnderReview,
			fmt.Sprintf("reached %d votes", reviewThreshold))
		notifyThread(*idea, 0, "🚀", "Idea flagged for review",
			fmt.Sprintf("\"%s\" passed %d votes and moved to Under review.", idea.Title, reviewThreshold))
		return true
	case idea.Status == models.IdeaUnderReview && crossed(approveThreshold):
		idea.Status = models.IdeaApproved
		logIdeaEvent(idea.ID, actorID, "vote_threshold", "status",
			models.IdeaUnderReview, models.IdeaApproved,
			fmt.Sprintf("reached %d votes", approveThreshold))
		notifyThread(*idea, 0, "🎉", "Idea approved",
			fmt.Sprintf("\"%s\" passed %d votes and was approved by the community.", idea.Title, approveThreshold))
		return true
	}
	return false
}
