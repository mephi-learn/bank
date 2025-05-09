package service

import (
	"bank/internal/models"
	"bank/pkg/errors"
	"context"
	"crypto/rand"
	"fmt"
	"time"
)

const createCardAttempts = 10

// CreateCard создание карты.
func (s *cardService) CreateCard(ctx context.Context, accountID uint) (*models.Card, error) {
	// Получаем счёт
	account, err := s.account.AccountByID(ctx, accountID)
	if err != nil {
		s.log.ErrorContext(ctx, "error search account", "error", err, "account id", accountID)
		return nil, errors.Errorf("failed to search account: %w", err)
	}

	// Получаем текущего клиента
	client, err := models.ClientFromContext(ctx)
	if err != nil {
		s.log.ErrorContext(ctx, "unknown client", "account id", accountID)
		return nil, errors.New("unknown client")
	}

	// Идентификатор клиента счёта должен совпадать с идентификатором клиента в счёте
	if client.ID != account.ClientID {
		s.log.ErrorContext(ctx, "create card incorrect client", "account id", accountID, "client id", account.ClientID)
		return nil, errors.Errorf("create card incorrect client: %w", err)
	}

	return s.createCard(ctx, account)
}

// CreateCardByAdmin создание карты.
func (s *cardService) CreateCardByAdmin(ctx context.Context, accountID uint) (*models.Card, error) {
	// Получаем счёт
	account, err := s.account.AccountByID(ctx, accountID)
	if err != nil {
		s.log.ErrorContext(ctx, "error search account", "error", err, "account id", accountID)
		return nil, errors.Errorf("failed to search account: %w", err)
	}

	return s.createCard(ctx, account)
}

// createCard создание карты.
func (s *cardService) createCard(ctx context.Context, account *models.Account) (*models.Card, error) {
	// Пробуем создать номер карты, но чтобы не зависнуть, например, в случае недоступности БД, ограничиваем количество попыток
	attempt := 0
	var cardNumber string
	for cardNumber == "" && attempt < createCardAttempts {
		newNumber := createCardNumber()
		if card, err := s.repo.CardExistsByNumber(ctx, newNumber); card == nil && err == nil {
			cardNumber = newNumber
		}
		attempt++
	}
	if cardNumber == "" {
		s.log.ErrorContext(ctx, "error create card number", "attempts", attempt, "client id", account.ClientID)
		return nil, errors.Errorf("failed to create card number")
	}

	expireDate := s.generateExpirationDate()

	cvv, err := s.generateCVV()
	if err != nil {
		return nil, fmt.Errorf("failed create CVV: %w", err)
	}

	// Формируем данные для создания
	cardRepo := models.Card{
		ClientID:  account.ClientID,
		AccountID: account.ID,
		Number:    cardNumber,
		Expire:    expireDate,
		CVV:       cvv,
	}

	// Создаём карту в хранилище
	card, err := s.repo.CreateCard(ctx, &cardRepo)
	if err != nil {
		s.log.ErrorContext(ctx, "error create card", "error", err, "account id", account.ID, "client id", account.ClientID)
		return nil, errors.Errorf("cannot create card: %w", err)
	}

	return card, nil
}

// generateExpirationDate генерирует дату истечения срока действия карты.
func (s *cardService) generateExpirationDate() time.Time {
	now := time.Now()
	expiryDate := now.AddDate(3, 0, 0) // Карта действительна 3 года

	return expiryDate
}

// generateCVV генерирует случайный CVV код.
func (s *cardService) generateCVV() (string, error) {
	// Генерируем случайные 3 цифры
	b := make([]byte, 2)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	// Конвертируем в число от 100 до 999
	cvv := 100 + (int(b[0])<<8|int(b[1]))%900

	return fmt.Sprintf("%03d", cvv), nil
}
