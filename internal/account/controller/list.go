package controller

import (
	"bank/internal/helpers"
	"bank/internal/models"
	"fmt"
	"net/http"
	"strconv"
)

func (h *handler) ListAccount(w http.ResponseWriter, r *http.Request) {
	client, err := models.ClientFromContext(r.Context())
	if err != nil {
		helpers.ErrorMessage(w, "authenticate need", http.StatusBadRequest, nil)
		return
	}

	// Получаем список счетов пользователя
	accounts, err := h.service.ListAccounts(r.Context(), client.ID)
	if err != nil {
		helpers.ErrorMessage(w, "failed list account", http.StatusBadRequest, err)
		return
	}

	helpers.SuccessMessage(w, "accounts", accounts)
}

func (h *handler) ListAccountAdmin(w http.ResponseWriter, r *http.Request) {
	// Получаем идентификатор пользователя из запроса
	clientId, err := strconv.Atoi(r.PathValue("client"))
	if err != nil {
		helpers.ErrorMessage(w, fmt.Sprintf("invalid client: %s", r.PathValue("client")), http.StatusBadRequest, err)
		return
	}

	// Получаем список счетов пользователя
	accounts, err := h.service.ListAccounts(r.Context(), uint(clientId))
	if err != nil {
		helpers.ErrorMessage(w, "failed list account", http.StatusBadRequest, err)
		return
	}

	helpers.SuccessMessage(w, "accounts", accounts)
}
