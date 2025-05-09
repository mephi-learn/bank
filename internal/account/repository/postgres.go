package repository

import (
	"bank/internal/models"
	"bank/pkg/errors"
	"context"
	"database/sql"
)

func (r *repository) CreateAccount(ctx context.Context, account *models.Account) (*models.Account, error) {
	var accountId uint
	if err := r.db.QueryRowContext(ctx, "insert into accounts(client_id, balance, number) values($1, $2, $3) returning id",
		account.ClientID, 0, account.Number).Scan(&accountId); err != nil {
		return nil, errors.Errorf("failed to create account: %w", err)
	}

	account = &models.Account{
		ID:       accountId,
		ClientID: account.ClientID,
		Number:   account.Number,
		Balance:  0,
	}

	return account, nil
}

func (r *repository) DeleteAccount(ctx context.Context, accountId uint) error {
	_, err := r.db.ExecContext(ctx, "delete from accounts where id = $1", accountId)
	if err != nil {
		return errors.Errorf("failed to delete account: %w", err)
	}

	return nil
}

func (r *repository) ListAccounts(ctx context.Context, clientId uint) ([]*models.Account, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, client_id, number, balance FROM accounts WHERE client_id = $1", clientId)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	var accounts []*models.Account
	for rows.Next() {
		var account models.Account
		if err := rows.Scan(&account.ID, &account.ClientID, &account.Number, &account.Balance); err != nil {
			return accounts, err
		}
		accounts = append(accounts, &account)
	}
	if err = rows.Err(); err != nil {
		return accounts, err
	}

	return accounts, nil
}

func (r *repository) AccountByID(ctx context.Context, accountId uint) (*models.Account, error) {
	account := models.Account{}
	if err := r.db.QueryRowContext(ctx, "SELECT id, client_id, number, balance FROM accounts WHERE id = $1", accountId).Scan(&account.ID, &account.ClientID, &account.Number, &account.Balance); err != nil {
		return nil, errors.Errorf("failed to get account info by id: %w", err)
	}

	return &account, nil
}

func (r *repository) AccountByNumber(ctx context.Context, accountNumber string) (*models.Account, error) {
	account := models.Account{}
	if err := r.db.QueryRowContext(ctx, "SELECT id, client_id, number, balance FROM accounts WHERE number = $1", accountNumber).Scan(&account.ID, &account.ClientID, &account.Number, &account.Balance); err != nil {
		return nil, errors.Errorf("failed to get account info by number: %w", err)
	}

	return &account, nil
}

func (r *repository) AccountExistsByNumber(ctx context.Context, accountNumber string) (*models.Account, error) {
	account, err := r.AccountByNumber(ctx, accountNumber)
	if err == nil {
		return account, nil
	}

	if !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	return nil, nil
}
