package repo

import (
	"context"
	"fmt"
	"practice7/internal/entity"
	"practice7/pkg/postgres"

	"github.com/google/uuid"
)

type UserRepo struct {
	PG *postgres.Postgres
}

func NewUserRepo(pg *postgres.Postgres) *UserRepo {
	return &UserRepo{pg}
}

func (u *UserRepo) RegisterUser(ctx context.Context, user *entity.User) (*entity.User, error) {
	err := u.PG.Conn.WithContext(ctx).Create(user).Error
	if err != nil {
		return nil, fmt.Errorf("register user: %w", err)
	}
	return user, nil
}

func (u *UserRepo) LoginUser(ctx context.Context, username string) (*entity.User, error) {
	var userFromDB entity.User
	err := u.PG.Conn.WithContext(ctx).Where("username = ?", username).First(&userFromDB).Error
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return &userFromDB, nil
}

func (u *UserRepo) GetUserByID(ctx context.Context, userID uuid.UUID) (*entity.User, error) {
	var user entity.User
	err := u.PG.Conn.WithContext(ctx).Where("id = ?", userID).First(&user).Error
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return &user, nil
}

func (u *UserRepo) UpdateUserRole(ctx context.Context, userID uuid.UUID) error {
	err := u.PG.Conn.WithContext(ctx).Model(&entity.User{}).
		Where("id = ?", userID).
		Update("role", "admin").Error
	if err != nil {
		return fmt.Errorf("failed to update user role: %w", err)
	}
	return nil
}