package middlewares

import (
	"errors"
	"fmt"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/samthehai/ml-backend-test-samthehai/internal/entity"
	"github.com/samthehai/ml-backend-test-samthehai/pkg/httperrors"
	"github.com/samthehai/ml-backend-test-samthehai/pkg/token"
	"github.com/samthehai/ml-backend-test-samthehai/pkg/utils"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
	currentUserKey          = "current_user"
)

func (mw *middlewareManager) AuthMiddleware(tokenMaker token.Maker) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authorizationHeader := c.Request().Header.Get(authorizationHeaderKey)
			if len(authorizationHeader) == 0 {
				err := errors.New("authorization header is not provided")
				utils.LogResponseError(c, mw.logger, err)
				return c.JSON(httperrors.ErrorResponse(httperrors.NewUnauthorizedError(err)))
			}

			fields := strings.Fields(authorizationHeader)
			if len(fields) < 2 {
				err := errors.New("invalid authorization header format")
				utils.LogResponseError(c, mw.logger, err)
				return c.JSON(httperrors.ErrorResponse(httperrors.NewUnauthorizedError(err)))
			}

			authorizationType := strings.ToLower(fields[0])
			if authorizationType != authorizationTypeBearer {
				err := fmt.Errorf("unsupported authorization type %s", authorizationType)
				utils.LogResponseError(c, mw.logger, err)
				return c.JSON(httperrors.ErrorResponse(httperrors.NewUnauthorizedError(err)))
			}

			accessToken := fields[1]
			payload, err := tokenMaker.VerifyToken(accessToken)
			if err != nil {
				utils.LogResponseError(c, mw.logger, err)
				return c.JSON(httperrors.ErrorResponse(httperrors.NewUnauthorizedError(err)))
			}

			user, err := mw.userUsecase.GetUserByEmail(c.Request().Context(), payload.Email)
			if err != nil || user == nil {
				utils.LogResponseError(c, mw.logger, err)
				return c.JSON(httperrors.ErrorResponse(httperrors.NewUnauthorizedError(httperrors.ErrNotFound.Error())))
			}

			c.Set(currentUserKey, user)
			c.Set(authorizationPayloadKey, payload)

			return next(c)
		}
	}
}

func (mw *middlewareManager) GetCurrentUser(c echo.Context) (*entity.User, error) {
	user, ok := c.Get(currentUserKey).(*entity.User)
	if !ok || user == nil {
		utils.LogResponseError(c, mw.logger, httperrors.ErrNotFound)
		return nil, httperrors.NewNotFoundError(httperrors.ErrNotFound)
	}

	return user, nil
}
