package service

import (
	"context"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/Mitra-Apps/be-user-service/domain/proto/user"
	"github.com/Mitra-Apps/be-user-service/domain/user/entity"
	epb "google.golang.org/genproto/googleapis/rpc/errdetails"
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
	user := &entity.User{
		Username:    req.Email,
		Password:    req.Password,
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
		Name:        req.Name,
		Address:     req.Address,
	}
	if err := s.userRepository.Create(ctx, user, req.RoleId); err != nil {
		fmt.Println(err.Error())
		st := status.New(codes.InvalidArgument, "Invalid input")
		ds, err := st.WithDetails(
			&epb.BadRequest{
				FieldViolations: []*epb.BadRequest_FieldViolation{
					{
						Field:       "Email or Phone Number",
						Description: "Email or Phone number already exist",
					},
				},
			},
		)
		if err != nil {
			return st.Err()
		}
		return ds.Err()
	}

	return nil
}

func (s *Service) CreateRole(ctx context.Context, role *entity.Role) error {
	return s.roleRepo.Create(ctx, role)
}

func checkPassword(password, hashedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err
}
