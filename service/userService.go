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
	util "github.com/Mitra-Apps/be-utility-service/service"
)

func (s *Service) GetAll(ctx context.Context) ([]*entity.User, error) {
	users, err := s.userRepository.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *Service) Login(ctx context.Context, payload entity.LoginRequest) (*entity.User, error) {

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
		return nil, util.NewError(ErrorCode, ErrorCodeDetail, ErrorMessage)
	}

	if user.WrongPasswordCounter >= 3 {
		ErrorCode = codes.InvalidArgument
		ErrorCodeDetail = pbErr.ErrorCode_AUTH_LOGIN_PASSWORD_INCORRECT_3X.String()
		ErrorMessage = "Anda telah melebihi limit kesalahan kata sandi, mohon ganti sandi anda"
		return nil, util.NewError(ErrorCode, ErrorCodeDetail, ErrorMessage)
	}

	if err := s.hashing.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password)); err != nil {
		user.WrongPasswordCounter++
		if err = s.userRepository.Save(ctx, user); err != nil {
			return nil, util.NewError(codes.Internal, codes.Unknown.String(), err.Error())
		}
		ErrorCode = codes.InvalidArgument
		if user.WrongPasswordCounter >= 3 {
			ErrorCodeDetail = pbErr.ErrorCode_AUTH_LOGIN_PASSWORD_INCORRECT_3X.String()
			ErrorMessage = "Anda telah melebihi limit kesalahan kata sandi, mohon ganti sandi anda"
			return nil, util.NewError(ErrorCode, ErrorCodeDetail, ErrorMessage)
		}
		ErrorCodeDetail = pbErr.ErrorCode_AUTH_LOGIN_PASSWORD_INCORRECT.String()
		ErrorMessage = "Data yang dimasukkan tidak sesuai"
		return nil, util.NewError(ErrorCode, ErrorCodeDetail, ErrorMessage)
	}

	if !user.IsVerified {
		ErrorCode = codes.InvalidArgument
		ErrorCodeDetail = pbErr.ErrorCode_AUTH_LOGIN_USER_UNVERIFIED.String()
		ErrorMessage = "Email sudah terdaftar, silahkan lakukan verifikasi OTP"
		return nil, util.NewError(ErrorCode, ErrorCodeDetail, ErrorMessage)
	}

	user.WrongPasswordCounter = 0
	if err = s.userRepository.Save(ctx, user); err != nil {
		return nil, util.NewError(codes.Internal, codes.Unknown.String(), err.Error())
	}

	return user, nil
}

func (s *Service) Register(ctx context.Context, req *pb.UserRegisterRequest) (*entity.OtpMailReq, error) {
	//hashing password
	hashedPassword, err := s.hashing.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
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
		return nil, util.NewError(codes.Internal, codes.Internal.String(), err.Error())
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
		return nil, util.NewError(ErrorCode, ErrorCodeDetail, ErrorMessage)
	}

	if err := s.userRepository.Create(ctx, user, req.RoleId); err != nil {
		ErrorCode = codes.Internal
		ErrorCodeDetail = pbErr.ErrorCode_UNKNOWN.String()
		ErrorMessage = err.Error()
		return nil, util.NewError(ErrorCode, ErrorCodeDetail, ErrorMessage)
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
		return nil, util.NewError(codes.Internal, codes.Internal.String(), err.Error())
	}

	err = s.redis.Set(ctx, redisKey, jsonData, time.Minute*5)
	if err != nil {
		log.Print("Error Set Value to Redis")
		return nil, util.NewError(codes.Internal, codes.Internal.String(), err.Error())
	}

	sendOtpReq := &entity.OtpMailReq{
		Name:    user.Name,
		Email:   user.Email,
		OtpCode: otp,
	}

	return sendOtpReq, nil
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

func (s *Service) VerifyOTP(ctx context.Context, otp int, redisKey string) (user *entity.User, err error) {

	email := strings.Replace(redisKey, tools.OtpRedisPrefix, "", -1)
	user, err = s.userRepository.GetByEmail(ctx, email)
	if err != nil {
		ErrorCode = codes.Internal
		ErrorCodeDetail = pbErr.ErrorCode_RECORD_NOT_FOUND.String()
		ErrorMessage = "Verify User Error"
		return nil, util.NewError(ErrorCode, ErrorCodeDetail, ErrorMessage)
	}
	if user.IsVerified {
		ErrorCode = codes.InvalidArgument
		ErrorCodeDetail = pbErr.ErrorCode_AUTH_OTP_ERROR_VERIFIED_USER.String()
		ErrorMessage = "Email sudah terverifikasi, silahkan login"
		return nil, util.NewError(ErrorCode, ErrorCodeDetail, ErrorMessage)
	}

	if err := verifyOtpFromRedis(s, otp, redisKey); err != nil {
		return nil, err
	}

	if _, err = s.userRepository.VerifyUserByEmail(ctx, user.Email); err != nil {
		ErrorCode = codes.Internal
		ErrorCodeDetail = pbErr.ErrorCode_UNKNOWN.String()
		ErrorMessage = "Update User Error"
		return nil, util.NewError(ErrorCode, ErrorCodeDetail, ErrorMessage)
	}

	return user, nil
}

