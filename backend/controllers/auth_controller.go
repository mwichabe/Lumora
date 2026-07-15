package controllers

import (
	"fmt"
	"os"
	"path/filepath"
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

	// Send the welcome email in the background — never block registration on it.
	go utils.SendWelcomeEmail(a.Cfg, user.Email, user.Name)
	// Drop an in-app welcome notification from Lumora.
	DeliverWelcome(user.ID)

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

	// Greet the returning user, and email a sign-in alert (throttled to ~3h).
	if DeliverLoginWelcome(user) {
		go utils.SendLoginEmail(a.Cfg, user.Email, user.Name)
	}

	return a.tokenResponse(c, user)
}

// Me returns the authenticated user (with hearts regenerated up to now). Also
// treats a session resume as a "sign-in": greets + emails, throttled to ~3h so
// frequent app opens / background refreshes don't spam.
func (a *AuthController) Me(c *fiber.Ctx) error {
	user := middleware.CurrentUser(c)
	if refreshHearts(user) {
		DeliverHeartsFull(user.ID)
	}
	database.DB.Save(user)

	if DeliverLoginWelcome(*user) {
		go utils.SendLoginEmail(a.Cfg, user.Email, user.Name)
	}
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
		EnsureEnrollment(user.ID, in.TargetLanguage)
	}
	if in.DailyGoalXP > 0 {
		user.DailyGoalXP = in.DailyGoalXP
	}
	database.DB.Save(user)
	return c.JSON(fiber.Map{"user": user})
}

type profileInput struct {
	Name        string `json:"name"`
	AvatarColor string `json:"avatarColor"`
	DailyGoalXP int    `json:"dailyGoalXp"`
}

// UpdateProfile edits the user's display name, avatar colour and daily goal.
func (a *AuthController) UpdateProfile(c *fiber.Ctx) error {
	user := middleware.CurrentUser(c)
	var in profileInput
	if err := c.BodyParser(&in); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
	}
	if name := strings.TrimSpace(in.Name); name != "" {
		user.Name = name
	}
	if in.AvatarColor != "" {
		user.AvatarColor = in.AvatarColor
	}
	if in.DailyGoalXP > 0 {
		user.DailyGoalXP = in.DailyGoalXP
	}
	database.DB.Save(user)
	return c.JSON(fiber.Map{"user": user})
}

var allowedImageExt = map[string]bool{
	".jpg": true, ".jpeg": true, ".png": true, ".webp": true, ".gif": true,
}

// UploadAvatar accepts a profile photo (multipart "file"), stores it on disk
// and points the user's AvatarURL at it. No third-party service required.
func (a *AuthController) UploadAvatar(c *fiber.Ctx) error {
	user := middleware.CurrentUser(c)

	fh, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "no file uploaded"})
	}
	if fh.Size > 5<<20 { // 5 MB
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "image must be under 5MB"})
	}
	ext := strings.ToLower(filepath.Ext(fh.Filename))
	if !allowedImageExt[ext] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "unsupported image type"})
	}

	dir := filepath.Join(a.Cfg.UploadsDir, "avatars")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not store image"})
	}
	name := fmt.Sprintf("user_%d_%d%s", user.ID, time.Now().Unix(), ext)
	if err := c.SaveFile(fh, filepath.Join(dir, name)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not save image"})
	}

	user.AvatarURL = "/uploads/avatars/" + name
	database.DB.Save(user)
	return c.JSON(fiber.Map{"user": user})
}

// RemoveAvatar clears the user's uploaded photo (reverting to the colour
// initial) and best-effort deletes the file from disk.
func (a *AuthController) RemoveAvatar(c *fiber.Ctx) error {
	user := middleware.CurrentUser(c)
	if user.AvatarURL != "" {
		// AvatarURL is a public path ("/uploads/avatars/x.png"); map it back
		// onto the configured directory to get the file on disk.
		rel := strings.TrimPrefix(user.AvatarURL, "/uploads/")
		_ = os.Remove(filepath.Join(a.Cfg.UploadsDir, rel))
		user.AvatarURL = ""
		database.DB.Save(user)
	}
	return c.JSON(fiber.Map{"user": user})
}

