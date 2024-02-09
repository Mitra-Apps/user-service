package postgre

import (
	"context"
	"fmt"

	"github.com/Mitra-Apps/be-user-service/domain/user/entity"
	"github.com/Mitra-Apps/be-user-service/domain/user/repository"
	"github.com/google/uuid"

	"gorm.io/gorm"
)

type userRepoImpl struct {
	db *gorm.DB
}

func NewUserRepoImpl(db *gorm.DB) repository.User {
	return &userRepoImpl{
		db: db,
	}
}

func (p *userRepoImpl) GetAll(ctx context.Context) ([]*entity.User, error) {
	var accounts []*entity.User
	if err := p.db.Order("created_at DESC").Find(&accounts).Error; err != nil {
		return nil, err
	}
	return accounts, nil
}

func (p *userRepoImpl) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user *entity.User
	if err := p.db.Preload("Roles").Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (p *userRepoImpl) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	var user *entity.User
	if err := p.db.Preload("Roles").Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (p *userRepoImpl) Create(ctx context.Context, user *entity.User, roleIds []string) error {
	err := p.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return err
		}
		for _, roleId := range roleIds {
			fmt.Println(roleId)
			if err := tx.Exec("Insert into user_roles (user_id,role_id) values (?,?)",
				user.Id, roleId).Error; err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

func (p *userRepoImpl) Save(ctx context.Context, user *entity.User) error {
	return p.db.WithContext(ctx).Save(user).Error
}

func (p *userRepoImpl) VerifyUserByEmail(ctx context.Context, email string) (bool, error) {
	var user *entity.User
	updatedFields := map[string]interface{}{
		"is_verified": true,
	}
	res := p.db.Model(user).Where("email = ?", email).Updates(updatedFields)
	if res.Error != nil {
		return false, res.Error
	}
	return true, nil
}
