package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/samthehai/ml-backend-test-samthehai/config"
	handlersusecase "github.com/samthehai/ml-backend-test-samthehai/internal/movie/interfaceadapters/http/usecase"
	"github.com/samthehai/ml-backend-test-samthehai/pkg/httperrors"
	"github.com/samthehai/ml-backend-test-samthehai/pkg/logger"
	"github.com/samthehai/ml-backend-test-samthehai/pkg/utils"
)

type movieHandlers struct {
	cfg          *config.Config
	movieUsecase handlersusecase.MovieUsecase
	logger       logger.Logger
}

func NewMovieHandlers(cfg *config.Config, movieUsecase handlersusecase.MovieUsecase, log logger.Logger) *movieHandlers {
	return &movieHandlers{cfg: cfg, movieUsecase: movieUsecase, logger: log}
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
