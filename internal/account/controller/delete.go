package controller

import (
	"bank/internal/helpers"
	"bank/internal/models"
	"fmt"
	"net/http"
	"strconv"
)

// DeleteAccount удаление счёта.
func (h *handler) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	client, err := models.ClientFromContext(r.Context())
	if err != nil {
		helpers.ErrorMessage(w, "authenticate need", http.StatusBadRequest, nil)
		return
	}

	// Получаем идентификатор пользователя из запроса
	accountId, err := strconv.Atoi(r.PathValue("account"))
	if err != nil {
		helpers.ErrorMessage(w, "failed create account", http.StatusBadRequest, err)
		return
	}

	// Удаляем аккаунт
	account := models.Account{
		ID:       uint(accountId),
		ClientID: client.ID,
	}
	if err = h.service.DeleteAccount(r.Context(), &account); err != nil {
		helpers.ErrorMessage(w, "failed delete account", http.StatusBadRequest, err)
		return
	}

	helpers.SuccessMessage(w, "account was deleted", map[string]any{"client id": client.ID})
}

// DeleteAccountAdmin удаление счёта с любого аккаунта.
func (h *handler) DeleteAccountAdmin(w http.ResponseWriter, r *http.Request) {
	// Получаем идентификатор пользователя из запроса
	clientId, err := strconv.Atoi(r.PathValue("client"))
	if err != nil {
		helpers.ErrorMessage(w, fmt.Sprintf("invalid client: %s", r.PathValue("client")), http.StatusBadRequest, err)
		return
	}

	// Получаем идентификатор пользователя из запроса
	accountId, err := strconv.Atoi(r.PathValue("account"))
	if err != nil {
		helpers.ErrorMessage(w, "failed create account", http.StatusBadRequest, err)
		return
	}

	// Удаляем аккаунт
	account := models.Account{
		ID:       uint(accountId),
		ClientID: uint(clientId),
	}
	if err = h.service.DeleteAccount(r.Context(), &account); err != nil {
		helpers.ErrorMessage(w, "failed delete account", http.StatusBadRequest, err)
		return
	}

	helpers.SuccessMessage(w, "account was deleted", map[string]any{"client id": clientId})
}
