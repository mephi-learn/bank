package controller

import (
	"bank/internal/helpers"
	"bank/internal/models"
	"fmt"
	"net/http"
	"strconv"
)

func (h *handler) ListCards(w http.ResponseWriter, r *http.Request) {
	client, err := models.ClientFromContext(r.Context())
	if err != nil {
		helpers.ErrorMessage(w, "authenticate need", http.StatusBadRequest, nil)
		return
	}

	// Получаем список карт пользователя
	cards, err := h.service.ListCards(r.Context(), client.ID)
	if err != nil {
		helpers.ErrorMessage(w, "failed list card", http.StatusBadRequest, err)
		return
	}

	helpers.SuccessMessage(w, "cards", cards)
}

func (h *handler) ListCardsByAdmin(w http.ResponseWriter, r *http.Request) {
	// Получаем идентификатор клиента из запроса
	clientId, err := strconv.Atoi(r.PathValue("client_id"))
	if err != nil {
		helpers.ErrorMessage(w, fmt.Sprintf("invalid client: %s", r.PathValue("client_id")), http.StatusBadRequest, err)
		return
	}

	// Получаем список карт пользователя
	cards, err := h.service.ListCards(r.Context(), uint(clientId))
	if err != nil {
		helpers.ErrorMessage(w, "failed list card", http.StatusBadRequest, err)
		return
	}

	helpers.SuccessMessage(w, "cards", cards)
}
