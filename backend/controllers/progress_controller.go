package controllers

import (
	"time"

	"github.com/gofiber/fiber/v2"

	"lumora/backend/database"
	"lumora/backend/middleware"
	"lumora/backend/models"
)

// ProgressController handles the home feed and lesson completion.
type ProgressController struct{}

// levelNames maps CEFR levels to Lumora-universe names.
var levelNames = map[string]string{
	"A1": "Spark", "A2": "Glow", "B1": "Flame", "B2": "Blaze", "C1": "Aurora", "C2": "Luminary",
}

// Home aggregates everything the home screen needs in one call.
func (p *ProgressController) Home(c *fiber.Ctx) error {
	user := middleware.CurrentUser(c)
	rollOverDay(user)
	if refreshHearts(user) {
		DeliverHeartsFull(user.ID)
	}
	database.DB.Save(user)

	lang := user.TargetLanguage
	if lang == "" {
		lang = "es"
	}

	// "Continue learning": first not-yet-completed lesson in order.
	var done []models.LessonProgress
	database.DB.Where("user_id = ? AND completed = ?", user.ID, true).Find(&done)
	completed := map[uint]bool{}
	for _, d := range done {
		completed[d.LessonID] = true
	}

	var skills []models.Skill
	database.DB.Where("language = ?", lang).Preload("Lessons").Order("order_index asc").Find(&skills)

	var nextLesson *models.Lesson
	var nextSkill *models.Skill
	for si := range skills {
		if user.XP < skills[si].RequiredXP {
			continue
		}
		for li := range skills[si].Lessons {
			if !completed[skills[si].Lessons[li].ID] {
				nextLesson = &skills[si].Lessons[li]
				nextSkill = &skills[si]
				break
			}
		}
		if nextLesson != nil {
			break
		}
	}

	quests := ensureDailyQuests(user.ID)

	return c.JSON(fiber.Map{
		"user":       user,
		"nextLesson": nextLesson,
		"nextSkill":  nextSkill,
		"quests":     quests,
	})
}

type completeInput struct {
	Accuracy int `json:"accuracy"`
}

// CompleteLesson awards XP, advances the streak, updates quests and (if needed)
// promotes the user's CEFR level.
func (p *ProgressController) CompleteLesson(c *fiber.Ctx) error {
	user := middleware.CurrentUser(c)
	rollOverDay(user)

	lessonID := c.Params("id")
	var lesson models.Lesson
	if err := database.DB.First(&lesson, lessonID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "lesson not found"})
	}

	var in completeInput
	_ = c.BodyParser(&in)
	if in.Accuracy < 0 || in.Accuracy > 100 {
		in.Accuracy = 100
	}

	// Record (or update) progress for this lesson.
	var prog models.LessonProgress
	first := database.DB.Where("user_id = ? AND lesson_id = ?", user.ID, lesson.ID).First(&prog).Error != nil
	prog.UserID = user.ID
	prog.LessonID = lesson.ID
	prog.Completed = true
	prog.Accuracy = in.Accuracy
	prog.XPEarned = lesson.XPReward
	prog.CompletedAt = time.Now()
	database.DB.Save(&prog)

	// Award XP (every completion gives XP; the loop is meant to be generous).
	xpGain := lesson.XPReward
	user.XP += xpGain
	user.XPToday += xpGain
	user.Gems += 5
	user.FluencyScore = clamp(user.FluencyScore+in.Accuracy/10, 0, 1000)

	prevStreak := user.Streak
	touchStreak(user)

	// CEFR promotion every 100 XP (demo-friendly thresholds).
	prevCEFR := user.CEFRLevel
	promoteLevel(user)
	user.League = leagueForXP(user.XP)

	database.DB.Save(user)

	// Update quest progress.
	updateQuestsOnLesson(user.ID, in.Accuracy)

	// Milestone notifications: level-up, streak, and finishing a whole unit.
	if user.CEFRLevel != prevCEFR {
		DeliverLevelUp(user.ID, user.CEFRLevel, user.LevelName)
	}
	if user.Streak > prevStreak {
		DeliverStreakMilestone(int(user.ID), user.Streak)
	}
	if unit := unitCompletedByLesson(user.ID, lesson.SkillID); unit != "" {
		DeliverUnitComplete(user.ID, unit)
	}

	leveledUp := first // surface a celebration the first time a lesson is cleared

	return c.JSON(fiber.Map{
		"xpEarned":   xpGain,
		"accuracy":   in.Accuracy,
		"user":       user,
		"firstClear": leveledUp,
	})
}

