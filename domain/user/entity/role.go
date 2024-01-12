package entity

import (
	"strconv"
	"time"

	pb "github.com/Mitra-Apps/be-user-service/domain/proto/user"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Role struct {
	gorm.Model
	RoleName    string         `gorm:"type:varchar(50);not null"`
	Description string         `gorm:"type:varchar(255);null"`
	IsActive    bool           `gorm:"type:bool;not null;default:TRUE"`
	Permission  datatypes.JSON `gorm:"null"`
}

func (r *Role) FromProto(role *pb.Role) error {
	if role.Id != "" {
		id, err := strconv.Atoi(role.Id)
		if err != nil {
			return err
		}
		r.ID = uint(id)
	}
	r.CreatedAt = time.Now()
	r.IsActive = role.IsActive
	r.RoleName = role.RoleName
	r.Description = role.Description
	r.Permission = datatypes.JSON(role.Permission)

	return nil
}
