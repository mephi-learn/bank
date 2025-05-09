package controller

import (
	"bank/internal/helpers"
	"fmt"
	"net/http"
	"time"
)

func (h *handler) CreateCredit(w http.ResponseWriter, r *http.Request) {
	// Создаём кредит
	credit, err := h.service.CreateCredit(r.Context())
	if err != nil {
		helpers.ErrorMessage(w, "failed create credit", http.StatusBadRequest, err)
		return
	}

	helpers.SuccessMessage(w, "credit was created", map[string]any{"credit": credit})
}

func (h *handler) PayCredit(w http.ResponseWriter, r *http.Request) {
	// Получаем дату оплаты
	payDate, err := time.Parse(time.RFC3339, r.PathValue("date"))
	if err != nil || payDate.IsZero() {
		helpers.ErrorMessage(w, fmt.Sprintf("invalid date: %s", r.PathValue("date")), http.StatusBadRequest, err)
		return
	}

	// Гасим кредит
	pay, err := h.service.PayCredit(r.Context(), payDate)
	if err != nil {
		helpers.ErrorMessage(w, "failed pay to credit", http.StatusBadRequest, err)
		return
	}

	helpers.SuccessMessage(w, "credit pay success", pay)
}

func (h *handler) ListCredit(w http.ResponseWriter, r *http.Request) {
	// Выводим список кредитов
	credits, err := h.service.ListCards(r.Context())
	if err != nil {
		helpers.ErrorMessage(w, "failed list credit", http.StatusBadRequest, err)
		return
	}

	helpers.SuccessMessage(w, "credits", credits)
}
