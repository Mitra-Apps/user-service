package postgre

import (
	"context"

	"github.com/Mitra-Apps/be-user-service/domain/user/entity"
	"github.com/Mitra-Apps/be-user-service/domain/user/repository"

	"gorm.io/gorm"
)

type RoleRepoImpl struct {
	db *gorm.DB
}

func NewRoleRepoImpl(db *gorm.DB) repository.Role {
	return &RoleRepoImpl{
		db: db,
	}
}

func (r *RoleRepoImpl) Create(ctx context.Context, role *entity.Role) error {
	return r.db.WithContext(ctx).Create(role).Error
}

func (r *RoleRepoImpl) GetRole(ctx context.Context) ([]entity.Role, error) {
	var roles []entity.Role
	if err := r.db.WithContext(ctx).Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}
