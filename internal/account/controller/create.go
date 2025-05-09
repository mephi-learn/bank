package controller

import (
	"bank/internal/helpers"
	"bank/internal/models"
	"fmt"
	"net/http"
	"strconv"
)

func (h *handler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	client, err := models.ClientFromContext(r.Context())
	if err != nil {
		helpers.ErrorMessage(w, "authenticate need", http.StatusBadRequest, nil)
		return
	}

	// Создаём счёт
	account, err := h.service.CreateAccount(r.Context(), client.ID)
	if err != nil {
		helpers.ErrorMessage(w, "failed create account", http.StatusBadRequest, err)
		return
	}

	helpers.SuccessMessage(w, "account was created", map[string]any{"account id": account.ID, "client id": client.ID})
}

func (h *handler) CreateAccountAdmin(w http.ResponseWriter, r *http.Request) {
	// Получаем идентификатор пользователя из запроса
	clientId, err := strconv.Atoi(r.PathValue("client"))
	if err != nil {
		helpers.ErrorMessage(w, fmt.Sprintf("invalid client: %s", r.PathValue("client")), http.StatusBadRequest, err)
		return
	}

	// Создаём аккаунт
	account, err := h.service.CreateAccount(r.Context(), uint(clientId))
	if err != nil {
		helpers.ErrorMessage(w, "failed create account", http.StatusBadRequest, err)
		return
	}

	helpers.SuccessMessage(w, "account was created", map[string]any{"account id": account.ID, "client id": uint(clientId)})
}
