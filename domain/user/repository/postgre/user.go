package postgre

import (
	"context"
	"fmt"

	"github.com/Mitra-Apps/be-user-service/domain/user/entity"
	"github.com/Mitra-Apps/be-user-service/domain/user/repository"

	"gorm.io/gorm"
)

type Postgre struct {
	db *gorm.DB
}

func NewUserRepoImpl(db *gorm.DB) repository.User {
	return &Postgre{
		db: db,
	}
}

func (p *Postgre) GetAll(ctx context.Context) ([]*entity.User, error) {
	var accounts []*entity.User
	res := p.db.Order("created_at DESC").Find(&accounts)
	if res.Error == gorm.ErrEmptySlice || res.RowsAffected == 0 {
		return nil, nil
	}
	if res.Error != nil {
		return nil, res.Error
	}
	return accounts, nil
}

func (p *Postgre) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user *entity.User
	res := p.db.Find(&user).Where("email = ?", email)
	if res.Error == gorm.ErrEmptySlice || res.RowsAffected == 0 {
		return nil, nil
	}
	if res.Error != nil {
		return nil, res.Error
	}
	return user, nil
}

func (p *Postgre) Create(ctx context.Context, user *entity.User, roleIds []string) error {
	err := p.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return err
		}
		for _, roleId := range roleIds {
			fmt.Println(roleId)
			if err := tx.Exec("Insert into user_roles (user_id,role_id) values (?,?)",
				user.Id, roleId).Debug().Error; err != nil {
				return err
			}
		}
		return nil
	})
	return err
}
