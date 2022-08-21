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

// GetByID godoc
// @Summary Get movie details information by its Id
// @Description Get movie details information by its Id, if the id is not exist returns http.StatusNotFound
// @Tags Movies
// @Accept json
// @Param id path uint64 true "id"
// @Produce json
// @Success 200 {object} entity.Movie
// @Failure 400 {object} httperrors.RestError
// @Failure 404 {object} httperrors.RestError
// @Failure 500 {object} httperrors.RestError
// @Router /movies/{id} [get]
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

// SearchByKeyword godoc
// @Summary Search movies by specific keyword. If do not specify keyword will return a list of popular movies.
// @Description Search movies by specific keyword. If do not specify keyword will return a list of popular movies.
// @Tags Movies
// @Accept json
// @Param search query string false "search query"
// @Produce json
// @Success 200 {object} []entity.Movie
// @Failure 400 {object} httperrors.RestError
// @Failure 500 {object} httperrors.RestError
// @Router /movies [get]
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

// AddFavoriteMovie godoc
// @Summary Add movie to user's favorite list.
// @Description Add movie to user's favorite list.
// 							If user is not login returns http.StatusUnauthorized.
// 							If the movie is already favorite returns http.StatusBadRequest.
// @Tags Movies
// @Accept json
// @Param id path uint64 true "id"
// @Param Authorization header string true "Format: Bearer accesstoken - which can be get when call /api/v1/login api"
// @Produce json
// @Security ApiKeyAuth
// @Success 200
// @Failure 400 {object} httperrors.RestError
// @Failure 401 {object} httperrors.RestError
// @Failure 404 {object} httperrors.RestError
// @Failure 500 {object} httperrors.RestError
// @Router /favorites/{id} [post]
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

// ListFavoriteMovies godoc
// @Summary List favorite movies of current login user.
// @Description List favorite movies of current login user.
// 							If user is not login returns http.StatusUnauthorized.
// @Tags Movies
// @Accept json
// @Param Authorization header string true "Format: Bearer accesstoken - which can be get when call /api/v1/login api"
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} []entity.Movie
// @Failure 400 {object} httperrors.RestError
// @Failure 401 {object} httperrors.RestError
// @Failure 500 {object} httperrors.RestError
// @Router /favorites [get]
func (h *movieHandlers) ListFavoriteMovies() echo.HandlerFunc {
	return func(c echo.Context) error {
		currentUser, err := h.getCurrentUserFn(c)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httperrors.ErrorResponse(err))
		}

		ctx := utils.GetRequestCtx(c)
		movies, err := h.movieUsecase.ListFavoriteMoviesByUserID(ctx, currentUser.ID)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httperrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, movies)
	}
}
