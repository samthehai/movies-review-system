package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/samthehai/ml-backend-test-samthehai/config"
	"github.com/samthehai/ml-backend-test-samthehai/internal/entity"
	handlersusecase "github.com/samthehai/ml-backend-test-samthehai/internal/movie/interfaceadapters/http/usecase"
	"github.com/samthehai/ml-backend-test-samthehai/internal/movie/usecase"
	"github.com/samthehai/ml-backend-test-samthehai/pkg/httperrors"
	"github.com/samthehai/ml-backend-test-samthehai/pkg/logger"
	"github.com/samthehai/ml-backend-test-samthehai/pkg/utils"
)

type movieHandlers struct {
	cfg              *config.Config
	movieUsecase     handlersusecase.MovieUsecase
	logger           logger.Logger
	getCurrentUserFn func(c echo.Context) (*entity.User, error)
}

func NewMovieHandlers(cfg *config.Config, movieUsecase handlersusecase.MovieUsecase,
	log logger.Logger, getCurrentUserFn func(c echo.Context) (*entity.User, error)) *movieHandlers {
	return &movieHandlers{cfg: cfg, movieUsecase: movieUsecase, logger: log,
		getCurrentUserFn: getCurrentUserFn}
}

type getByIDRequest struct {
	ID uint64 `param:"id"`
}

func (h *movieHandlers) GetByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		req := &getByIDRequest{}
		if err := utils.ReadRequest(c, req); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httperrors.ErrorResponse(err))
		}

		ctx := utils.GetRequestCtx(c)
		movie, err := h.movieUsecase.GetMovieByID(ctx, req.ID)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httperrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, movie)
	}
}

type searchByKeywordRequest struct {
	Keyword string `query:"search"`
}

func (h *movieHandlers) SearchByKeyword() echo.HandlerFunc {
	return func(c echo.Context) error {
		req := &searchByKeywordRequest{}
		if err := utils.ReadRequest(c, req); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httperrors.ErrorResponse(err))
		}

		ctx := utils.GetRequestCtx(c)
		movies, err := h.movieUsecase.SearchByKeyword(ctx, req.Keyword)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httperrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, movies)
	}
}

type addFavoriteMovieRequest struct {
	MovieID uint64 `param:"id"`
}

func (h *movieHandlers) AddFavoriteMovie() echo.HandlerFunc {
	return func(c echo.Context) error {
		req := &addFavoriteMovieRequest{}
		if err := utils.ReadRequest(c, req); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httperrors.ErrorResponse(err))
		}

		currentUser, err := h.getCurrentUserFn(c)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httperrors.ErrorResponse(err))
		}

		ctx := utils.GetRequestCtx(c)
		if err := h.movieUsecase.AddFavoriteMovie(ctx, usecase.AddFavoriteMovieParams{
			UserID:  currentUser.ID,
			MovieID: req.MovieID,
		}); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httperrors.ErrorResponse(err))
		}

		return c.NoContent(http.StatusOK)
	}
}
