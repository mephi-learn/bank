package models

import "time"

const (
	TransactionDeposit TransactionType = iota + 1
	TransactionTransfer
)

const (
	TransactionStatusPending TransactionStatus = iota + 1
	TransactionStatusCompleted
	TransactionStatusFailed
)

type (
	TransactionType   int
	TransactionStatus int
)

type Transaction struct {
	ID        uint              `json:"id"`
	AccountID uint              `json:"account_id"`
	Amount    float64           `json:"amount"`
	Type      TransactionType   `json:"type"`
	Status    TransactionStatus `json:"status"`
	CreatedAt time.Time         `json:"created_at"`
}
