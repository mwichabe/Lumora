package models

import "time"

// Payment records a Paystack transaction. The `Reference` is our idempotency key
// (unique per attempt); status moves pending → success/failed as we verify or
// receive the webhook.
type Payment struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index" json:"userId"`
	Reference string    `gorm:"uniqueIndex" json:"reference"`
	Product   string    `json:"product"`  // e.g. "exam_attempt"
	Level     string    `json:"level"`    // exam level this attempt unlocks (A1..C2, FINAL)
	Amount    int       `json:"amount"`   // in the currency subunit (e.g. KES cents)
	Currency  string    `json:"currency"` // e.g. "KES"
	Status      string `json:"status"`   // pending | success | failed
	Consumed    bool   `json:"consumed"` // a successful attempt is consumed when the exam is submitted
	ReceiptSent bool   `json:"-"`        // the receipt email has been sent (exactly-once)
	Channel     string `json:"channel"`  // card, mobile_money, etc. (from Paystack)
	CreatedAt time.Time `json:"createdAt"`
	PaidAt    time.Time `json:"paidAt"`
}
