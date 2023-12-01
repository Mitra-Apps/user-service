package repository

import (
	"context"

	"github.com/Mitra-Apps/user-service/domain/user/entity"
)

type UserInterface interface {
	GetAll(ctx context.Context) ([]*entity.User, error)
}
