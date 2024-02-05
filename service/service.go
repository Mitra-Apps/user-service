package service

import (
	"context"

	"github.com/Mitra-Apps/be-user-service/config/tools"
	pb "github.com/Mitra-Apps/be-user-service/domain/proto/user"
	"github.com/Mitra-Apps/be-user-service/domain/user/entity"
	"github.com/Mitra-Apps/be-user-service/domain/user/repository"
	"github.com/go-redis/redis/v8"
)

type Service struct {
	userRepository repository.User
	roleRepo       repository.Role
	hashing        tools.BcryptInterface
	redis          *redis.Client
	auth           Authentication
}

func New(
	userRepository repository.User,
	roleRepo repository.Role, hashing tools.BcryptInterface,
	redis *redis.Client, auth Authentication) *Service {
	return &Service{
		userRepository: userRepository,
		roleRepo:       roleRepo,
		hashing:        hashing,
		redis:          redis,
		auth:           auth,
	}
}

//go:generate mockgen -source=service.go -destination=mock/service.go -package=mock
type ServiceInterface interface {
	GetAll(ctx context.Context) ([]*entity.User, error)
	Login(ctx context.Context, payload entity.LoginRequest) (*entity.User, error)
	Register(ctx context.Context, req *pb.UserRegisterRequest) (string, error)
	CreateRole(ctx context.Context, role *entity.Role) error
	GetRole(ctx context.Context) ([]entity.Role, error)
	VerifyOTP(ctx context.Context, otp int, redisKey string) (result bool, err error)
	ResendOTP(ctx context.Context, email string) (otp int, err error)
}
