package usecase

import (
	"context"

	"github.com/samthehai/ml-backend-test-samthehai/internal/entity"
	"github.com/samthehai/ml-backend-test-samthehai/internal/movie/usecase"
)

type MovieUsecase interface {
	GetMovieByID(ctx context.Context, movieID uint64) (*entity.Movie, error)
	SearchByKeyword(ctx context.Context, keyword string) ([]*entity.Movie, error)
	AddFavoriteMovie(ctx context.Context, args usecase.AddFavoriteMovieParams) error
	ListFavoriteMoviesByUserID(ctx context.Context, userID uint64) ([]*entity.Movie, error)
}
