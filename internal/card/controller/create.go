package controller

import (
	"bank/internal/helpers"
	"fmt"
	"net/http"
	"strconv"
)

func (h *handler) CreateCard(w http.ResponseWriter, r *http.Request) {
	// Получаем идентификатор счёта пользователя из запроса
	accountID, err := strconv.Atoi(r.PathValue("account_id"))
	if err != nil || accountID <= 0 {
		helpers.ErrorMessage(w, fmt.Sprintf("invalid account: %s", r.PathValue("account_id")), http.StatusBadRequest, err)
		return
	}

	// Создаём карту
	card, err := h.service.CreateCard(r.Context(), uint(accountID))
	if err != nil {
		helpers.ErrorMessage(w, "failed create card", http.StatusBadRequest, err)
		return
	}

	helpers.SuccessMessage(w, "card was created", map[string]any{"card": card})
}

func (h *handler) CreateCardByAdmin(w http.ResponseWriter, r *http.Request) {
	// Получаем идентификатор пользователя из запроса
	accountID, err := strconv.Atoi(r.PathValue("account_id"))
	if err != nil || accountID <= 0 {
		helpers.ErrorMessage(w, fmt.Sprintf("invalid account: %s", r.PathValue("account_id")), http.StatusBadRequest, err)
		return
	}

	// Создаём карту
	card, err := h.service.CreateCardByAdmin(r.Context(), uint(accountID))
	if err != nil {
		helpers.ErrorMessage(w, "failed create card", http.StatusBadRequest, err)
		return
	}

	helpers.SuccessMessage(w, "card was created", map[string]any{"card": card})
}