func verifyOtpFromRedis(s *Service, otp int, redisKey string) error {
	storedJSON, err := s.redis.GetStringKey(s.redis.GetContext(), redisKey)
	if err != nil {
		ErrorCode = codes.InvalidArgument
		ErrorCodeDetail = pbErr.ErrorCode_AUTH_OTP_INVALID.String()
		ErrorMessage = "Kode OTP Tidak Berlaku"
		return util.NewError(ErrorCode, ErrorCodeDetail, ErrorMessage)
	}

	fmt.Println("Retrieved JSON string from Redis:", storedJSON)
	var retrievedObject map[string]interface{}

	if err := json.Unmarshal([]byte(storedJSON), &retrievedObject); err != nil {
		ErrorCode = codes.Internal
		ErrorCodeDetail = pbErr.ErrorCode_UNKNOWN.String()
		ErrorMessage = "Unmarshal Error"
		return util.NewError(ErrorCode, ErrorCodeDetail, ErrorMessage)
	}
	if retrievedObject["OTP"] != strconv.Itoa(otp) {
		ErrorCode = codes.InvalidArgument
		ErrorCodeDetail = pbErr.ErrorCode_AUTH_OTP_INVALID.String()
		ErrorMessage = "Kode OTP Tidak Berlaku"
		return util.NewError(ErrorCode, ErrorCodeDetail, ErrorMessage)
	}
	return nil
}

func (s *Service) ResendOTP(ctx context.Context, email string) (*entity.OtpMailReq, error) {
	otp := generateRandom4DigitNumber()
	otpString := strconv.Itoa(otp)
	redisPayload := map[string]interface{}{
		"OTP": otpString,
	}
	redisKey := tools.OtpRedisPrefix + email

	user, err := s.userRepository.GetByEmail(ctx, email)
	if err != nil {
		return nil, util.NewError(codes.Internal, pbErr.ErrorCode_RECORD_NOT_FOUND.String(), err.Error())
	}
	// Marshal the JSON data
	jsonData, err := json.Marshal(redisPayload)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return nil, util.NewError(codes.Internal, codes.Internal.String(), err.Error())
	}
	err = s.redis.Set(ctx, redisKey, jsonData, time.Minute*5)
	if err != nil {
		ErrorCode = codes.Internal
		ErrorCodeDetail = pbErr.ErrorCode_UNKNOWN.String()
		ErrorMessage = "Set Value Redis Error"
		return nil, util.NewError(ErrorCode, ErrorCodeDetail, ErrorMessage)
	}

	sendOtpReq := &entity.OtpMailReq{
		Name:    user.Name,
		Email:   user.Email,
		OtpCode: otp,
	}

	return sendOtpReq, nil
}

func (s *Service) ChangePassword(ctx context.Context, req *pb.ChangePasswordRequest) (*entity.User, error) {
	user, err := s.userRepository.GetByEmail(ctx, req.Email)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			ErrorCode = codes.NotFound
			ErrorCodeDetail = pbErr.ErrorCode_RECORD_NOT_FOUND.String()
			ErrorMessage = "Email belum terdaftar, mohon registrasi"
		} else {
			ErrorCode = codes.Internal
			ErrorCodeDetail = pbErr.ErrorCode_UNKNOWN.String()
			ErrorMessage = err.Error()
		}
		return nil, util.NewError(ErrorCode, ErrorCodeDetail, ErrorMessage)
	}
	redisKey := "otp:" + req.Email
	if err = verifyOtpFromRedis(s, int(req.OtpCode), redisKey); err != nil {
		ErrorMessage = "Data yang dimasukkan tidak sesuai"
		return nil, util.NewError(ErrorCode, ErrorCodeDetail, ErrorMessage)
	}
	hashedPassword, err := s.hashing.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, util.NewError(codes.Internal, codes.Unknown.String(), err.Error())
	}
	user.Password = string(hashedPassword)
	user.WrongPasswordCounter = 0
	user.IsVerified = true
	if err = s.userRepository.Save(ctx, user); err != nil {
		return nil, util.NewError(codes.Internal, codes.Unknown.String(), err.Error())
	}
	return user, nil
}

func (s *Service) Logout(ctx context.Context, req *pb.LogoutRequest) error {
	user, err := s.userRepository.GetByEmail(ctx, req.Email)
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
		return util.NewError(ErrorCode, ErrorCodeDetail, ErrorMessage)
	}
	user.AccessToken = ""
	user.RefreshToken = ""
	if err = s.userRepository.Save(ctx, user); err != nil {
		return util.NewError(codes.Internal, codes.Unknown.String(), err.Error())
	}
	return nil
}

func generateRandom4DigitNumber() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(9000) + 1000 // Ensure a 4-digit number
}
