package usecase

import (
	"context"

	"github.com/samthehai/ml-backend-test-samthehai/internal/entity"
	"github.com/samthehai/ml-backend-test-samthehai/internal/user/usecase"
)

type UserUsecase interface {
	Register(ctx context.Context, args usecase.RegisterParams) (*entity.User, error)
	Login(ctx context.Context, args usecase.LoginParams) (*entity.UserWithAccessToken, error)
}
