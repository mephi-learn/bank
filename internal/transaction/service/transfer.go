package service

import (
	"bank/internal/models"
	"bank/pkg/errors"
	"context"
)

// CardByCard перевод с карты на карту.
func (s *transactionService) CardByCard(ctx context.Context, fromCard string, toCard string, amount float64) ([]*models.Transaction, error) {
	// Получаем текущего клиента
	client, err := models.ClientFromContext(ctx)
	if err != nil {
		return nil, errors.New("authenticate need")
	}

	// Проверяем карту источник
	cardFrom, err := s.card.ValidateCardByNum(ctx, fromCard)
	if err != nil {
		return nil, errors.Errorf("validate source card failed: %w", err)
	}

	// Запрещаем переводить деньги с чужой карты
	if cardFrom.ClientID != client.ID {
		return nil, errors.Errorf("insufficient privileges: someone card %s", fromCard)
	}

	// Проверяем карту приёмник
	cardTo, err := s.card.ValidateCardByNum(ctx, toCard)
	if err != nil {
		return nil, errors.Errorf("validate destination card failed: %w", err)
	}

	// Переводим деньги
	transactions, err := s.repo.TransferBetweenAccounts(ctx, cardFrom.AccountID, cardTo.AccountID, amount)
	if err != nil {
		return nil, errors.Errorf("card to card transaction failed: %w", err)
	}

	return transactions, nil
}

// Deposit пополнение счёта.
func (s *transactionService) Deposit(ctx context.Context, cardNum string, amount float64) (*models.Transaction, error) {
	// Получаем текущего клиента
	client, err := models.ClientFromContext(ctx)
	if err != nil {
		return nil, errors.New("authenticate need")
	}

	// Проверяем карту источник
	card, err := s.card.ValidateCardByNum(ctx, cardNum)
	if err != nil {
		return nil, errors.Errorf("validate source card failed: %w", err)
	}

	// Запрещаем пополнять чужие карты
	if card.ClientID != client.ID {
		return nil, errors.Errorf("insufficient privileges: someone card %s", cardNum)
	}

	transactions, err := s.repo.Deposit(ctx, card.AccountID, amount)
	if err != nil {
		return nil, errors.Errorf("card to card transaction failed: %w", err)
	}

	return transactions, nil
}
