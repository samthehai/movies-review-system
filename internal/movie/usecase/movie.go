package usecase

import (
	"context"
	"fmt"

	"github.com/samthehai/ml-backend-test-samthehai/config"
	"github.com/samthehai/ml-backend-test-samthehai/internal/entity"
	"github.com/samthehai/ml-backend-test-samthehai/internal/movie/usecase/repository"
	"github.com/samthehai/ml-backend-test-samthehai/pkg/httperrors"
	"github.com/samthehai/ml-backend-test-samthehai/pkg/logger"
)

type movieUsecase struct {
	cfg             config.Config
	movieRepository repository.MovieRepository
	logger          logger.Logger
}

func NewMovieUsecase(cfg config.Config, movieRepository repository.MovieRepository, log logger.Logger) *movieUsecase {
	return &movieUsecase{cfg: cfg, movieRepository: movieRepository, logger: log}
}

func (u *movieUsecase) GetMovieByID(ctx context.Context, movieID uint64) (*entity.Movie, error) {
	movie, err := u.movieRepository.FindByID(ctx, movieID)
	if err != nil {
		return nil, httperrors.NewInternalServerError(fmt.Errorf("movieUsecase.GetMovieByID.movieRepository.FindByID: %w", err))
	}

	if movie == nil {
		return nil, httperrors.NewNotFoundError(fmt.Errorf("movieUsecase.GetMovieByID.movieRepository.FindByID: %w", err))
	}

	return movie, nil
}

func (u *movieUsecase) SearchByKeyword(ctx context.Context, keyword string) ([]*entity.Movie, error) {
	if len(keyword) == 0 {
		// TODO:
		return nil, nil
	}

	movies, err := u.movieRepository.FindByKeyword(ctx, keyword)
	if err != nil {
		return nil, httperrors.NewInternalServerError(fmt.Errorf("movieRepository.FindByKeyword: %w", err))
	}

	return movies, nil
}
