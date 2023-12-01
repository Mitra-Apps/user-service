package entity

import "time"

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