type passwordInput struct {
	CurrentPassword string `json:"currentPassword"`
	NewPassword     string `json:"newPassword"`
}

// ChangePassword verifies the current password and sets a new one.
func (a *AuthController) ChangePassword(c *fiber.Ctx) error {
	user := middleware.CurrentUser(c)
	var in passwordInput
	if err := c.BodyParser(&in); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
	}
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(in.CurrentPassword)) != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "current password is incorrect"})
	}
	if len(in.NewPassword) < 6 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "new password must be at least 6 characters"})
	}
	hash, _ := bcrypt.GenerateFromPassword([]byte(in.NewPassword), bcrypt.DefaultCost)
	user.PasswordHash = string(hash)
	database.DB.Save(user)
	return c.JSON(fiber.Map{"ok": true})
}

type forgotInput struct {
	Email string `json:"email"`
}

// ForgotPassword issues a single-use reset link by email. It always responds OK
// (never revealing whether an account exists) to avoid account enumeration.
func (a *AuthController) ForgotPassword(c *fiber.Ctx) error {
	var in forgotInput
	if err := c.BodyParser(&in); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
	}
	email := strings.ToLower(strings.TrimSpace(in.Email))

	var user models.User
	if email != "" && database.DB.Where("email = ?", email).First(&user).Error == nil {
		token := utils.RandomToken(32)
		database.DB.Create(&models.PasswordReset{
			UserID: user.ID, Token: token,
			ExpiresAt: time.Now().Add(time.Hour),
		})
		resetURL := a.Cfg.AppURL + "/reset-password?token=" + token
		go utils.SendPasswordResetEmail(a.Cfg, user.Email, user.Name, resetURL)
	}
	return c.JSON(fiber.Map{"ok": true})
}

type resetInput struct {
	Token    string `json:"token"`
	Password string `json:"password"`
}

// ResetPassword consumes a valid, unexpired token and sets a new password.
func (a *AuthController) ResetPassword(c *fiber.Ctx) error {
	var in resetInput
	if err := c.BodyParser(&in); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
	}
	if len(in.Password) < 6 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "password must be at least 6 characters"})
	}

	var pr models.PasswordReset
	if database.DB.Where("token = ? AND used = ?", in.Token, false).First(&pr).Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "this reset link is invalid or has already been used"})
	}
	if time.Now().After(pr.ExpiresAt) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "this reset link has expired — please request a new one"})
	}

	var user models.User
	if database.DB.First(&user, pr.UserID).Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "account not found"})
	}
	hash, _ := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	user.PasswordHash = string(hash)
	database.DB.Save(&user)
	database.DB.Model(&pr).Update("used", true)

	return c.JSON(fiber.Map{"ok": true})
}

type deleteInput struct {
	Password string `json:"password"`
}

// DeleteAccount permanently removes the user and all of their data.
func (a *AuthController) DeleteAccount(c *fiber.Ctx) error {
	user := middleware.CurrentUser(c)
	var in deleteInput
	if err := c.BodyParser(&in); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
	}
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(in.Password)) != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "password is incorrect"})
	}

	uid := user.ID
	database.DB.Where("user_id = ?", uid).Delete(&models.Enrollment{})
	database.DB.Where("user_id = ?", uid).Delete(&models.Mistake{})
	database.DB.Where("user_id = ?", uid).Delete(&models.Notification{})
	database.DB.Where("user_id = ?", uid).Delete(&models.LessonProgress{})
	database.DB.Where("user_id = ?", uid).Delete(&models.Friendship{})
	database.DB.Where("user_id = ?", uid).Delete(&models.UserQuest{})
	database.DB.Delete(&models.User{}, uid)

	return c.JSON(fiber.Map{"ok": true})
}

func (a *AuthController) tokenResponse(c *fiber.Ctx, user models.User) error {
	token, err := utils.GenerateToken(user.ID, a.Cfg.JWTSecret)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not issue token"})
	}
	return c.JSON(fiber.Map{"token": token, "user": user})
}
