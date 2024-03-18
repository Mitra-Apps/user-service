package service

import (
	"context"

	"github.com/Mitra-Apps/be-user-service/config/tools"
	"github.com/Mitra-Apps/be-user-service/config/tools/redis"
	pb "github.com/Mitra-Apps/be-user-service/domain/proto/user"
	"github.com/Mitra-Apps/be-user-service/domain/user/entity"
	"github.com/Mitra-Apps/be-user-service/domain/user/repository"
	"google.golang.org/grpc/codes"
)

type Service struct {
	userRepository repository.User
	roleRepo       repository.Role
	hashing        tools.BcryptInterface
	redis          redis.RedisInterface
	auth           Authentication
}

var (
	ErrorCode       codes.Code
	ErrorCodeDetail string
	ErrorMessage    string
)

func New(
	userRepository repository.User,
	roleRepo repository.Role,
	hashing tools.BcryptInterface,
	redis redis.RedisInterface,
	auth Authentication) *Service {
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
	Register(ctx context.Context, req *pb.UserRegisterRequest) (*entity.OtpMailReq, error)
	CreateRole(ctx context.Context, role *entity.Role) error
	GetRole(ctx context.Context) ([]entity.Role, error)
	VerifyOTP(ctx context.Context, otp int, redisKey string) (user *entity.User, err error)
	ResendOTP(ctx context.Context, email string) (*entity.OtpMailReq, error)
	ChangePassword(ctx context.Context, req *pb.ChangePasswordRequest) (*entity.User, error)
	Logout(ctx context.Context, req *pb.LogoutRequest) error
}