// --- helpers ---

// touchStreak records activity for "today" and keeps the streak correct:
//   - first ever activity (or after a gap) → streak = 1
//   - activity on the day after the last one → streak + 1
//   - more activity the same day → unchanged (but never left at 0)
// It is self-healing: even if LastActiveDate was pre-set to today (as it is at
// registration) a studying user's streak still becomes at least 1.
func touchStreak(user *models.User) {
	today := time.Now().Format("2006-01-02")
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")

	if user.LastActiveDate == today {
		if user.Streak < 1 {
			user.Streak = 1
		}
		return
	}
	if user.LastActiveDate == yesterday {
		user.Streak++
	} else {
		user.Streak = 1
	}
	user.LastActiveDate = today
}

// rollOverDay resets XPToday when a new calendar day begins. If a full day was
// missed the streak breaks.
func rollOverDay(user *models.User) {
	today := time.Now().Format("2006-01-02")
	if user.LastActiveDate == today {
		return
	}
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	if user.LastActiveDate != "" && user.LastActiveDate != yesterday {
		user.Streak = 0 // missed more than a day
	}
	user.XPToday = 0
	// Hearts no longer refill daily — they regenerate over time (see
	// hearts_controller) or can be purchased.
	database.DB.Save(user)
}

func promoteLevel(user *models.User) {
	order := []string{"A1", "A2", "B1", "B2", "C1", "C2"}
	idx := 0
	for i, l := range order {
		if l == user.CEFRLevel {
			idx = i
		}
	}
	target := user.XP / 100
	if target > len(order)-1 {
		target = len(order) - 1
	}
	if target > idx {
		user.CEFRLevel = order[target]
		user.LevelName = levelNames[user.CEFRLevel]
	}
}

// unitCompletedByLesson returns the unit name if every lesson in every skill of
// the just-completed lesson's unit is now done — otherwise "". (Dedup for the
// notification itself is handled by its campaign key.)
func unitCompletedByLesson(userID, skillID uint) string {
	var skill models.Skill
	if database.DB.First(&skill, skillID).Error != nil || skill.Unit == "" {
		return ""
	}

	var unitSkills []models.Skill
	database.DB.Preload("Lessons").
		Where("language = ? AND unit = ?", skill.Language, skill.Unit).
		Find(&unitSkills)

	completed := map[uint]bool{}
	var done []models.LessonProgress
	database.DB.Where("user_id = ? AND completed = ?", userID, true).Find(&done)
	for _, d := range done {
		completed[d.LessonID] = true
	}

	anyLessons := false
	for _, s := range unitSkills {
		for _, l := range s.Lessons {
			anyLessons = true
			if !completed[l.ID] {
				return ""
			}
		}
	}
	if !anyLessons {
		return ""
	}
	return skill.Unit
}

func clamp(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

// --- leagues -----------------------------------------------------------------

// Leagues are XP bands. A user's league is derived from their total XP so the
// "Bronze League" label always matches who they're actually competing against.
var leagueOrder = []string{"Bronze", "Silver", "Gold", "Sapphire", "Ruby", "Emerald", "Amethyst", "Obsidian"}
var leagueFloors = []int{0, 100, 300, 600, 1000, 1500, 2500, 4000}

func leagueIndexForXP(xp int) int {
	idx := 0
	for i, f := range leagueFloors {
		if xp >= f {
			idx = i
		}
	}
	return idx
}

func leagueForXP(xp int) string {
	return leagueOrder[leagueIndexForXP(xp)]
}

// leagueBand returns the [lo, hi) XP range for the league containing xp.
func leagueBand(xp int) (int, int) {
	i := leagueIndexForXP(xp)
	hi := 1 << 30
	if i+1 < len(leagueFloors) {
		hi = leagueFloors[i+1]
	}
	return leagueFloors[i], hi
}
