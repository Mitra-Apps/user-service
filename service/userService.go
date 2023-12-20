package service

import (
	"context"

	"golang.org/x/crypto/bcrypt"

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
	user, err := s.userRepository.GetByEmail(ctx, payload.Username)
	if err != nil {
		return nil, err
	}
	err = checkPassword(payload.Password, user.Password)
	if err != nil {
		err = ErrWrongPassword
		return nil, err
	}
	return user, nil
}

func checkPassword(password, hashedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err
}
