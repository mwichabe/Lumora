package controllers

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"

	"lumora/backend/database"
	"lumora/backend/models"
)

// patchStatus sends PATCH /ideas/:id {"status": to} as `as` and returns the
// HTTP status and the reloaded idea.
func patchStatus(t *testing.T, as *models.User, ideaID uint, to string) (int, models.Idea) {
	t.Helper()
	app := fiber.New()
	ic := &IdeaController{}
	app.Patch("/ideas/:id", func(c *fiber.Ctx) error {
		c.Locals("user", as)
		return ic.Update(c)
	})
	req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/ideas/%d", ideaID),
		strings.NewReader(fmt.Sprintf(`{"status":%q}`, to)))
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("request: %v", err)
	}
	_, _ = io.ReadAll(res.Body)
	var idea models.Idea
	database.DB.First(&idea, ideaID)
	return res.StatusCode, idea
}

func workflowUsers(t *testing.T) (owner, other *models.User) {
	t.Helper()
	owner = &models.User{Email: "owner@test.dev", Name: "Owner"}
	other = &models.User{Email: "other@test.dev", Name: "Other"}
	database.DB.Create(owner)
	database.DB.Create(other)
	return owner, other
}

// The full happy path: draft → review → approved → in progress → completed.
func TestIdeaWalksTheWholeLifecycle(t *testing.T) {
	newIdeaDB(t)
	owner, other := workflowUsers(t)
	idea := makeIdea(t, owner.ID, "Offline lessons", "Download a unit to study without data.")
	applyVote(idea, owner.ID, 1)

	if code, got := patchStatus(t, owner, idea.ID, models.IdeaUnderReview); code != 200 || got.Status != models.IdeaUnderReview {
		t.Fatalf("submit: code %d status %s", code, got.Status)
	}
	if code, got := patchStatus(t, other, idea.ID, models.IdeaApproved); code != 200 || got.Status != models.IdeaApproved {
		t.Fatalf("approve: code %d status %s", code, got.Status)
	}
	if code, got := patchStatus(t, other, idea.ID, models.IdeaInProgress); code != 200 || got.Status != models.IdeaInProgress {
		t.Fatalf("start: code %d status %s", code, got.Status)
	}
	code, got := patchStatus(t, other, idea.ID, models.IdeaCompleted)
	if code != 200 || got.Status != models.IdeaCompleted {
		t.Fatalf("complete: code %d status %s", code, got.Status)
	}
	// Status-only updates must leave the content alone.
	if got.Description != "Download a unit to study without data." {
		t.Errorf("description changed to %q by a status move", got.Description)
	}
}

func TestOnlyTheAuthorCanSubmitADraft(t *testing.T) {
	newIdeaDB(t)
	owner, other := workflowUsers(t)
	idea := makeIdea(t, owner.ID, "Offline lessons", "")

	if code, got := patchStatus(t, other, idea.ID, models.IdeaUnderReview); code != http.StatusConflict || got.Status != models.IdeaDraft {
		t.Fatalf("a non-author submitted the draft: code %d status %s", code, got.Status)
	}
}

func TestAuthorsCannotApproveTheirOwnIdea(t *testing.T) {
	newIdeaDB(t)
	owner, _ := workflowUsers(t)
	idea := makeIdea(t, owner.ID, "Offline lessons", "")
	applyVote(idea, owner.ID, 1)
	database.DB.Model(idea).Update("status", models.IdeaUnderReview)

	if code, got := patchStatus(t, owner, idea.ID, models.IdeaApproved); code != http.StatusConflict || got.Status != models.IdeaUnderReview {
		t.Fatalf("author approved their own idea: code %d status %s", code, got.Status)
	}
}

func TestApprovalNeedsNetSupport(t *testing.T) {
	newIdeaDB(t)
	owner, other := workflowUsers(t)
	idea := makeIdea(t, owner.ID, "Offline lessons", "")
	applyVote(idea, other.ID, -1) // score -1
	database.DB.Model(idea).Update("status", models.IdeaUnderReview)

	if code, _ := patchStatus(t, other, idea.ID, models.IdeaApproved); code != http.StatusConflict {
		t.Fatalf("approved an idea with negative support: code %d", code)
	}
}

