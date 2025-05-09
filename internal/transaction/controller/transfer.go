package controller

import (
	"bank/internal/helpers"
	"fmt"
	"net/http"
	"strconv"
)

// CardToCard перевод с карты на карту.
func (h *handler) CardToCard(w http.ResponseWriter, r *http.Request) {
	// Получаем номер карты источника
	fromCardNum := r.PathValue("from_card_num")
	if len(fromCardNum) != 16 {
		helpers.ErrorMessage(w, fmt.Sprintf("invalid source card number: %s", r.PathValue("from_card_num")), http.StatusBadRequest, nil)
		return
	}

	// Получаем номер карты приёмника
	toCardNum := r.PathValue("to_card_num")
	if len(toCardNum) != 16 {
		helpers.ErrorMessage(w, fmt.Sprintf("invalid destination card number: %s", r.PathValue("to_card_num")), http.StatusBadRequest, nil)
		return
	}

	// Получаем сумму перевода
	amount, err := strconv.ParseFloat(r.PathValue("amount"), 64)
	if err != nil || amount <= 0 {
		helpers.ErrorMessage(w, fmt.Sprintf("invalid amount: %s", r.PathValue("amount")), http.StatusBadRequest, err)
		return
	}

	// Делаем перевод
	transactions, err := h.service.CardByCard(r.Context(), fromCardNum, toCardNum, amount)
	if err != nil {
		helpers.ErrorMessage(w, "failed transaction from card to card", http.StatusBadRequest, err)
		return
	}

	helpers.SuccessMessage(w, "transfer card to card success", transactions)
}

// Deposit пополнение по номеру карты.
func (h *handler) Deposit(w http.ResponseWriter, r *http.Request) {
	// Получаем номер карты
	cardNum := r.PathValue("card_num")
	if len(cardNum) != 16 {
		helpers.ErrorMessage(w, fmt.Sprintf("invalid destination card number: %s", r.PathValue("card_num")), http.StatusBadRequest, nil)
		return
	}

	// Получаем сумму перевода
	amount, err := strconv.ParseFloat(r.PathValue("amount"), 64)
	if err != nil || amount <= 0 {
		helpers.ErrorMessage(w, fmt.Sprintf("invalid amount: %s", r.PathValue("amount")), http.StatusBadRequest, err)
		return
	}

	// Пополняем карту
	transactions, err := h.service.Deposit(r.Context(), cardNum, amount)
	if err != nil {
		helpers.ErrorMessage(w, "failed deposit to card", http.StatusBadRequest, err)
		return
	}

	helpers.SuccessMessage(w, "deposit success", transactions)
}
