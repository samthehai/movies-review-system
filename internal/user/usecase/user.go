package usecase

import (
	"context"
	"fmt"
	"net/http"

	"github.com/samthehai/ml-backend-test-samthehai/config"
	"github.com/samthehai/ml-backend-test-samthehai/internal/entity"
	"github.com/samthehai/ml-backend-test-samthehai/internal/user/usecase/repository"
	"github.com/samthehai/ml-backend-test-samthehai/pkg/httperrors"
	"github.com/samthehai/ml-backend-test-samthehai/pkg/logger"
	"github.com/samthehai/ml-backend-test-samthehai/pkg/token"
	"github.com/samthehai/ml-backend-test-samthehai/pkg/utils"
)

type userUsecase struct {
	cfg            config.Config
	userRepository repository.UserRepository
	logger         logger.Logger
	tokenMaker     token.Maker
}

func NewUserUsecase(cfg config.Config, userRepository repository.UserRepository, log logger.Logger, tokenMaker token.Maker) *userUsecase {
	return &userUsecase{cfg: cfg, userRepository: userRepository, logger: log, tokenMaker: tokenMaker}
}

type RegisterParams struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (u *userUsecase) Register(ctx context.Context, args RegisterParams) (*entity.User, error) {
	existsUser, err := u.userRepository.FindByEmail(ctx, args.Email)
	if err != nil {
		return nil, httperrors.NewInternalServerError(fmt.Errorf("userUsecase.Register.FindByEmail: %w", err))
	}

	if existsUser != nil {
		return nil, httperrors.NewRestError(http.StatusBadRequest, httperrors.ErrExistsEmail.Error(), nil)
	}

	hashedPassword, err := utils.HashedPassword(args.Password)
	if err != nil {
		return nil, httperrors.NewBadRequestError(fmt.Errorf("userUsecase.Register.HashedPassword: %w", err))
	}

	user, err := u.userRepository.Register(ctx, repository.RegisterParams{
		Username:       args.Username,
		Email:          args.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		return nil, httperrors.NewInternalServerError(fmt.Errorf("userUsecase.Register.userRepository.Register: %w", err))
	}

	return user, nil
}

type LoginParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (u *userUsecase) Login(ctx context.Context, args LoginParams) (*entity.UserWithAccessToken, error) {
	user, err := u.userRepository.FindByEmail(ctx, args.Email)
	if err != nil {
		return nil, httperrors.NewInternalServerError(fmt.Errorf("userUsecase.Login.FindByEmail: %w", err))
	}

	if user == nil {
		return nil, httperrors.NewRestError(http.StatusBadRequest, httperrors.ErrBadRequest.Error(), nil)
	}

	if err := utils.CheckPassword(args.Password, user.HashedPassword); err != nil {
		return nil, httperrors.NewRestError(http.StatusBadRequest, httperrors.ErrBadRequest.Error(), nil)
	}

	accessToken, err := u.tokenMaker.CreateToken(user.Username, u.cfg.Server.AccessTokenDuration)
	if err != nil {
		return nil, httperrors.NewRestError(http.StatusInternalServerError, httperrors.ErrInternalServer.Error(), nil)
	}

	return &entity.UserWithAccessToken{
		ID:             user.ID,
		Username:       user.Username,
		Email:          user.Email,
		HashedPassword: user.HashedPassword,
		AccessToken:    accessToken,
	}, nil
}
