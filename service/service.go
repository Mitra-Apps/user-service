package service

import (
	"context"

	pb "github.com/Mitra-Apps/be-user-service/domain/proto/user"
	"github.com/Mitra-Apps/be-user-service/domain/user/entity"
	"github.com/Mitra-Apps/be-user-service/domain/user/repository"
)

type Service struct {
	userRepository repository.User
	roleRepo       repository.Role
}

func New(userRepository repository.User, roleRepo repository.Role) *Service {
	return &Service{
		userRepository: userRepository,
		roleRepo:       roleRepo,
	}
}

//go:generate mockgen -source=service.go -destination=mock/service.go -package=mock
type ServiceInterface interface {
	GetAll(ctx context.Context) ([]*entity.User, error)
	Login(ctx context.Context, payload entity.LoginRequest) (*entity.User, []int64, error)
	Register(ctx context.Context, req *pb.UserRegisterRequest) error
	CreateRole(ctx context.Context, role *entity.Role) error
	GetRole(ctx context.Context) ([]entity.Role, error)
}
