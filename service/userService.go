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
	"github.com/go-redis/redis"
)

func (s *Service) GetAll(ctx context.Context) ([]*entity.User, error) {
	users, err := s.userRepository.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *Service) Login(ctx context.Context, payload entity.LoginRequest) (*entity.LoginResponse, error) {
	var (
		code        codes.Code
		errResponse *tools.ErrorResponse
	)

	user, err := s.userRepository.GetByEmail(ctx, payload.Email)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			code = codes.NotFound
			errResponse = &tools.ErrorResponse{
				Code:       code.String(),
				CodeDetail: pbErr.ErrorCode_AUTH_LOGIN_NOT_FOUND.String(),
				Message:    err.Error(),
			}
		} else {
			code = codes.Internal
			errResponse = &tools.ErrorResponse{
				Code:       code.String(),
				CodeDetail: pbErr.ErrorCode_UNKNOWN.String(),
				Message:    err.Error(),
			}
		}
		return nil, NewError(code, errResponse)
	}

	if !user.IsActive {
		code = codes.InvalidArgument
		errResponse = &tools.ErrorResponse{
			Code:       code.String(),
			CodeDetail: pbErr.ErrorCode_AUTH_LOGIN_USER_UNVERIFIED.String(),
			Message:    "Email sudah terdaftar, silahkan lakukan verifikasi OTP",
		}
		return nil, NewError(code, errResponse)
	}

	if err := s.hashing.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password)); err != nil {
		code = codes.InvalidArgument
		errResponse = &tools.ErrorResponse{
			Code:       code.String(),
			CodeDetail: pbErr.ErrorCode_AUTH_LOGIN_PASSWORD_INCORRECT.String(),
			Message:    err.Error(),
		}
		return nil, NewError(code, errResponse)
	}

	res := &entity.LoginResponse{
		AccessToken:  "TODO:will add later",
		RefreshToken: "TODO:will add later",
	}

	return res, nil
}

func (s *Service) Register(ctx context.Context, req *pb.UserRegisterRequest) (string, error) {
	var errResponse *tools.ErrorResponse

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
		IsActive:    false,
	}

	data, err := s.userRepository.GetByEmail(ctx, req.Email)
	if err != nil && !strings.Contains(err.Error(), "not found") {
		errResponse = &tools.ErrorResponse{
			Code:       codes.Internal.String(),
			CodeDetail: pbErr.ErrorCode_UNKNOWN.String(),
			Message:    err.Error(),
		}
		return "", NewError(codes.Internal, errResponse)
	}

	if data != nil {
		switch data.IsVerified {
		case false:
			errResponse = &tools.ErrorResponse{
				Code:       codes.InvalidArgument.String(),
				CodeDetail: pbErr.ErrorCode_AUTH_REGISTER_USER_UNVERIFIED.String(), //TODO:, check any detail error code needed
				Message:    "Email sudah terdaftar, mohon ke login page.",
			}
		case true:
			errResponse = &tools.ErrorResponse{
				Code:       codes.InvalidArgument.String(),
				CodeDetail: pbErr.ErrorCode_AUTH_REGISTER_USER_VERIFIED.String(), //TODO:, check any detail error code needed
				Message:    "Email dan/atau No. Telp sudah terdaftar.",
			}
		}
		return "", NewError(codes.InvalidArgument, errResponse)
	}

	if err := s.userRepository.Create(ctx, user, req.RoleId); err != nil {
		errResponse = &tools.ErrorResponse{
			Code:       codes.Internal.String(),
			CodeDetail: pbErr.ErrorCode_UNKNOWN.String(),
			Message:    err.Error(),
		}
		return "", NewError(codes.Internal, errResponse)
	}

	otp := generateRandom4DigitNumber()
	otpString := strconv.Itoa(otp)
	redisPayload := map[string]interface{}{
		"OTP": otpString,
	}
	if otp != 0 {
		redisKey := otpRedisPrefix + data.Email
		err = s.redis.Set(ctx, redisKey, redisPayload, time.Minute*5).Err()
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

func (s *Service) VerifyOTP(ctx context.Context, otp int, redisKey string) (result bool, err error) {
	storedJSON, err := s.redis.Get(s.redis.Context(), redisKey).Result()
	if err == redis.Nil {
		errResponse := &tools.ErrorResponse{
			Code:       codes.InvalidArgument.String(),
			CodeDetail: codes.InvalidArgument.String(),
			Message:    "Key Not Found",
		}
		return false, NewError(codes.InvalidArgument, errResponse)
	} else if err != nil {
		errResponse := &tools.ErrorResponse{
			Code:       codes.InvalidArgument.String(),
			CodeDetail: codes.InvalidArgument.String(),
			Message:    "Redis Error",
		}
		return false, NewError(codes.InvalidArgument, errResponse)
	} else {
		fmt.Println("Retrieved JSON string from Redis:", storedJSON)
		var retrievedObject map[string]interface{}

		if err := json.Unmarshal([]byte(storedJSON), &retrievedObject); err != nil {
			log.Print("Error unmarshaling JSON:", err)
			errResponse := &tools.ErrorResponse{
				Code:       codes.InvalidArgument.String(),
				CodeDetail: codes.InvalidArgument.String(),
				Message:    "Unmarshal Error",
			}
			return false, NewError(codes.InvalidArgument, errResponse)
		}
		if retrievedObject["OTP"] == strconv.Itoa(otp) {
			email := strings.Replace(redisKey, otpRedisPrefix, "", -1)
			_, err := s.userRepository.VerifyUserByEmail(ctx, email)
			if err != nil {
				errResponse := &tools.ErrorResponse{
					Code:       codes.InvalidArgument.String(),
					CodeDetail: codes.InvalidArgument.String(),
					Message:    "Verify User Error",
				}
				return false, NewError(codes.InvalidArgument, errResponse)
			}
		} else {
			errResponse := &tools.ErrorResponse{
				Code:       codes.InvalidArgument.String(),
				CodeDetail: codes.InvalidArgument.String(),
				Message:    "OTP Is Not Match",
			}
			return false, NewError(codes.InvalidArgument, errResponse)
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
	redisKey := otpRedisPrefix + email
	// Marshal the JSON data
	jsonData, err := json.Marshal(redisPayload)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}
	err = s.redis.Set(ctx, redisKey, jsonData, time.Minute*5).Err()
	if err != nil {
		errResponse := &tools.ErrorResponse{
			Code:       codes.InvalidArgument.String(),
			CodeDetail: codes.InvalidArgument.String(),
			Message:    "Error Set New Value To Redis",
		}
		return 0, NewError(codes.InvalidArgument, errResponse)
	}

	return otp, nil
}

func generateRandom4DigitNumber() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(9000) + 1000 // Ensure a 4-digit number
}
