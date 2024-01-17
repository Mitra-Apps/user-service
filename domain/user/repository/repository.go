package repository

import (
	"context"

	"github.com/Mitra-Apps/be-user-service/domain/user/entity"
	"github.com/google/uuid"
)

//go:generate mockgen -source=repository.go -destination=mock/repository.go -package=mock
type User interface {
	GetAll(ctx context.Context) ([]*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	Create(ctx context.Context, user *entity.User, roleIds []string) error
}

type Role interface {
	Create(ctx context.Context, role *entity.Role) error
	GetRole(ctx context.Context) ([]entity.Role, error)
	GetRoleByUserId(ctx context.Context, userID uuid.UUID) ([]*entity.Role, error)
}
