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

// Register godoc
// @Summary Register new user
// @Description register new user, returns username and email
// @Tags Users
// @Accept json
// @Param  registerRequest body registerRequest true "registerRequest body"
// @Produce json
// @Success 201 {object} registerResponse
// @Failure 400 {object} httperrors.RestError
// @Failure 404 {object} httperrors.RestError
// @Failure 500 {object} httperrors.RestError
// @Router /users/register [post]
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

// Login godoc
// @Summary Login user
// @Description login user, returns user information and accesstoken with default expired time is 15 minutes
// @Tags Users
// @Accept json
// @Param  loginRequest body loginRequest true "loginRequest body"
// @Produce json
// @Success 201 {object} loginResponse
// @Failure 400 {object} httperrors.RestError
// @Failure 404 {object} httperrors.RestError
// @Failure 500 {object} httperrors.RestError
// @Router /users/login [post]
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
