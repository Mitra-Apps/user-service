package service

import (
	"context"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Mitra-Apps/be-user-service/config"
	pb "github.com/Mitra-Apps/be-user-service/domain/proto/user"
	"github.com/Mitra-Apps/be-user-service/domain/user/entity"
)

func (s *Service) GetAll(ctx context.Context) ([]*entity.User, error) {
	users, err := s.userRepository.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *Service) Login(ctx context.Context, payload entity.LoginRequest) (*entity.User, error) {
	if strings.Trim(payload.Username, " ") == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Username is required")
	}
	if strings.Trim(payload.Password, " ") == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Password is required")
	}
	user, err := s.userRepository.GetByEmail(ctx, payload.Username)
	if user == nil && err == nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid username")
	}
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error getting user by email")
	}
	err = checkPassword(payload.Password, user.Password)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid password")
	}
	return user, nil
}

func (s *Service) Register(ctx context.Context, req *pb.UserRegisterRequest) error {
	fmt.Println("register service")
	//hashing password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &entity.User{
		Username:    req.Email,
		Password:    string(hashedPassword),
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
		Name:        req.Name,
		Address:     req.Address,
		IsActive:    false,
	}

	if err := s.userRepository.Create(ctx, user, req.RoleId); err != nil {
		errResponse := &config.ErrorResponse{
			Code:       codes.InvalidArgument.String(),
			CodeDetail: codes.InvalidArgument.String(), //TODO:, check any detail error code needed
			Message:    "Email dan/atau No. Telp sudah terdaftar",
		}
		return NewError(codes.InvalidArgument, errResponse)
	}
	return nil
}

func (s *Service) CreateRole(ctx context.Context, role *entity.Role) error {
	return s.roleRepo.Create(ctx, role)
}

func (s *Service) GetRole(ctx context.Context) ([]entity.Role, error) {
	roles, err := s.roleRepo.GetRole(ctx)
	if err != nil {
		return nil, err
	}
	return roles, nil
}

func checkPassword(password, hashedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err
}
