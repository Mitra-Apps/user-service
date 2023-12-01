package repository

import (
	"context"
	"user-service/domain/user/entity"
)

type UserInterface interface {
	GetAll(ctx context.Context) ([]*entity.User, error)
}
