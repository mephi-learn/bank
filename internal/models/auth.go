package models

import (
	"bank/pkg/errors"
	"context"
)

const (
	errContextIsNil = "context is nil"
	errNoUser       = "no user"
)

var ctxUser ctxAuthKey = struct{}{}

type ctxAuthKey struct{}

type Client struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Admin    bool   `json:"admin"`
}

func (u *Client) ClientToContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, ctxUser, u)
}

func ClientFromContext(ctx context.Context) (*Client, error) {
	if ctx == nil {
		return nil, errors.New(errContextIsNil)
	}

	rawUser := ctx.Value(ctxUser)
	if rawUser == nil {
		return nil, errors.New(errNoUser)
	}

	user := rawUser.(*Client)
	if user == nil {
		return nil, errors.New(errNoUser)
	}

	return user, nil
}
