package controllers

import (
	"log"
	"strings"

	"lumora/backend/database"
	"lumora/backend/models"
)

// The learner's CEFR level.
//
// A CEFR level is a claim about what someone can do in a language, so it is
// only ever derived from evidence in that language — never from XP. XP counts
// activity: replaying one beginner lesson, practice drills and listening
// sessions all earn it, in any language, and none of that shows proficiency.
//
// A learner stands at the highest of:
//
//   - the stage of the course they've reached. A course whose units are tagged
//     with a level ("A1 · Grundlagen") moves the learner up one level each time
//     they complete EVERY lesson of the level they're on. Courses without
//     tagged units are beginner material and never move anyone past A1.
//   - the highest level they hold a certificate for, i.e. a passed exam.
//
// Everyone starts at A1, and the level is per language: it is recomputed when
// the learner switches course.

var cefrOrder = []string{"A1", "A2", "B1", "B2", "C1", "C2"}

// cefrIndex returns a level's position in cefrOrder, or -1 if it isn't one.
func cefrIndex(level string) int {
	for i, l := range cefrOrder {
		if l == level {
			return i
		}
	}
	return -1
}

// unitLevel extracts the CEFR level a unit is tagged with ("A1 · Grundlagen"
// → "A1"), or "" when the unit isn't tagged.
func unitLevel(unit string) string {
	if len(unit) < 2 {
		return ""
	}
	if code := unit[:2]; cefrIndex(code) >= 0 && (len(unit) == 2 || unit[2] == ' ') {
		return code
	}
	return ""
}

// earnedLevel works out the level userID has earned in lang.
func earnedLevel(userID uint, lang string) string {
	// How far through the levelled course are they?
	var skills []models.Skill
	database.DB.Where("language = ?", lang).Preload("Lessons").Find(&skills)

	var done []models.LessonProgress
	database.DB.Where("user_id = ? AND completed = ?", userID, true).Find(&done)
	completed := map[uint]bool{}
	for _, d := range done {
		completed[d.LessonID] = true
	}

	total := make([]int, len(cefrOrder))
	finished := make([]int, len(cefrOrder))
	for _, s := range skills {
		i := cefrIndex(unitLevel(s.Unit))
		if i < 0 {
			continue
		}
		for _, l := range s.Lessons {
			total[i]++
			if completed[l.ID] {
				finished[i]++
			}
		}
	}
	stage := 0
	for stage < len(cefrOrder)-1 && total[stage] > 0 && finished[stage] == total[stage] {
		stage++
	}

	// A passed exam proves the level outright.
	var certs []models.Certificate
	database.DB.Where("user_id = ? AND language = ?", userID, lang).Find(&certs)
	for _, c := range certs {
		level := c.Level
		if level == "FINAL" { // the comprehensive A1→C2 exam
			level = "C2"
		}
		if i := cefrIndex(level); i > stage {
			stage = i
		}
	}
	return cefrOrder[stage]
}

// syncLevel sets the user's level to what they've earned in their current
// language. It reports whether the level went up, so the caller can celebrate
// it. The caller is responsible for saving the user.
func syncLevel(user *models.User) (advanced bool) {
	lang := user.TargetLanguage
	if lang == "" {
		lang = "es"
	}
	level := earnedLevel(user.ID, lang)
	advanced = cefrIndex(level) > cefrIndex(user.CEFRLevel) && user.CEFRLevel != ""
	user.CEFRLevel = level
	user.LevelName = levelNames[level]
	return advanced
}

// RepairLevels corrects levels that were handed out by the old rule (one CEFR
// level per 100 XP) and withdraws the "Level up! You reached …" notifications
// that rule sent, since they congratulated learners on levels they hadn't
// reached. It runs at startup and is a no-op once everything is consistent.
func RepairLevels() {
	var users []models.User
	database.DB.Select("id", "target_language", "cefr_level", "level_name").Find(&users)
	corrected := 0
	for i := range users {
		u := &users[i]
		before := u.CEFRLevel
		syncLevel(u)
		if u.CEFRLevel != before {
			database.DB.Model(&models.User{}).Where("id = ?", u.ID).
				Updates(map[string]interface{}{"cefr_level": u.CEFRLevel, "level_name": u.LevelName})
			corrected++
		}
	}

	// The old rule's notifications were keyed "levelup_<LEVEL>"; honest ones
	// are keyed per language ("levelup_de_A2") and are left alone.
	oldKeys := make([]string, len(cefrOrder))
	for i, l := range cefrOrder {
		oldKeys[i] = "levelup_" + l
	}
	withdrawn := database.DB.Where("key IN ?", oldKeys).Delete(&models.Notification{}).RowsAffected

	if corrected > 0 || withdrawn > 0 {
		log.Printf("[levels] corrected %d user level(s), withdrew %d XP-based level-up notification(s)", corrected, withdrawn)
	}
}

// levelUpKey is the dedup key for a level-up notification in one language.
func levelUpKey(lang, level string) string {
	return "levelup_" + strings.ToLower(lang) + "_" + level
}
