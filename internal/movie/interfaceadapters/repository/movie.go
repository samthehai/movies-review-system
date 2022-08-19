package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/samthehai/ml-backend-test-samthehai/internal/entity"
)

type movieRepository struct {
	connManager ConnManager
}

func NewMovieRepository(connManager ConnManager) *movieRepository {
	return &movieRepository{connManager: connManager}
}

const findByID = `SELECT id, original_title, original_language, overview, poster_path, backdrop_path,
adult, release_date, budget, revenue, created_at, updated_at FROM movies WHERE id = ?`

func (r *movieRepository) FindByID(ctx context.Context, movieID uint64) (*entity.Movie, error) {
	foundMovie := &Movie{}
	if err := r.connManager.GetReader().QueryRowxContext(ctx, findByID, movieID).StructScan(foundMovie); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("userRepository.FindByID.QueryRowxContext: %w", err)
	}

	return &entity.Movie{
		ID:               foundMovie.ID,
		OriginalTitle:    foundMovie.OriginalTitle,
		OriginalLanguage: foundMovie.OriginalLanguage,
		Overview:         foundMovie.Overview,
		PosterPath:       foundMovie.PosterPath,
		BackdropPath:     foundMovie.BackdropPath,
		Adult:            foundMovie.Adult,
		ReleaseDate:      foundMovie.ReleaseDate,
		Budget:           foundMovie.Budget,
		Revenue:          foundMovie.Revenue,
		CreatedAt:        foundMovie.CreatedAt,
		UpdatedAt:        foundMovie.UpdatedAt,
	}, nil
}

const findByKeyword = `SELECT id, original_title, original_language, overview, poster_path, backdrop_path,
adult, release_date, budget, revenue, created_at, updated_at
FROM movies
WHERE MATCH (original_title, overview, original_language) AGAINST ('%s*' IN BOOLEAN MODE)
ORDER BY id ASC`

func (r *movieRepository) FindByKeyword(ctx context.Context, keyword string) ([]*entity.Movie, error) {
	query := fmt.Sprintf(findByKeyword, keyword)
	movies := make([]*entity.Movie, 0)

	rows, err := r.connManager.GetReader().QueryxContext(ctx, query)
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

const findPopularMovies = `SELECT id, original_title, original_language, overview, poster_path, backdrop_path,
adult, release_date, budget, revenue, created_at, updated_at, IFNULL(favorite_numbers.favorite_number, 0) as favorite_number
FROM movies
LEFT JOIN (SELECT movie_id, count(*) AS favorite_number FROM favorites GROUP BY movie_id) AS favorite_numbers
ON movies.id = favorite_numbers.movie_id
ORDER BY favorite_number DESC
LIMIT %v`

func (r *movieRepository) FindPopularMovies(ctx context.Context, limit uint) ([]*entity.Movie, error) {
	query := fmt.Sprintf(findPopularMovies, limit)
	movies := make([]*entity.Movie, 0)

	rows, err := r.connManager.GetReader().QueryxContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("QueryxContext: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		movie := &struct {
			*Movie
			FavoriteNumber uint `json:"favorite_number" db:"favorite_number"`
		}{}
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
