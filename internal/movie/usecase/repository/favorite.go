//go:generate mockgen -source favorite.go -destination ../testdata/mock_repository/favorite_gen.go
package repository

import (
	"context"

	"github.com/samthehai/ml-backend-test-samthehai/internal/entity"
)

type AddFavoriteMovieParams struct {
	UserID  uint64 `json:"user_id"`
	MovieID uint64 `json:"email"`
}

type CheckIsFavoriteMovieParams struct {
	UserID  uint64 `json:"user_id"`
	MovieID uint64 `json:"email"`
}

type FavoriteRepository interface {
	AddFavoriteMovie(ctx context.Context, args AddFavoriteMovieParams) error
	CheckIsFavoriteMovie(ctx context.Context, args CheckIsFavoriteMovieParams) (bool, error)
	FindFavoriteMoviesByUserID(ctx context.Context, userID uint64) ([]*entity.Movie, error)
}
