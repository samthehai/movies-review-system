package usecase

import (
	"context"
	"fmt"
	"net/http"

	"github.com/samthehai/ml-backend-test-samthehai/config"
	"github.com/samthehai/ml-backend-test-samthehai/internal/entity"
	"github.com/samthehai/ml-backend-test-samthehai/internal/movie/usecase/repository"
	"github.com/samthehai/ml-backend-test-samthehai/pkg/httperrors"
	"github.com/samthehai/ml-backend-test-samthehai/pkg/logger"
)

const limitPopularMovieNumber = 100

type movieUsecase struct {
	cfg                config.Config
	movieRepository    repository.MovieRepository
	favoriteRepository repository.FavoriteRepository
	logger             logger.Logger
}

func NewMovieUsecase(cfg config.Config, log logger.Logger, movieRepository repository.MovieRepository, favoriteRepository repository.FavoriteRepository) *movieUsecase {
	return &movieUsecase{cfg: cfg, logger: log, movieRepository: movieRepository, favoriteRepository: favoriteRepository}
}

func (u *movieUsecase) GetMovieByID(ctx context.Context, movieID uint64) (*entity.Movie, error) {
	movie, err := u.movieRepository.FindByID(ctx, movieID)
	if err != nil {
		return nil, httperrors.NewInternalServerError(fmt.Errorf("movieRepository.FindByID: %w", err))
	}

	if movie == nil {
		return nil, httperrors.NewNotFoundError(fmt.Errorf("movieRepository.FindByID: not found"))
	}

	return movie, nil
}

func (u *movieUsecase) SearchByKeyword(ctx context.Context, keyword string) ([]*entity.Movie, error) {
	if len(keyword) == 0 {
		movies, err := u.movieRepository.FindPopularMovies(ctx, limitPopularMovieNumber)
		if err != nil {
			return nil, httperrors.NewInternalServerError(fmt.Errorf("movieRepository.FindPopularMovies: %w", err))
		}
		return movies, nil
	}

	movies, err := u.movieRepository.FindByKeyword(ctx, keyword)
	if err != nil {
		return nil, httperrors.NewInternalServerError(fmt.Errorf("movieRepository.FindByKeyword: %w", err))
	}

	return movies, nil
}

type AddFavoriteMovieParams struct {
	UserID  uint64 `json:"user_id"`
	MovieID uint64 `json:"email"`
}

func (u *movieUsecase) AddFavoriteMovie(ctx context.Context, args AddFavoriteMovieParams) error {
	movie, err := u.movieRepository.FindByID(ctx, args.MovieID)
	if err != nil {
		return httperrors.NewInternalServerError(fmt.Errorf("movieRepository.FindByID: %w", err))
	}

	if movie == nil {
		return httperrors.NewNotFoundError(fmt.Errorf("movieRepository.FindByID: not found"))
	}

	isFavorite, err := u.favoriteRepository.CheckIsFavoriteMovie(ctx, repository.CheckIsFavoriteMovieParams{
		UserID:  args.UserID,
		MovieID: args.MovieID,
	})
	if err != nil {
		return httperrors.NewInternalServerError(fmt.Errorf("favoriteRepository.CheckIsFavoriteMovie: %w", err))
	}

	if isFavorite {
		return httperrors.NewRestError(http.StatusBadRequest, "already is favorited", nil)
	}

	if err := u.favoriteRepository.AddFavoriteMovie(ctx, repository.AddFavoriteMovieParams{
		UserID:  args.UserID,
		MovieID: args.MovieID,
	}); err != nil {
		return httperrors.NewRestError(http.StatusInternalServerError,
			fmt.Errorf("favoriteRepository.AddFavoriteMovie: %w", err).Error(), nil)
	}

	return nil
}

func (u *movieUsecase) ListFavoriteMoviesByUserID(ctx context.Context, userID uint64) ([]*entity.Movie, error) {
	movies, err := u.favoriteRepository.FindFavoriteMoviesByUserID(ctx, userID)
	if err != nil {
		return nil, httperrors.NewInternalServerError(fmt.Errorf("favoriteRepository.FindFavoriteMoviesByUserID: %w", err))
	}

	return movies, nil
}
