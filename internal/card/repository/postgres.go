package repository

import (
	"bank/internal/models"
	"bank/pkg/errors"
	"context"
	"database/sql"
	"time"
)

func (r *repository) CreateCard(ctx context.Context, card *models.Card) (*models.Card, error) {
	// Без PGP ключа работа невозможна
	if r.config.PGPkey == "" {
		r.log.ErrorContext(ctx, "missing PGP key")
		return nil, errors.New("missing PGP key")
	}

	cvv, err := card.CryptCVV()
	if err != nil {
		return nil, errors.Errorf("failed to crypt cvv: %w", err)
	}

	var cardID uint
	if err := r.db.QueryRowContext(ctx, "insert into cards(client_id, account_id, cardnum, expire, cvv) values($1, $2, pgp_sym_encrypt($3, $4), $5, $6) returning id",
		card.ClientID, card.AccountID, card.Number, r.config.PGPkey, card.Expire, cvv).Scan(&cardID); err != nil {
		return nil, errors.Errorf("failed to create card: %w", err)
	}

	card = &models.Card{
		ID:        cardID,
		ClientID:  card.ClientID,
		AccountID: card.AccountID,
		Number:    card.Number,
		Expire:    time.Time{},
		CVV:       cvv,
	}

	return card, nil
}

func (r *repository) DeleteCard(ctx context.Context, cardID uint) error {
	_, err := r.db.ExecContext(ctx, "delete from cards where id = $1", cardID)
	if err != nil {
		return errors.Errorf("failed to delete card: %w", err)
	}

	return nil
}

func (r *repository) ListCards(ctx context.Context, clientID uint) ([]*models.Card, error) {
	// Без PGP ключа работа невозможна
	if r.config.PGPkey == "" {
		r.log.ErrorContext(ctx, "missing PGP key")
		return nil, errors.New("missing PGP key")
	}

	rows, err := r.db.QueryContext(ctx, "SELECT id, client_id, account_id, pgp_sym_decrypt(cardnum, $1), expire, cvv FROM cards WHERE client_id = $2", r.config.PGPkey, clientID)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	var cards []*models.Card
	for rows.Next() {
		var card models.Card
		var cvvCrypt string
		if err = rows.Scan(&card.ID, &card.ClientID, &card.AccountID, &card.Number, &card.Expire, &cvvCrypt); err != nil {
			return cards, err
		}

		if err = card.DecryptCVV(cvvCrypt); err != nil {
			return nil, errors.Errorf("failed to decrypt cvv number: %w", err)
		}

		cards = append(cards, &card)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return cards, nil
}

func (r *repository) CardById(ctx context.Context, cardID uint) (*models.Card, error) {
	// Без PGP ключа работа невозможна
	if r.config.PGPkey == "" {
		r.log.ErrorContext(ctx, "missing PGP key")
		return nil, errors.New("missing PGP key")
	}

	var cvvCrypt string
	card := models.Card{}
	if err := r.db.QueryRowContext(ctx, "SELECT id, client_id, account_id, pgp_sym_decrypt(cardnum, $1), expire, cvv FROM cards WHERE id = $2", r.config.PGPkey, cardID).
		Scan(&card.ID, &card.ClientID, &card.AccountID, &card.Number, &card.Expire, &cvvCrypt); err != nil {
		return nil, errors.Errorf("failed to get account info by id: %w", err)
	}

	if err := card.DecryptCVV(cvvCrypt); err != nil {
		return nil, errors.Errorf("failed to decrypt cvv number: %w", err)
	}

	return &card, nil
}

func (r *repository) CardByNumber(ctx context.Context, cardNumber string) (*models.Card, error) {
	// Без PGP ключа работа невозможна
	if r.config.PGPkey == "" {
		r.log.ErrorContext(ctx, "missing PGP key")
		return nil, errors.New("missing PGP key")
	}

	card := models.Card{Number: cardNumber}

	var cvvCrypt string
	if err := r.db.QueryRowContext(ctx, "SELECT id, client_id, account_id, pgp_sym_decrypt(cardnum, $1) cardnum, expire, cvv FROM cards WHERE pgp_sym_decrypt(cardnum, $2) = $3", r.config.PGPkey, r.config.PGPkey, card.Number).
		Scan(&card.ID, &card.ClientID, &card.AccountID, &card.Number, &card.Expire, &cvvCrypt); err != nil {
		return nil, errors.Errorf("failed to get account info by number: %w", err)
	}

	return &card, nil
}

func (r *repository) CardExistsByNumber(ctx context.Context, cardNumber string) (*models.Card, error) {
	card, err := r.CardByNumber(ctx, cardNumber)
	if err == nil {
		return card, nil
	}

	if !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	return nil, nil
}
