package usecase

import (
	"context"
	"fmt"
	"practice7/internal/entity"
	"practice7/internal/usercase/repo"
	"practice7/utils"

	"github.com/google/uuid"
)

type UserUseCase struct {
	repo *repo.UserRepo
}

func NewUserUseCase(r *repo.UserRepo) *UserUseCase {
	return &UserUseCase{
		repo: r,
	}
}

func (u *UserUseCase) RegisterUser(ctx context.Context, user *entity.User) (*entity.User, error) {
	createdUser, err := u.repo.RegisterUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("register user: %w", err)
	}
	return createdUser, nil
}

func (u *UserUseCase) LoginUser(ctx context.Context, input *entity.LoginUserDTO) (string, error) {
	userFromRepo, err := u.repo.LoginUser(ctx, input.Username)
	if err != nil {
		return "", fmt.Errorf("user not found: %w", err)
	}

	if !utils.CheckPassword(userFromRepo.Password, input.Password) {
		return "", fmt.Errorf("invalid password")
	}

	token, err := utils.GenerateJWT(userFromRepo.ID, userFromRepo.Role)
	if err != nil {
		return "", fmt.Errorf("generate JWT: %w", err)
	}

	return token, nil
}

func (u *UserUseCase) GetMe(ctx context.Context, userID uuid.UUID) (*entity.User, error) {
	user, err := u.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get user by id: %w", err)
	}
	return user, nil
}

func (u *UserUseCase) PromoteUser(ctx context.Context, userID uuid.UUID) error {
	err := u.repo.UpdateUserRole(ctx, userID)
	if err != nil {
		return fmt.Errorf("promote user: %w", err)
	}
	return nil
}