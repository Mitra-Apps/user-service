package service

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"

	pbErr "github.com/Mitra-Apps/be-user-service/domain/proto"
	pb "github.com/Mitra-Apps/be-user-service/domain/proto/user"
	"github.com/Mitra-Apps/be-user-service/domain/user/entity"
	"github.com/Mitra-Apps/be-user-service/handler/middleware"
	"github.com/google/uuid"
)

func (s *Service) GetAll(ctx context.Context) ([]*entity.User, error) {
	users, err := s.userRepository.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *Service) Login(ctx context.Context, payload entity.LoginRequest) (uuid.UUID, error) {

	user, err := s.userRepository.GetByEmail(ctx, payload.Email)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			ErrorCode = codes.NotFound
			ErrorCodeDetail = pbErr.ErrorCode_AUTH_LOGIN_NOT_FOUND.String()
			ErrorMessage = "Email belum terdaftar, mohon registrasi"
		} else {
			ErrorCode = codes.Internal
			ErrorCodeDetail = pbErr.ErrorCode_UNKNOWN.String()
			ErrorMessage = err.Error()
		}
		return uuid.Nil, NewError(ErrorCode, ErrorCodeDetail, ErrorMessage)
	}

	if !user.IsVerified {
		ErrorCode = codes.InvalidArgument
		ErrorCodeDetail = pbErr.ErrorCode_AUTH_LOGIN_USER_UNVERIFIED.String()
		ErrorMessage = "Email sudah terdaftar, silahkan lakukan verifikasi OTP"
		return uuid.Nil, NewError(ErrorCode, ErrorCodeDetail, ErrorMessage)
	}

	if err := s.hashing.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password)); err != nil {
		ErrorCode = codes.InvalidArgument
		ErrorCodeDetail = pbErr.ErrorCode_AUTH_LOGIN_PASSWORD_INCORRECT.String()
		ErrorMessage = "Data yang dimasukkan tidak sesuai"
		return uuid.Nil, NewError(ErrorCode, ErrorCodeDetail, ErrorMessage)
	}

	return user.Id, nil
}

func (s *Service) Register(ctx context.Context, req *pb.UserRegisterRequest) (string, error) {
	//hashing password
	hashedPassword, err := s.hashing.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	user := &entity.User{
		Email:       req.Email,
		Password:    string(hashedPassword),
		Username:    req.Email,
		PhoneNumber: req.PhoneNumber,
		Name:        req.Name,
		Address:     req.Address,
	}

	data, err := s.userRepository.GetByEmail(ctx, req.Email)
	if err != nil && !strings.Contains(err.Error(), "not found") {
		ErrorCode = codes.Internal
		ErrorCodeDetail = pbErr.ErrorCode_UNKNOWN.String()
		ErrorMessage = err.Error()
		return "", NewError(ErrorCode, ErrorCodeDetail, ErrorMessage)
	}

	if data != nil {
		switch data.IsVerified {
		case false:
			ErrorCode = codes.InvalidArgument
			ErrorCodeDetail = pbErr.ErrorCode_AUTH_REGISTER_USER_UNVERIFIED.String()
			ErrorMessage = "Email sudah terdaftar, mohon ke halaman login."
		case true:
			ErrorCode = codes.InvalidArgument
			ErrorCodeDetail = pbErr.ErrorCode_AUTH_REGISTER_USER_VERIFIED.String()
			ErrorMessage = "Email dan/atau No. Telp sudah terdaftar."
		}
		return "", NewError(ErrorCode, ErrorCodeDetail, ErrorMessage)
	}

	if err := s.userRepository.Create(ctx, user, req.RoleId); err != nil {
		ErrorCode = codes.Internal
		ErrorCodeDetail = pbErr.ErrorCode_UNKNOWN.String()
		ErrorMessage = err.Error()
		return "", NewError(ErrorCode, ErrorCodeDetail, ErrorMessage)
	}

	otp := 0
	generateNumber, err := s.generateUnique4DigitNumber()
	if err == nil {
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
	fmt.Println("get role service", middleware.GetUserIDValue(ctx))
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
