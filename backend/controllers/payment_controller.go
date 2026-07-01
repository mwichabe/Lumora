package controllers

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"

	"lumora/backend/config"
	"lumora/backend/database"
	"lumora/backend/middleware"
	"lumora/backend/models"
	"lumora/backend/utils"
)

// PaymentController handles Paystack checkout for exam attempts. Pricing is
// per-level (higher levels cost more; the comprehensive FINAL costs the most),
// and every attempt — including retakes — must be paid for.
type PaymentController struct {
	Cfg config.Config
}

const paystackBase = "https://api.paystack.co"
const productExamAttempt = "exam_attempt"
const productHearts = "hearts_refill"
const heartsRefillPriceKES = 100 // one full refill of 5 hearts

// examPricesKES is the price (in whole KES) of one attempt at each level.
// Higher level → higher price; the FINAL mastery exam is the priciest.
var examPricesKES = map[string]int{
	"A1": 300, "A2": 400, "B1": 550, "B2": 700, "C1": 900, "C2": 1100,
	"FINAL": 2500,
}

func examPriceKES(level string) int {
	if p, ok := examPricesKES[level]; ok {
		return p
	}
	return 500
}

// PaymentsEnabled reports whether a Paystack secret key is configured. When it
// isn't, exams stay free (local development).
func PaymentsEnabled() bool {
	return os.Getenv("PAYSTACK_SECRET_KEY") != ""
}

// hasUnconsumedAttempt reports whether the user has a paid, not-yet-used attempt
// for a level.
func hasUnconsumedAttempt(userID uint, level string) bool {
	var n int64
	database.DB.Model(&models.Payment{}).
		Where("user_id = ? AND level = ? AND status = ? AND consumed = ?",
			userID, level, "success", false).
		Count(&n)
	return n > 0
}

// consumePaidAttempt marks one paid attempt for a level as used. Returns false
// if there's nothing to consume.
func consumePaidAttempt(userID uint, level string) bool {
	var p models.Payment
	if database.DB.
		Where("user_id = ? AND level = ? AND status = ? AND consumed = ?",
			userID, level, "success", false).
		Order("created_at asc").First(&p).Error != nil {
		return false
	}
	database.DB.Model(&p).Update("consumed", true)
	return true
}

// Status tells the frontend whether payments are on, the price of each level
// (in KES and an approximate USD equivalent for international users), and which
// levels the user has already paid for (an unconsumed attempt).
func (pc *PaymentController) Status(c *fiber.Ctx) error {
	user := middleware.CurrentUser(c)

	rate := pc.Cfg.KESPerUSD
	if rate <= 0 {
		rate = 130
	}

	prices := fiber.Map{}
	pricesUsd := fiber.Map{}
	paid := fiber.Map{}
	for level := range levelPassMark {
		kes := examPriceKES(level)
		prices[level] = kes
		pricesUsd[level] = math.Round(float64(kes)/float64(rate)*100) / 100
		paid[level] = hasUnconsumedAttempt(user.ID, level)
	}

	return c.JSON(fiber.Map{
		"paymentsEnabled": pc.Cfg.PaystackSecret != "",
		"currency":        "KES",
		"prices":          prices,
		"pricesUsd":       pricesUsd,
		"paid":            paid,
	})
}

type initInput struct {
	Product string `json:"product"` // "hearts" for a hearts refill; otherwise an exam attempt
	Level   string `json:"level"`
}

type paystackInitData struct {
	AuthorizationURL string `json:"authorization_url"`
	AccessCode       string `json:"access_code"`
	Reference        string `json:"reference"`
}

// Initialize creates a pending payment (an exam attempt, or a hearts refill) and
// asks Paystack for a checkout URL.
func (pc *PaymentController) Initialize(c *fiber.Ctx) error {
	user := middleware.CurrentUser(c)
	if pc.Cfg.PaystackSecret == "" {
		return c.Status(fiber.StatusServiceUnavailable).
			JSON(fiber.Map{"error": "payments are not configured"})
	}

	var in initInput
	_ = c.BodyParser(&in)

	// Decide product, price and reference.
	var product, level, reference string
	var priceKes int
	if in.Product == "hearts" {
		product = productHearts
		priceKes = heartsRefillPriceKES
		reference = fmt.Sprintf("LUM-HEARTS-%d-%d", user.ID, time.Now().UnixNano())
	} else {
		product = productExamAttempt
		level = in.Level
		if _, ok := levelPassMark[level]; !ok {
			level = "A1"
		}
		priceKes = examPriceKES(level)
		reference = fmt.Sprintf("LUM-%s-%d-%d", level, user.ID, time.Now().UnixNano())
	}
	amount := priceKes * 100 // Paystack expects the currency subunit

	body := map[string]interface{}{
		"email":        user.Email,
		"amount":       amount,
		"currency":     "KES",
		"reference":    reference,
		"callback_url": pc.Cfg.AppURL + "/payment/callback",
		"metadata": map[string]interface{}{
			"userId":  user.ID,
			"product": product,
			"level":   level,
		},
	}

	var res struct {
		Status bool             `json:"status"`
		Data   paystackInitData `json:"data"`
	}
	if err := pc.paystack(http.MethodPost, "/transaction/initialize", body, &res); err != nil || !res.Status {
		return c.Status(fiber.StatusBadGateway).
			JSON(fiber.Map{"error": "could not start payment"})
	}

	database.DB.Create(&models.Payment{
		UserID: user.ID, Reference: reference, Product: product,
		Level: level, Amount: amount, Currency: "KES", Status: "pending",
	})

	return c.JSON(fiber.Map{
		"authorizationUrl": res.Data.AuthorizationURL,
		"reference":        res.Data.Reference,
	})
}

