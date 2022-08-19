package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/samthehai/ml-backend-test-samthehai/internal/movie/usecase/repository"
)

type favoriteRepository struct {
	connManager ConnManager
}

func NewFavoriteRepository(connManager ConnManager) *favoriteRepository {
	return &favoriteRepository{connManager: connManager}
}

const addFavoriteMovieQuery = `INSERT INTO favorites(user_id, movie_id) VALUES (?,?)`

func (r *favoriteRepository) AddFavoriteMovie(ctx context.Context, args repository.AddFavoriteMovieParams) error {
	_, err := r.connManager.GetWriter().ExecContext(ctx, addFavoriteMovieQuery, args.UserID, args.MovieID)
	if err != nil {
		return fmt.Errorf("ExecContext: %w", err)
	}

	return nil
}

const checkIsFavoriteMovieQuery = `SELECT user_id, movie_id, created_at, updated_at FROM favorites WHERE user_id = ? AND movie_id = ?`

func (r *favoriteRepository) CheckIsFavoriteMovie(ctx context.Context, args repository.CheckIsFavoriteMovieParams) (bool, error) {
	favorite := &Favorite{}
	if err := r.connManager.GetReader().QueryRowxContext(ctx, checkIsFavoriteMovieQuery, args.UserID,
		args.MovieID).StructScan(favorite); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("QueryRowxContext: %w", err)
	}

	return true, nil
}
