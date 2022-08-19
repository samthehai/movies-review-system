package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/samthehai/ml-backend-test-samthehai/internal/entity"
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

const findFavoriteMoviesByUserIDQuery = `SELECT movies.id, movies.original_title, movies.original_language,
movies.overview, movies.poster_path, movies.backdrop_path,
movies.adult, movies.release_date, movies.budget, movies.revenue, movies.created_at, movies.updated_at
FROM movies
INNER JOIN favorites
ON movies.id = favorites.movie_id
WHERE favorites.user_id = ?
ORDER BY movies.id ASC`

func (r *favoriteRepository) FindFavoriteMoviesByUserID(ctx context.Context, userID uint64) ([]*entity.Movie, error) {
	movies := make([]*entity.Movie, 0)

	rows, err := r.connManager.GetReader().QueryxContext(ctx, findFavoriteMoviesByUserIDQuery, userID)
	if err != nil {
		return nil, fmt.Errorf("QueryxContext: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		movie := &Movie{}
		if err = rows.StructScan(movie); err != nil {
			return nil, fmt.Errorf("StructScan: %w", err)
		}

		movies = append(movies, &entity.Movie{
			ID:               movie.ID,
			OriginalTitle:    movie.OriginalTitle,
			OriginalLanguage: movie.OriginalLanguage,
			Overview:         movie.Overview,
			PosterPath:       movie.PosterPath,
			BackdropPath:     movie.BackdropPath,
			Adult:            movie.Adult,
			ReleaseDate:      movie.ReleaseDate,
			Budget:           movie.Budget,
			Revenue:          movie.Revenue,
			CreatedAt:        movie.CreatedAt,
			UpdatedAt:        movie.UpdatedAt,
		})
	}

	return movies, nil
}
