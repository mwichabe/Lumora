package database

import (
	"log"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"lumora/backend/models"
)

// DB is the shared GORM handle used by controllers.
var DB *gorm.DB

// Connect opens the SQLite database, runs migrations and seeds starter content.
func Connect(path string) {
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
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
	); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	DB = db
	Seed(db)
	log.Println("database ready")
}
