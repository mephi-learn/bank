package repository

import (
	"bank/internal/models"
	"bank/pkg/errors"
	"context"
)

func (r *repository) Create(ctx context.Context, clientData *models.Client) (uint, error) {
	if clientData.Admin {
		var userCount uint
		if err := r.db.QueryRowContext(ctx, "select count(*) cn from clients where admin=$1", true).Scan(&userCount); err != nil {
			return 0, errors.Errorf("failed to check exist user: %w", err)
		}
		if userCount > 0 {
			return 0, errors.Errorf("only one user can be admin")
		}
	}

	var userId uint
	if err := r.db.QueryRowContext(ctx, "insert into clients(name, username, password, email, admin) values($1, $2, $3, $4, $5) returning id",
		clientData.Name, clientData.Username, clientData.Password, clientData.Email, clientData.Admin).Scan(&userId); err != nil {
		return 0, errors.Errorf("failed to create user: %w", err)
	}

	return userId, nil
}

func (r *repository) GetByUsernameAndPassword(ctx context.Context, clientData *models.Client) (*models.Client, error) {
	client := models.Client{}
	if err := r.db.QueryRowContext(ctx, "select id, name, username, email, admin from clients where username=$1 and password=$2", clientData.Username, clientData.Password).
		Scan(&client.ID, &client.Name, &client.Username, &client.Email, &client.Admin); err != nil {
		return nil, errors.Errorf("client not found: %w", err)
	}

	return &client, nil
}

func (r *repository) GetById(ctx context.Context, clientId uint) (*models.Client, error) {
	client := models.Client{}
	if err := r.db.QueryRowContext(ctx, "select id, name, username, email, admin from clients where id=$1", clientId).
		Scan(&client.ID, &client.Name, &client.Username, &client.Email, &client.Admin); err != nil {
		return nil, errors.Errorf("client not found: %w", err)
	}

	return &client, nil
}

func (r *repository) List(ctx context.Context) ([]*models.Client, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, name, username, email, admin FROM clients")
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	var clients []*models.Client
	for rows.Next() {
		var client models.Client
		if err = rows.Scan(&client.ID, &client.Name, &client.Username, &client.Email, &client.Admin); err != nil {
			return clients, err
		}
		clients = append(clients, &client)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return clients, nil
}
