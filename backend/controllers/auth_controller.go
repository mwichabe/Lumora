package controllers

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"

	"lumora/backend/config"
	"lumora/backend/database"
	"lumora/backend/middleware"
	"lumora/backend/models"
	"lumora/backend/utils"
)

// AuthController groups authentication endpoints.
type AuthController struct {
	Cfg config.Config
}

type registerInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type loginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// avatarColors cycle through the brand palette for placeholder avatars.
var avatarColors = []string{"#6C3FC5", "#F5A623", "#00C2A8", "#FF5C5C", "#17A3DD"}

// Register creates a new account with starter gamification state.
func (a *AuthController) Register(c *fiber.Ctx) error {
	var in registerInput
	if err := c.BodyParser(&in); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
	}
	in.Email = strings.ToLower(strings.TrimSpace(in.Email))
	if in.Email == "" || len(in.Password) < 6 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "email and a 6+ char password are required"})
	}

	var existing models.User
	if err := database.DB.Where("email = ?", in.Email).First(&existing).Error; err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "email already registered"})
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	user := models.User{
		Email:          in.Email,
		PasswordHash:   string(hash),
		Name:           in.Name,
		AvatarColor:    avatarColors[len(in.Email)%len(avatarColors)],
		NativeLanguage: "en",
		CEFRLevel:      "A1",
		LevelName:      "Spark",
		DailyGoalXP:    20,
		Hearts:         5,
		Gems:           50,
		League:         "Bronze",
		LastActiveDate: time.Now().Format("2006-01-02"),
	}
	if err := database.DB.Create(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not create user"})
	}

	return a.tokenResponse(c, user)
}

// Login authenticates an existing account.
func (a *AuthController) Login(c *fiber.Ctx) error {
	var in loginInput
	if err := c.BodyParser(&in); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
	}
	in.Email = strings.ToLower(strings.TrimSpace(in.Email))

	var user models.User
	if err := database.DB.Where("email = ?", in.Email).First(&user).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid credentials"})
	}
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(in.Password)) != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid credentials"})
	}

	return a.tokenResponse(c, user)
}

// Me returns the authenticated user.
func (a *AuthController) Me(c *fiber.Ctx) error {
	user := middleware.CurrentUser(c)
	return c.JSON(fiber.Map{"user": user})
}

type setupInput struct {
	TargetLanguage string `json:"targetLanguage"`
	Reason         string `json:"reason"`
	DailyGoalXP    int    `json:"dailyGoalXp"`
}

// Setup persists the onboarding choices (language + daily goal).
func (a *AuthController) Setup(c *fiber.Ctx) error {
	user := middleware.CurrentUser(c)
	var in setupInput
	if err := c.BodyParser(&in); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
	}
	if in.TargetLanguage != "" {
		user.TargetLanguage = in.TargetLanguage
	}
	if in.DailyGoalXP > 0 {
		user.DailyGoalXP = in.DailyGoalXP
	}
	database.DB.Save(user)
	return c.JSON(fiber.Map{"user": user})
}

func (a *AuthController) tokenResponse(c *fiber.Ctx, user models.User) error {
	token, err := utils.GenerateToken(user.ID, a.Cfg.JWTSecret)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not issue token"})
	}
	return c.JSON(fiber.Map{"token": token, "user": user})
}
