package postgre

import (
	"context"

	"github.com/Mitra-Apps/be-user-service/domain/user/entity"
	"github.com/Mitra-Apps/be-user-service/domain/user/repository"
	"github.com/google/uuid"

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

func (r *RoleRepoImpl) GetRoleByUserId(ctx context.Context, userID uuid.UUID) ([]*entity.Role, error) {
	var roles []*entity.Role
	if tx := r.db.WithContext(ctx).Raw(`SELECT r.id, r.created_at, r.updated_at, r.deleted_at, r.role_name, r.description, `+
		`r.is_active, r."permission" `+
		`FROM roles r inner join user_roles ur on r.id = ur.role_id `+
		`where  ur.user_id = ? and r.is_active = TRUE`, userID).
		Scan(&roles); tx.Error != nil {
		return nil, tx.Error
	}
	return roles, nil
}