// Verify confirms a reference with Paystack and marks the attempt paid on
// success. Used by the callback page; the webhook is the authoritative backup.
func (pc *PaymentController) Verify(c *fiber.Ctx) error {
	user := middleware.CurrentUser(c)
	reference := c.Query("reference")
	if reference == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "missing reference"})
	}

	var res struct {
		Status bool `json:"status"`
		Data   struct {
			Status    string `json:"status"`
			Reference string `json:"reference"`
			Channel   string `json:"channel"`
		} `json:"data"`
	}
	if err := pc.paystack(http.MethodGet, "/transaction/verify/"+reference, nil, &res); err != nil || !res.Status {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"error": "could not verify payment"})
	}

	success := res.Data.Status == "success"
	if success {
		pc.fulfillPayment(reference, res.Data.Channel)
	} else {
		DeliverPaymentFailed(user.ID, reference)
	}

	var p models.Payment
	database.DB.Where("reference = ?", reference).First(&p)
	return c.JSON(fiber.Map{
		"status":  res.Data.Status,
		"success": success,
		"level":   p.Level,
		"product": p.Product,
	})
}

// Webhook receives Paystack events, verifies the signature, and marks the
// attempt paid on charge.success (works even if the user closes the tab).
func (pc *PaymentController) Webhook(c *fiber.Ctx) error {
	body := c.Body()

	mac := hmac.New(sha512.New, []byte(pc.Cfg.PaystackSecret))
	mac.Write(body)
	expected := hex.EncodeToString(mac.Sum(nil))
	if !hmac.Equal([]byte(expected), []byte(c.Get("x-paystack-signature"))) {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	var evt struct {
		Event string `json:"event"`
		Data  struct {
			Reference string `json:"reference"`
			Status    string `json:"status"`
			Channel   string `json:"channel"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &evt); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if evt.Event == "charge.success" && evt.Data.Status == "success" {
		pc.fulfillPayment(evt.Data.Reference, evt.Data.Channel)
	}
	return c.SendStatus(fiber.StatusOK) // always 200 so Paystack stops retrying
}

// fulfillPayment marks a payment successful (an unconsumed, ready-to-use exam
// attempt or a hearts refill), once, notifies the user in-app and emails a
// receipt on the first successful fulfilment.
func (pc *PaymentController) fulfillPayment(reference, channel string) {
	var p models.Payment
	if database.DB.Where("reference = ?", reference).First(&p).Error != nil {
		return // unknown reference — ignore
	}
	alreadyDone := p.Status == "success"
	if !alreadyDone {
		p.Status = "success"
		p.Channel = channel
		p.PaidAt = time.Now()
		database.DB.Save(&p)
	}

	log.Printf("[payment] fulfilling %s product=%s alreadyDone=%v receiptSent=%v",
		reference, p.Product, alreadyDone, p.ReceiptSent)

	switch p.Product {
	case productHearts:
		grantFullHearts(p.UserID)
		if !alreadyDone {
			DeliverHeartsPurchased(p.UserID)
		}
	default: // exam attempt
		DeliverPaymentSuccess(p.UserID, reference)
	}

	// Email a receipt exactly once. Sent synchronously so failures are logged
	// and we only mark it sent on success (a later verify/webhook call retries).
	if !p.ReceiptSent {
		var u models.User
		if database.DB.First(&u, p.UserID).Error == nil {
			item := "Exam attempt"
			if p.Level != "" {
				item = p.Level + " exam attempt"
			}
			if p.Product == productHearts {
				item = "Hearts refill (full set)"
			}
			amountLabel := fmt.Sprintf("%s %d", p.Currency, p.Amount/100)
			if err := utils.SendPaymentEmail(pc.Cfg, u.Email, u.Name, item, amountLabel); err == nil {
				database.DB.Model(&p).Update("receipt_sent", true)
			}
		}
	}
}

// paystack performs an authenticated JSON request against the Paystack API.
func (pc *PaymentController) paystack(method, path string, body interface{}, out interface{}) error {
	var reqBody io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return err
		}
		reqBody = bytes.NewReader(b)
	}

	req, err := http.NewRequest(method, paystackBase+path, reqBody)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+pc.Cfg.PaystackSecret)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 20 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if out != nil {
		return json.Unmarshal(data, out)
	}
	return nil
}
