package usecase

import (
	"context"
	"practice7/internal/entity"

	"github.com/google/uuid"
)

type UserInterface interface {
	RegisterUser(ctx context.Context, user *entity.User) (*entity.User, error)
	LoginUser(ctx context.Context, input *entity.LoginUserDTO) (string, error)
	GetMe(ctx context.Context, userID uuid.UUID) (*entity.User, error)
	PromoteUser(ctx context.Context, userID uuid.UUID) error
}