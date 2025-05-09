package models

import "time"

const (
	PaymentStatusSuccess PaymentStatus = iota + 1
	PaymentStatusPending
	PaymentStatusFailed
)

type PaymentStatus int

type Payment struct {
	ID        int64         `json:"id"`
	CreditID  int64         `json:"credit_id"`
	Amount    float64       `json:"account"`
	StartDate time.Time     `json:"start_date"`
	Status    PaymentStatus `json:"status"`
	CreatedAt time.Time     `json:"created_at"`
}
