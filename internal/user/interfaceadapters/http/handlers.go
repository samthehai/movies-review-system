package http

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/samthehai/ml-backend-test-samthehai/config"
	handlersusecase "github.com/samthehai/ml-backend-test-samthehai/internal/user/interfaceadapters/http/usecase"
	"github.com/samthehai/ml-backend-test-samthehai/internal/user/usecase"
	"github.com/samthehai/ml-backend-test-samthehai/pkg/httperrors"
	"github.com/samthehai/ml-backend-test-samthehai/pkg/logger"
	"github.com/samthehai/ml-backend-test-samthehai/pkg/utils"
)

type userHandlers struct {
	cfg         *config.Config
	userUsecase handlersusecase.UserUsecase
	logger      logger.Logger
}

func NewUserHandlers(cfg *config.Config, userUsecase handlersusecase.UserUsecase, log logger.Logger) *userHandlers {
	return &userHandlers{cfg: cfg, userUsecase: userUsecase, logger: log}
}

type registerRequest struct {
	Username string `json:"username" validate:"required,lte=254,gte=3"`
	Email    string `json:"email" validate:"required,lte=254,email"`
	Password string `json:"password" validate:"required,gte=6"`
}

type registerResponse struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

func (h *userHandlers) Register() echo.HandlerFunc {
	return func(c echo.Context) error {
		req := &registerRequest{}
		if err := utils.ReadRequest(c, req); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httperrors.ErrorResponse(err))
		}

		ctx := utils.GetRequestCtx(c)
		createdUser, err := h.userUsecase.Register(ctx, usecase.RegisterParams{
			Username: req.Username,
			Email:    req.Email,
			Password: req.Password,
		})
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httperrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusCreated, registerResponse{
			Username: createdUser.Username,
			Email:    createdUser.Email,
		})
	}
}

type loginRequest struct {
	Email    string `json:"email" validate:"required,lte=254,email"`
	Password string `json:"password" validate:"required,gte=6"`
}

type loginResponse struct {
	Username    string `json:"username"`
	Email       string `json:"email"`
	AccessToken string `json:"access_token"`
}

func (h *userHandlers) Login() echo.HandlerFunc {
	return func(c echo.Context) error {
		req := &loginRequest{}
		if err := utils.ReadRequest(c, req); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httperrors.ErrorResponse(err))
		}

		ctx := utils.GetRequestCtx(c)
		res, err := h.userUsecase.Login(ctx, usecase.LoginParams{
			Email:    req.Email,
			Password: req.Password,
		})
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httperrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusCreated, loginResponse{
			Username:    res.Username,
			Email:       res.Email,
			AccessToken: res.AccessToken,
		})
	}
}
