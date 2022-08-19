package middlewares

import (
	"context"

	"github.com/samthehai/ml-backend-test-samthehai/config"
	"github.com/samthehai/ml-backend-test-samthehai/internal/entity"
	"github.com/samthehai/ml-backend-test-samthehai/pkg/logger"
)

type UserUsecase interface {
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
}

type middlewareManager struct {
	cfg         *config.Config
	logger      logger.Logger
	userUsecase UserUsecase
}

func NewMiddlewareManager(cfg *config.Config, logger logger.Logger, userUsecase UserUsecase) *middlewareManager {
	return &middlewareManager{cfg: cfg, logger: logger, userUsecase: userUsecase}
}
