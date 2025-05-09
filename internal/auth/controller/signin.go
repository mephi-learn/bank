package controller

import (
	"bank/internal/helpers"
	"bank/internal/models"
	"encoding/json"
	"net/http"
)

func (h *handler) signIn(w http.ResponseWriter, r *http.Request) {
	var signIn models.Client

	err := json.NewDecoder(r.Body).Decode(&signIn)
	if err != nil {
		helpers.ErrorMessage(w, "invalid json data", http.StatusBadRequest, err)
		return
	}

	signedToken, err := h.service.SignIn(r.Context(), &signIn)
	if err != nil {
		helpers.ErrorMessage(w, "user was not found", http.StatusBadRequest, nil)
		return
	}

	helpers.SuccessMessage(w, "token create", signedToken)
}
