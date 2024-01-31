package service

import (
	"context"
	"encoding/json"
	"fmt"
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
	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo"
	"github.com/rs/zerolog/log"
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
	redisKey := "otp:" + req.Email
	generateNumber := generateRandom4DigitNumber()
	otpString := strconv.Itoa(generateNumber)
	redisPayload := RedisOTP{
		Email: req.Email,
		OTP:   otpString,
	}
	err = s.redis.Set(ctx, redisKey, redisPayload, time.Minute*5).Err()
	if err != nil {
		log.Print("Error Set Value to Redis")
		errResponse := &config.ErrorResponse{
			Code:       codes.InvalidArgument.String(),
			CodeDetail: codes.InvalidArgument.String(),
			Message:    "Error Set Value to Redis",
		}
		return "", NewError(codes.InvalidArgument, errResponse)
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

func (s *Service) VerifyOTP(ctx context.Context, otp int, redisKey string) (result bool, err error) {
	storedJSON, err := s.redis.Get(s.redis.Context(), redisKey).Result()
	if err == redis.Nil {
		errResponse := &config.ErrorResponse{
			Code:       codes.InvalidArgument.String(),
			CodeDetail: codes.InvalidArgument.String(),
			Message:    "Key Not Found",
		}
		return false, NewError(codes.InvalidArgument, errResponse)
	} else if err != nil {
		errResponse := &config.ErrorResponse{
			Code:       codes.InvalidArgument.String(),
			CodeDetail: codes.InvalidArgument.String(),
			Message:    "Redis Error",
		}
		return false, NewError(codes.InvalidArgument, errResponse)
	} else {
		fmt.Println("Retrieved JSON string from Redis:", storedJSON)
		// Convert the JSON string back to a JSON object
		var retrievedObject map[string]interface{}
		if err := json.Unmarshal([]byte(storedJSON), &retrievedObject); err != nil {
			fmt.Println("Error unmarshaling JSON:", err)
			errResponse := &config.ErrorResponse{
				Code:       codes.InvalidArgument.String(),
				CodeDetail: codes.InvalidArgument.String(),
				Message:    "Unmarshal Error",
			}
			return false, NewError(codes.InvalidArgument, errResponse)
		}

		fmt.Println("Parsed JSON object:", retrievedObject)
		return true, nil
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
