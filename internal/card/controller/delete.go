package controller

import (
	"bank/internal/helpers"
	"bank/internal/models"
	"fmt"
	"net/http"
	"strconv"
)

func (h *handler) DeleteCard(w http.ResponseWriter, r *http.Request) {
	client, err := models.ClientFromContext(r.Context())
	if err != nil {
		helpers.ErrorMessage(w, "authenticate need", http.StatusBadRequest, nil)
		return
	}

	// Получаем идентификатор карты из запроса
	cardId, err := strconv.Atoi(r.PathValue("card_id"))
	if err != nil {
		helpers.ErrorMessage(w, fmt.Sprintf("invalid card id: %s", r.PathValue("card_id")), http.StatusBadRequest, err)
		return
	}

	// Удаляем карту
	card := models.Card{
		ID:       uint(cardId),
		ClientID: client.ID,
	}
	if err = h.service.DeleteCard(r.Context(), &card); err != nil {
		helpers.ErrorMessage(w, "failed delete card", http.StatusBadRequest, err)
		return
	}

	helpers.SuccessMessage(w, "card has been delete", map[string]any{"client id": client.ID, "card id": cardId})
}

func (h *handler) DeleteCardByAdmin(w http.ResponseWriter, r *http.Request) {
	// Получаем идентификатор карты из запроса
	cardId, err := strconv.Atoi(r.PathValue("card_id"))
	if err != nil {
		helpers.ErrorMessage(w, fmt.Sprintf("invalid card id: %s", r.PathValue("card_id")), http.StatusBadRequest, err)
		return
	}

	// Удаляем карту
	card := models.Card{
		ID: uint(cardId),
	}
	if err = h.service.DeleteCardByAdmin(r.Context(), &card); err != nil {
		helpers.ErrorMessage(w, "failed delete card", http.StatusBadRequest, err)
		return
	}

	helpers.SuccessMessage(w, "card has been delete", map[string]any{"card id": cardId})
}
