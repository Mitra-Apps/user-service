package service

import (
	"context"
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
				Message:    "Email belum terdaftar, mohon registrasi",
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

	if !user.IsVerified {
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
			Message:    "Data yang dimasukkan tidak sesuai",
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
