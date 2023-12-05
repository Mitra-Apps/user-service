package entity

import (
	"time"

	pb "github.com/Mitra-Apps/be-user-service/domain/user/proto"
)

type User struct {
	Id          int64 `gorm:"primaryKey"`
	Username    string
	Email       string
	PhoneNumber string
	Password    string
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (u *User) ToProto() *pb.User {
	return &pb.User{
		Id:          u.Id,
		Username:    u.Username,
		Email:       u.Email,
		PhoneNumber: u.PhoneNumber,
		Password:    u.Password,
		IsActive:    u.IsActive,
	}
}
