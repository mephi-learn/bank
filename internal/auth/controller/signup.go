package controller

import (
	"bank/internal/helpers"
	"bank/internal/models"
	"encoding/json"
	"net/http"
)

func (h *handler) signUp(w http.ResponseWriter, r *http.Request) {
	var user models.Client

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		helpers.ErrorMessage(w, "invalid json data", http.StatusBadRequest, err)
		return
	}

	userId, err := h.service.SignUp(r.Context(), &user)
	if err != nil {
		helpers.ErrorMessage(w, "sign up error", http.StatusBadRequest, err)
		return
	}

	helpers.SuccessMessageWithCode(w, "user was created", userId, http.StatusCreated)
}
