package repository

import (
	"context"

	"github.com/samthehai/ml-backend-test-samthehai/internal/entity"
)

type MovieRepository interface {
	FindByID(ctx context.Context, movieID uint64) (*entity.Movie, error)
	FindByKeyword(ctx context.Context, keyword string) ([]*entity.Movie, error)
}
