package service

import (
	"bank/internal/models"
	"bank/pkg/errors"
	"context"
)

func (h *authService) ClientById(ctx context.Context, userId uint) (*models.Client, error) {
	user, err := h.repo.GetById(ctx, userId)
	if err != nil {
		return nil, errors.Errorf("failed get client: %w", err)
	}

	return user, nil
}
