package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"

	"github.com/Mitra-Apps/be-user-service/config/tools"
	pbErr "github.com/Mitra-Apps/be-user-service/domain/proto"
	pb "github.com/Mitra-Apps/be-user-service/domain/proto/user"
	"github.com/Mitra-Apps/be-user-service/domain/user/entity"
	"github.com/Mitra-Apps/be-user-service/handler/middleware"
	"github.com/go-redis/redis"
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

	otp := generateRandom4DigitNumber()
	otpString := strconv.Itoa(otp)
	redisPayload := map[string]interface{}{
		"OTP": otpString,
	}
	redisKey := tools.OtpRedisPrefix + req.Email
	jsonData, err := json.Marshal(redisPayload)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return "", NewError(codes.Internal, codes.Internal.String(), err.Error())
	}

	err = s.redis.Set(ctx, redisKey, jsonData, time.Minute*5).Err()
	if err != nil {
		log.Print("Error Set Value to Redis")
		return "", NewError(codes.Internal, codes.Internal.String(), err.Error())
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

func (s *Service) VerifyOTP(ctx context.Context, otp int, redisKey string) (result bool, err error) {
	storedJSON, err := s.redis.Get(s.redis.Context(), redisKey).Result()
	if err == redis.Nil {
		ErrorCode = codes.InvalidArgument
		ErrorCodeDetail = pbErr.ErrorCode_AUTH_OTP_INVALID.String()
		ErrorMessage = "Kode Otp Tidak Berlaku"
		return false, NewError(ErrorCode, ErrorCodeDetail, ErrorMessage)
	} else if err != nil {
		ErrorCode = codes.Internal
		ErrorCodeDetail = pbErr.ErrorCode_UNKNOWN.String()
		ErrorMessage = "Redis Error"
		return false, NewError(ErrorCode, ErrorCodeDetail, ErrorMessage)
	} else {
		fmt.Println("Retrieved JSON string from Redis:", storedJSON)
		var retrievedObject map[string]interface{}

		if err := json.Unmarshal([]byte(storedJSON), &retrievedObject); err != nil {
			ErrorCode = codes.Internal
			ErrorCodeDetail = pbErr.ErrorCode_UNKNOWN.String()
			ErrorMessage = "Unmarshal Error"
			return false, NewError(ErrorCode, ErrorCodeDetail, ErrorMessage)
		}
		if retrievedObject["OTP"] == strconv.Itoa(otp) {
			email := strings.Replace(redisKey, tools.OtpRedisPrefix, "", -1)
			_, err := s.userRepository.VerifyUserByEmail(ctx, email)
			if err != nil {
				ErrorCode = codes.Internal
				ErrorCodeDetail = pbErr.ErrorCode_UNKNOWN.String()
				ErrorMessage = "Verify User Error"
				return false, NewError(ErrorCode, ErrorCodeDetail, ErrorMessage)
			}
		} else {
			ErrorCode = codes.Internal
			ErrorCodeDetail = pbErr.ErrorCode_UNKNOWN.String()
			ErrorMessage = "OTP is not match"
			return false, NewError(ErrorCode, ErrorCodeDetail, ErrorMessage)
		}

		return true, nil
	}
}

func (s *Service) ResendOTP(ctx context.Context, email string) (otp int, err error) {
	otp = generateRandom4DigitNumber()
	otpString := strconv.Itoa(otp)
	redisPayload := map[string]interface{}{
		"OTP": otpString,
	}
	redisKey := tools.OtpRedisPrefix + email
	// Marshal the JSON data
	jsonData, err := json.Marshal(redisPayload)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}
	err = s.redis.Set(ctx, redisKey, jsonData, time.Minute*5).Err()
	if err != nil {
		ErrorCode = codes.Internal
		ErrorCodeDetail = pbErr.ErrorCode_UNKNOWN.String()
		ErrorMessage = "Set Value Redis Error"
		return 0, NewError(ErrorCode, ErrorCodeDetail, ErrorMessage)
	}
	return otp, nil
}

func generateRandom4DigitNumber() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(9000) + 1000 // Ensure a 4-digit number
}
