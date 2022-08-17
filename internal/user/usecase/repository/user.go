package repository

import (
	"context"

	"github.com/samthehai/ml-backend-test-samthehai/internal/entity"
)

type RegisterParams struct {
	Username       string `json:"username"`
	Email          string `json:"email"`
	HashedPassword string `json:"hashed_password"`
}

type UserRepository interface {
	Register(ctx context.Context, args RegisterParams) (*entity.User, error)
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
}
