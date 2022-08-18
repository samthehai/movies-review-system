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
