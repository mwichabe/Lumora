package controllers

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"

	"lumora/backend/database"
	"lumora/backend/utils"
)

// Message translation, shared by direct messages and idea threads.
//
// The flow is the same in both places:
//
//	on send   → detect the language offline (free, instant, stored)
//	          → if it isn't English, translate in the background and store it
//	on read   → serve the stored translation alongside the original
//
// Detection is separated from translation on purpose. Detection is free, so it
// runs on every message and the UI can label the language even when
// translation is switched off. Translation costs an API call, so it happens
// once per message and the result is persisted — never re-run per view.
//
// The background goroutine is why a message appears instantly and its
// translation lands a moment later, rather than the sender waiting on a model
// round-trip to see their own message posted.

// translator is the process-wide client, set once at boot by
// InitTranslation. Nil-safe: every method on it tolerates a disabled state.
var translator *utils.Translator

// InitTranslation wires up the translator at startup. Safe to call with no API
// key configured — translation is then simply off and detection still runs.
func InitTranslation(t *utils.Translator) { translator = t }

// translationDTO is what the client renders under a foreign-language message.
type translationDTO struct {
	Lang     string `json:"lang"`     // ISO 639-1 of the original
	LangName string `json:"langName"` // "Spanish"
	Text     string `json:"text"`     // English, "" while still pending
	Pending  bool   `json:"pending"`  // translation was requested but hasn't landed
}

// translationFor builds the DTO for a message body. Returns nil when the text
// is English, undetermined, or too short to judge — in which case the client
// shows nothing at all, which is the right outcome for the overwhelming
// majority of messages.
func translationFor(body, lang, translated string, at *time.Time) *translationDTO {
	if lang == "" || lang == "en" || body == "" {
		return nil
	}
	dto := &translationDTO{
		Lang:     lang,
		LangName: utils.LanguageName(lang),
		Text:     translated,
	}
	// Pending means "we know this isn't English and a translation is coming".
	// If translation is disabled entirely, it isn't coming — say so by not
	// claiming to be pending, so the UI shows just the language badge.
	dto.Pending = translated == "" && at == nil && translator.Enabled()
	return dto
}

// detectAndQueue records the message's language and, when it isn't English,
// kicks off a background translation.
//
// Both models are handled through the same path because the work is identical;
// only the table differs. The write-back is a targeted column update rather
// than a full Save so it can't clobber an edit or a delete that landed while
// the model call was in flight.
func detectAndQueue(table string, id uint, body string) string {
	lang, needs := utils.NeedsTranslation(body)
	if lang.Code == "" {
		return ""
	}

	// Store the detection immediately — it's free and the UI wants the label
	// whether or not a translation follows.
	database.DB.Table(table).Where("id = ?", id).
		Update("detected_lang", lang.Code)

	if !needs || !translator.Enabled() {
		return lang.Code
	}

	go translateInBackground(table, id, body, lang.Code)
	return lang.Code
}

func translateInBackground(table string, id uint, body, from string) {
	defer func() {
		// A panic in a detached goroutine would take the whole server down,
		// and a failed translation is never worth that.
		if r := recover(); r != nil {
			log.Printf("[translate] panic translating %s/%d: %v", table, id, r)
		}
	}()

	out, err := translator.Translate(context.Background(), body, from)
	if err != nil {
		if !errors.Is(err, utils.ErrTranslationDisabled) {
			log.Printf("[translate] %s/%d (%s): %v", table, id, from, err)
		}
		return
	}

	now := time.Now()
	database.DB.Table(table).Where("id = ?", id).Updates(map[string]interface{}{
		"translated_body": out,
		"translated_at":   now,
	})
}

// --- on-demand translation ---------------------------------------------------

// TranslateChatMessage produces (or re-produces) the English translation of a
// direct message. The background pass covers the normal case; this endpoint is
// the retry path for when it failed — a timeout, a transient API error, or a
// message that predates the feature.
func (cc *ChatController) TranslateChatMessage(c *fiber.Ctx) error {
	msg, ok := loadChatMessage(c)
	if !ok {
		return nil // the helper already wrote the error response
	}
	dto, err := translateNow("messages", msg.ID, msg.Body, msg.DetectedLang, msg.TranslatedBody, msg.TranslatedAt)
	if err != nil {
		return translationError(c, err)
	}
	return c.JSON(fiber.Map{"translation": dto})
}

// TranslateIdeaMessage is the same retry path for a message in an idea thread.
func (ic *IdeaController) TranslateIdeaMessage(c *fiber.Ctx) error {
	msg, ok := loadMessage(c)
	if !ok {
		return nil // the helper already wrote the error response
	}
	dto, err := translateNow("idea_messages", msg.ID, msg.Body, msg.DetectedLang, msg.TranslatedBody, msg.TranslatedAt)
	if err != nil {
		return translationError(c, err)
	}
	return c.JSON(fiber.Map{"translation": dto})
}

// translateNow returns a cached translation if there is one, otherwise
// translates synchronously and stores the result.
func translateNow(table string, id uint, body, lang, cached string, at *time.Time) (*translationDTO, error) {
	if cached != "" {
		return translationFor(body, lang, cached, at), nil
	}
	if body == "" {
		return nil, errors.New("nothing to translate")
	}

	// A message stored before this feature existed has no detected language;
	// work it out now rather than refusing.
	if lang == "" {
		detected, needs := utils.NeedsTranslation(body)
		if !needs {
			return nil, errors.New("this message already looks like English")
		}
		lang = detected.Code
		database.DB.Table(table).Where("id = ?", id).Update("detected_lang", lang)
	}

	out, err := translator.Translate(context.Background(), body, lang)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	database.DB.Table(table).Where("id = ?", id).Updates(map[string]interface{}{
		"translated_body": out,
		"translated_at":   now,
	})
	return translationFor(body, lang, out, &now), nil
}

func translationError(c *fiber.Ctx, err error) error {
	if errors.Is(err, utils.ErrTranslationDisabled) {
		return c.Status(fiber.StatusServiceUnavailable).
			JSON(fiber.Map{"error": "translation isn't configured on this server"})
	}
	return c.Status(fiber.StatusBadGateway).
		JSON(fiber.Map{"error": "could not translate that message right now"})
}
