package service

import (
	"context"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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

func checkPassword(password, hashedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err
}
