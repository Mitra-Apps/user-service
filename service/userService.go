package service

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Mitra-Apps/be-user-service/auth"
	"github.com/Mitra-Apps/be-user-service/config"
	pb "github.com/Mitra-Apps/be-user-service/domain/proto/user"
	"github.com/Mitra-Apps/be-user-service/domain/user/entity"
	"github.com/labstack/echo"
)

func (s *Service) GetAll(ctx context.Context) ([]*entity.User, error) {
	users, err := s.userRepository.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *Service) Login(ctx context.Context, payload entity.LoginRequest) (string, error) {
	if strings.Trim(payload.Username, " ") == "" {
		return "", status.Errorf(codes.InvalidArgument, "Username is required")
	}
	if strings.Trim(payload.Password, " ") == "" {
		return "", status.Errorf(codes.InvalidArgument, "Password is required")
	}
	user, err := s.userRepository.GetByEmail(ctx, payload.Username)
	if user == nil && err == nil {
		return "", status.Errorf(codes.InvalidArgument, "Invalid username")
	}

	if err != nil {
		return "", status.Errorf(codes.Internal, "Error getting user by email")
	}
	err = checkPassword(payload.Password, user.Password)
	if err != nil {
		return "", status.Errorf(codes.InvalidArgument, "Invalid password")
	}
	jwt, err := auth.GenerateToken(ctx, user)
	if err != nil {
		echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return jwt, nil
}

func (s *Service) Register(ctx context.Context, req *pb.UserRegisterRequest) (string, error) {
	fmt.Println("register service")
	//hashing password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
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
		return "", NewError(codes.InvalidArgument, errResponse)
	}

	otp := 0
	generateNumber, err := s.generateUnique4DigitNumber()
	if err != nil {
		fmt.Print("Error Generate OTP TOKEN ", err)
	} else {
		otp = generateNumber
	}
	otpString := strconv.Itoa(otp)
	if otp != 0 {
		redisKey := "otp:" + otpString
		err = s.redis.Set(ctx, redisKey, "", time.Minute*5).Err()
		if err != nil {
			log.Print("Error Set Value to Redis")
		}
	}
	return otpString, nil
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

func (s *Service) generateUnique4DigitNumber() (int, error) {
	for {
		randomNumber := generateRandom4DigitNumber()
		// Check if the number exists in Redis
		key := "otp:" + strconv.Itoa(randomNumber)
		exists, err := s.redis.Exists(s.redis.Context(), key).Result()
		if err != nil {
			return 0, err
		}

		// If the number doesn't exist in Redis, return it
		if exists == 0 {
			return randomNumber, nil
		}
	}
}

func generateRandom4DigitNumber() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(9000) + 1000 // Ensure a 4-digit number
}

func checkPassword(password, hashedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err
}
