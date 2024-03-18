package entity

import (
	"time"

	pb "github.com/Mitra-Apps/be-user-service/domain/proto/user"
	"github.com/google/uuid"
)

type User struct {
	Id                   uuid.UUID     `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Username             string        `gorm:"type:varchar(255);not null;unique"`
	Password             string        `gorm:"type:varchar(255);not null"`
	Email                string        `gorm:"type:varchar(255);not null;unique"`
	PhoneNumber          string        `gorm:"type:varchar(50);not null;unique"`
	AvatarImageId        uuid.NullUUID `gorm:"type:varchar(255);null"`
	AccessToken          string        `gorm:"type:varchar(255);null"`
	RefreshToken         string        `gorm:"type:varchar(255);null"`
	IsActive             bool          `gorm:"type:bool;not null;default:TRUE"`
	CreatedAt            time.Time     `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP"`
	CreatedBy            uuid.UUID     `gorm:"type:uuid;not null"`
	UpdatedAt            time.Time     `gorm:"type:timestamptz;null;default:CURRENT_TIMESTAMP"`
	UpdatedBy            uuid.NullUUID `gorm:"type:uuid;null"`
	Name                 string        `gorm:"type:varchar(255);not null"`
	Roles                []Role        `gorm:"many2many:user_roles;"`
	Address              string        `gorm:"type:varchar(255);null"`
	IsVerified           bool          `gorm:"type:bool;not null;default:FALSE"`
	WrongPasswordCounter uint
}

func (u *User) ToProto() *pb.User {
	var avatarImageId string
	if u.AvatarImageId.Valid {
		avatarImageId = u.AvatarImageId.UUID.String()
	}
	return &pb.User{
		Id:            u.Id.String(),
		Username:      u.Username,
		Email:         u.Email,
		PhoneNumber:   u.PhoneNumber,
		Password:      u.Password,
		IsActive:      u.IsActive,
		IsVerified:    u.IsVerified,
		AvatarImageId: avatarImageId,
		Name:          u.Name,
		Address:       u.Address,
	}
}

type LoginRequest struct {
	Email    string
	Password string
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type OtpMailReq struct {
	Name    string
	Email   string
	OtpCode int
}