func TestStatusCannotSkipSteps(t *testing.T) {
	newIdeaDB(t)
	owner, _ := workflowUsers(t)
	idea := makeIdea(t, owner.ID, "Offline lessons", "")

	for _, to := range []string{models.IdeaApproved, models.IdeaCompleted, models.IdeaArchived} {
		if code, got := patchStatus(t, owner, idea.ID, to); code == 200 || got.Status != models.IdeaDraft {
			t.Errorf("draft → %s was allowed (code %d, status %s)", to, code, got.Status)
		}
	}
}

func TestCompletingNeedsEveryTaskDone(t *testing.T) {
	newIdeaDB(t)
	owner, other := workflowUsers(t)
	idea := makeIdea(t, owner.ID, "Offline lessons", "")
	database.DB.Model(idea).Update("status", models.IdeaInProgress)
	database.DB.Create(&models.IdeaTask{IdeaID: idea.ID, Title: "Build it", Status: "todo"})

	if code, got := patchStatus(t, other, idea.ID, models.IdeaCompleted); code != http.StatusConflict || got.Status != models.IdeaInProgress {
		t.Fatalf("completed with an open task: code %d status %s", code, got.Status)
	}
}

// Votes carry an idea all the way to approval without anyone having to act.
func TestVotesEscalateDraftThroughToApproved(t *testing.T) {
	newIdeaDB(t)
	idea := makeIdea(t, 1, "Dark mode", "")

	for i := 0; i < approveThreshold; i++ {
		applyVote(idea, uint(100+i), 1)
	}
	if idea.Status != models.IdeaApproved {
		t.Fatalf("status = %s at %d votes, want approved", idea.Status, idea.Score)
	}
}

// A withdrawn idea that already has plenty of votes mustn't bounce straight
// back into review on the next vote — escalation is on crossing, not sitting
// above, the threshold.
func TestWithdrawnIdeaStaysInDraft(t *testing.T) {
	newIdeaDB(t)
	idea := makeIdea(t, 1, "Dark mode", "")
	for i := 0; i < reviewThreshold+1; i++ {
		applyVote(idea, uint(100+i), 1)
	}
	idea.Status = models.IdeaDraft
	database.DB.Save(idea)

	applyVote(idea, 999, 1)
	if idea.Status != models.IdeaDraft {
		t.Errorf("withdrawn idea re-escalated to %s", idea.Status)
	}
}

func TestRestoreReturnsToTheStepItLeftFrom(t *testing.T) {
	newIdeaDB(t)
	owner, _ := workflowUsers(t)
	idea := makeIdea(t, owner.ID, "Dark mode", "")
	database.DB.Model(idea).Update("status", models.IdeaApproved)

	app := fiber.New()
	ic := &IdeaController{}
	app.Post("/ideas/:id/:action", func(c *fiber.Ctx) error {
		c.Locals("user", owner)
		if c.Params("action") == "archive" {
			return ic.Archive(c)
		}
		return ic.Restore(c)
	})
	for _, action := range []string{"archive", "restore"} {
		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/ideas/%d/%s", idea.ID, action),
			strings.NewReader(`{"reason":"parked for now"}`))
		req.Header.Set("Content-Type", "application/json")
		res, _ := app.Test(req, -1)
		if res.StatusCode != 200 {
			body, _ := io.ReadAll(res.Body)
			t.Fatalf("%s: status %d (%s)", action, res.StatusCode, body)
		}
	}
	var got models.Idea
	database.DB.First(&got, idea.ID)
	if got.Status != models.IdeaApproved || got.ArchivedAt != nil {
		t.Errorf("restored to %s (archived=%v), want approved", got.Status, got.ArchivedAt != nil)
	}
}

func TestTasksNeedAnApprovedIdea(t *testing.T) {
	newIdeaDB(t)
	owner, _ := workflowUsers(t)
	idea := makeIdea(t, owner.ID, "Dark mode", "")

	app := fiber.New()
	ic := &IdeaController{}
	app.Post("/ideas/:id/tasks", func(c *fiber.Ctx) error {
		c.Locals("user", owner)
		return ic.CreateTask(c)
	})
	res, _ := app.Test(httptest.NewRequest(http.MethodPost,
		fmt.Sprintf("/ideas/%d/tasks", idea.ID), nil), -1)
	if res.StatusCode != http.StatusConflict {
		t.Fatalf("draft became a task: status %d", res.StatusCode)
	}
}
