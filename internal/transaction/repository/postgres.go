package repository

import (
	"bank/internal/models"
	"bank/pkg/errors"
	"context"
	"time"
)

// TransferBetweenAccounts переводит деньги с одного счёта на другой.
func (r *repository) TransferBetweenAccounts(ctx context.Context, accountFrom uint, accountTo uint, amount float64) ([]*models.Transaction, error) {
	now := time.Now()

	// Списываем деньги источника
	transactionFrom, err := r.NewTransaction(ctx, accountFrom, amount*-1, models.TransactionTransfer, now)
	if err != nil {
		return nil, errors.Errorf("failed to initiate source transaction: %w", err)
	}

	// Добавляем деньги приёмнику
	transactionTo, err := r.NewTransaction(ctx, accountTo, amount, models.TransactionTransfer, now)
	if err != nil {
		_ = r.RollbackTransaction(ctx, transactionFrom)
		return nil, errors.Errorf("failed to initiate destination transaction: %w", err)
	}

	if err = r.UpdateTransactionStatus(ctx, transactionFrom, models.TransactionStatusCompleted); err != nil {
		_ = r.RollbackTransaction(ctx, transactionFrom)
		_ = r.RollbackTransaction(ctx, transactionTo)
	}

	if err = r.UpdateTransactionStatus(ctx, transactionTo, models.TransactionStatusCompleted); err != nil {
		_ = r.RollbackTransaction(ctx, transactionFrom)
		_ = r.RollbackTransaction(ctx, transactionTo)
	}

	return []*models.Transaction{transactionFrom, transactionTo}, err
}

// Deposit пополнение счёт.
func (r *repository) Deposit(ctx context.Context, account uint, amount float64) (*models.Transaction, error) {
	now := time.Now()

	// Добавляем денег на счёт
	transaction, err := r.NewTransaction(ctx, account, amount, models.TransactionDeposit, now)
	if err != nil {
		return nil, errors.Errorf("failed to initiate transaction: %w", err)
	}

	if err = r.UpdateTransactionStatus(ctx, transaction, models.TransactionStatusCompleted); err != nil {
		_ = r.RollbackTransaction(ctx, transaction)
	}

	return transaction, err
}

// NewTransaction создаёт транзакцию и изменяет баланс исходного счёта.
func (r *repository) NewTransaction(ctx context.Context, accountId uint, amount float64, transactionType models.TransactionType, now time.Time) (*models.Transaction, error) {
	transaction := &models.Transaction{
		AccountID: accountId,
		Amount:    amount,
		Type:      transactionType,
		Status:    models.TransactionStatusPending,
		CreatedAt: now,
	}

	tr, err := r.db.Begin()
	defer func() {
		_ = tr.Rollback()
	}()
	if err != nil {
		return nil, errors.Errorf("failed to initialize storage transaction: %w", err)
	}
	if err := r.db.QueryRowContext(ctx, "insert into transactions(account_id, amount, transaction_type, transaction_status, created_at) values($1, $2, $3, $4, $5) returning id",
		transaction.AccountID, transaction.Amount, transaction.Type, transaction.Status, now).Scan(&transaction.ID); err != nil {
		return nil, errors.Errorf("failed to create transaction: %w", err)
	}

	var balance float64
	if err := r.db.QueryRowContext(ctx, "update accounts set balance = balance + $1 where id = $2 returning balance", amount, accountId).Scan(&balance); err != nil {
		return nil, errors.Errorf("failed to modify account balance: %w", err)
	}

	if balance < 0 {
		return nil, errors.New("no money on account")
	}

	if err = tr.Commit(); err != nil {
		return nil, errors.Errorf("failed to commit transaction: %w", err)
	}

	return transaction, nil
}

// UpdateTransactionStatus изменяет статус транзакции.
func (r *repository) UpdateTransactionStatus(ctx context.Context, transaction *models.Transaction, status models.TransactionStatus) error {
	_, err := r.db.ExecContext(ctx, "update transactions set transaction_status = $1 where id = $2", status, transaction.ID)
	if err != nil {
		return errors.Errorf("failed update transaction: %w", err)
	}
	transaction.Status = status

	return nil
}

// RollbackTransaction откатывает транзакцию и восстанавливает счёт.
func (r *repository) RollbackTransaction(ctx context.Context, transaction *models.Transaction) error {
	status := models.TransactionStatusFailed
	_, err := r.db.ExecContext(ctx, "update transactions set transaction_status = $1, balance = balance - $2 where id = $2", status, transaction.ID)
	if err != nil {
		return errors.Errorf("failed update transaction: %w", err)
	}
	transaction.Status = status

	return nil
}
