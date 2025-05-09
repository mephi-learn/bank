package models

import "time"

const (
	CreditStatusActive CreditStatus = iota + 1
	CreditStatusClose
	CreditStatusOverdue
)

type CreditStatus int

type Credit struct {
	ID           int64        `json:"id"`
	AccountID    int64        `json:"account_id"`
	Principal    float64      `json:"principal"`
	InterestRate float64      `json:"interest_rate"`
	TermMonths   int          `json:"term_months"`
	StartDate    time.Time    `json:"start_date"`
	Status       CreditStatus `json:"status"`
	CreatedAt    time.Time    `json:"created_at"`
}
