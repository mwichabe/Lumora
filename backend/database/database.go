package database

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"lumora/backend/models"
)

// DB is the shared GORM handle used by controllers.
var DB *gorm.DB

// Connect opens the database, runs migrations and seeds starter content.
//
// The driver is chosen by databaseURL: when it's set (production — Render passes
// a Postgres connection string) we use Postgres; when it's empty we fall back to
// a local SQLite file at path, so `go run .` still works with no setup.
func Connect(databaseURL, path string) {
	// On a managed host the SQLite fallback is a trap: the filesystem is
	// ephemeral, so the app would boot, look perfectly healthy, and silently
	// discard every account on the next deploy. Refuse to start instead.
	if databaseURL == "" && os.Getenv("RENDER") == "true" {
		log.Fatal("DATABASE_URL is not set. Refusing to start on the SQLite " +
			"fallback: Render's filesystem is ephemeral and all data would be " +
			"lost on the next deploy. Set DATABASE_URL to your Postgres " +
			"connection string (Neon), including ?sslmode=require.")
	}

	var (
		db  *gorm.DB
		err error
	)
	if databaseURL != "" {
		log.Println("connecting to postgres")
		db, err = gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	} else {
		log.Printf("connecting to sqlite at %s", path)
		db, err = gorm.Open(sqlite.Open(path), &gorm.Config{})
	}
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	if err := db.AutoMigrate(
		&models.User{},
		&models.Skill{},
		&models.Lesson{},
		&models.Exercise{},
		&models.VocabItem{},
		&models.ListeningSession{},
		&models.ListeningMatch{},
		&models.ListeningLine{},
		&models.ListeningQuestion{},
		&models.ReadingSession{},
		&models.ReadingLine{},
		&models.ReadingQuestion{},
		&models.Enrollment{},
		&models.Mistake{},
		&models.Notification{},
		&models.Certificate{},
		&models.Message{},
		&models.LessonProgress{},
		&models.Character{},
		&models.Friendship{},
		&models.Quest{},
		&models.UserQuest{},
		&models.Payment{},
		&models.PasswordReset{},
		&models.LeagueSeason{},
		&models.LeagueMembership{},
		&models.LeagueDaily{},
		&models.LeagueReport{},
		&models.Idea{},
		&models.IdeaVote{},
		&models.IdeaStar{},
		&models.IdeaTag{},
		&models.IdeaMessage{},
		&models.IdeaReaction{},
		&models.IdeaEvent{},
		&models.IdeaTask{},
		&models.BrainstormSession{},
	); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	DB = db
	healLegacyAvatars(db)
	Seed(db)
	log.Println("database ready")
}

// healLegacyAvatars repairs accounts left over from when profile photos were
// written to disk under /uploads/avatars and referenced by that path.
//
// Photos now live in the database and are served by GET /api/avatars/:id —
// nothing serves /uploads any more, so those rows render a broken image. On
// Render the files are gone for good (ephemeral filesystem), but on a machine
// that still has them we can rescue the picture rather than dropping it:
//
//	file still on disk -> read it into the row, repoint the URL
//	file gone          -> clear the URL so the UI falls back to the initial
//
// Idempotent: once a row is healed it no longer matches, so restarts are free.
func healLegacyAvatars(db *gorm.DB) {
	var users []models.User
	if err := db.Where("avatar_url LIKE ?", "/uploads/%").Find(&users).Error; err != nil {
		return
	}
	if len(users) == 0 {
		return
	}

	rescued, cleared := 0, 0
	for _, u := range users {
		// Only the base name is trusted — the stored path is user-influenced and
		// must never be able to walk out of the uploads directory.
		name := filepath.Base(u.AvatarURL)
		data, err := os.ReadFile(filepath.Join("uploads", "avatars", name))

		if err != nil || len(data) == 0 {
			db.Model(&models.User{}).Where("id = ?", u.ID).
				Update("avatar_url", "")
			cleared++
			continue
		}

		mime := "image/jpeg"
		if strings.EqualFold(filepath.Ext(name), ".png") {
			mime = "image/png"
		}
		db.Model(&models.User{}).Where("id = ?", u.ID).Updates(map[string]interface{}{
			"avatar_data": data,
			"avatar_mime": mime,
			// The ?v= stamp busts the immutable cache GetAvatar sets.
			"avatar_url": fmt.Sprintf("/api/avatars/%d?v=%d", u.ID, time.Now().Unix()),
		})
		rescued++
	}
	log.Printf("[avatars] healed %d legacy rows: %d imported from disk, %d cleared to initials",
		len(users), rescued, cleared)
}
