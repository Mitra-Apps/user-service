package postgre

import (
	"context"

	"github.com/Mitra-Apps/be-user-service/domain/user/entity"

	"gorm.io/gorm"
)

type Postgre struct {
	db *gorm.DB
}

func NewPostgre(db *gorm.DB) *Postgre {
	return &Postgre{db}
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
