package usecase

import (
	"context"

	"github.com/samthehai/ml-backend-test-samthehai/internal/entity"
)

type MovieUsecase interface {
	GetMovieByID(ctx context.Context, movieID uint64) (*entity.Movie, error)
	SearchByKeyword(ctx context.Context, keyword string) ([]*entity.Movie, error)
}
