package models

type Account struct {
	ID       uint    `json:"id"`        // Идентификатор счёта
	ClientID uint    `json:"client_id"` // Привязано к пользователю
	Number   string  `json:"number"`    // Номер счёта
	Balance  float64 `json:"balance"`   // Баланс счёта
}
