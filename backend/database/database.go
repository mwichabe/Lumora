package database

import (
	"log"

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
	); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	DB = db
	Seed(db)
	log.Println("database ready")
}
