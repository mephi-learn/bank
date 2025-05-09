package service

import (
	"bank/internal/models"
	"bank/pkg/errors"
	"context"
)

func (h *authService) List(ctx context.Context) ([]*models.Client, error) {
	user, err := h.repo.List(ctx)
	if err != nil {
		return nil, errors.Errorf("failed list clients: %w", err)
	}

	return user, nil
}
