package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/samthehai/ml-backend-test-samthehai/internal/entity"
	"github.com/samthehai/ml-backend-test-samthehai/internal/user/usecase/repository"
)

type userRepository struct {
	connManager ConnManager
}

func NewUserRepository(connManager ConnManager) *userRepository {
	return &userRepository{connManager: connManager}
}

const registerQuery = `INSERT INTO users(username, email, hashed_password) VALUES (?,?,?)`
const findByID = `SELECT id, username, email, hashed_password FROM users WHERE id = ?`

func (r *userRepository) Register(ctx context.Context, args repository.RegisterParams) (*entity.User, error) {
	res, err := r.connManager.GetWriter().ExecContext(ctx, registerQuery, args.Username, args.Email, args.HashedPassword)
	if err != nil {
		return nil, fmt.Errorf("userRepository.Register.ExecContext: %w", err)
	}

	createdUserID, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("userRepository.Register.LastInsertId: %w", err)
	}

	u := &User{}
	if err := r.connManager.GetWriter().QueryRowxContext(ctx, findByID, createdUserID).StructScan(u); err != nil {
		return nil, fmt.Errorf("userRepository.Register.QueryRowxContext: %w", err)
	}

	return &entity.User{
		ID:             u.ID,
		Username:       u.Username,
		Email:          u.Email,
		HashedPassword: u.HashedPassword,
	}, nil
}

const findByEmail = `SELECT id, username, email, hashed_password FROM users WHERE email = ?`

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	foundUser := &User{}
	if err := r.connManager.GetReader().QueryRowxContext(ctx, findByEmail, email).StructScan(foundUser); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("userRepository.FindByEmail.QueryRowxContext: %w", err)
	}
	return &entity.User{
		ID:             foundUser.ID,
		Username:       foundUser.Username,
		Email:          foundUser.Email,
		HashedPassword: foundUser.HashedPassword,
	}, nil
}
