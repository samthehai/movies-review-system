package repository

import "context"

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
}
