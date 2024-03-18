package service

import (
	"context"
	"errors"
	"reflect"
	"testing"

	mockTools "github.com/Mitra-Apps/be-user-service/config/tools/mock"
	mockRedis "github.com/Mitra-Apps/be-user-service/config/tools/redis/mock"
	pb "github.com/Mitra-Apps/be-user-service/domain/proto/user"
	"github.com/Mitra-Apps/be-user-service/domain/user/entity"
	"github.com/Mitra-Apps/be-user-service/domain/user/repository/mock"
	r "github.com/go-redis/redis"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestService_GetAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockUser := mock.NewMockUser(ctrl)
	mockUserRecord := mockUser.EXPECT()

	users := []*entity.User{
		{
			Name: "test1",
		},
		{
			Name: "test2",
		},
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		s       *Service
		args    args
		want    []*entity.User
		wantErr bool
		mocks   *gomock.Call
	}{
		{
			name: "error get all data",
			s: &Service{
				userRepository: mockUser,
			},
			args: args{
				ctx: context.Background(),
			},
			want:    nil,
			wantErr: true,
			mocks:   mockUserRecord.GetAll(gomock.Any()).Return(nil, errors.New("any error")),
		},
		{
			name: "success",
			s: &Service{
				userRepository: mockUser,
			},
			args: args{
				ctx: context.Background(),
			},
			want:    users,
			wantErr: false,
			mocks:   mockUserRecord.GetAll(gomock.Any()).Return(users, nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.GetAll(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.GetAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.GetAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mock.NewMockUser(ctrl)
	mockLogin := func(user *entity.User, err error) func(m *mock.MockUser) {
		return func(m *mock.MockUser) {
			m.EXPECT().GetByEmail(gomock.Any(), gomock.Any()).Return(user, err)
		}
	}
	mockSaveUser := func(err error) func(m *mock.MockUser) {
		return func(m *mock.MockUser) {
			m.EXPECT().Save(gomock.Any(), gomock.Any()).Return(err)
		}
	}
	mockHash := mockTools.NewMockBcryptInterface(ctrl)
	mockCompareHash := func(err error) func(m *mockTools.MockBcryptInterface) {
		return func(m *mockTools.MockBcryptInterface) {
			m.EXPECT().CompareHashAndPassword(gomock.Any(), gomock.Any()).Return(err)
		}
	}
	userId := uuid.New()
	loginRequest := &entity.LoginRequest{
		Email:    "test@email.com",
		Password: "test@123",
	}
	unverifiedUser := &entity.User{
		Email:      "test@email.com",
		Password:   "test@123",
		IsVerified: false,
	}
	verifiedUser := &entity.User{
		Id:                   userId,
		Email:                "test@email.com",
		Password:             "test@123",
		IsVerified:           true,
		WrongPasswordCounter: 0,
	}
	verifiedUserWrongPass2x := &entity.User{
		Id:                   userId,
		Email:                "test@email.com",
		Password:             "test@123",
		WrongPasswordCounter: 2,
	}
	verifiedUserWrongPass3x := &entity.User{
		Id:                   userId,
		Email:                "test@email.com",
		Password:             "test@123",
		WrongPasswordCounter: 3,
	}
	type args struct {
		ctx     context.Context
		payload entity.LoginRequest
	}
	tests := []struct {
		name    string
		s       *Service
		args    args
		want    *entity.User
		wantErr bool
	}{
		{
			name: "error record not found",
			s: &Service{
				userRepository: mockRepo,
			},
			args: args{
				ctx:     context.Background(),
				payload: *loginRequest,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "unexpected error",
			s: &Service{
				userRepository: mockRepo,
			},
			args: args{
				ctx:     context.Background(),
				payload: *loginRequest,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error unverified account",
			s: &Service{
				userRepository: mockRepo,
				hashing:        mockHash,
			},
			args: args{
				ctx:     context.Background(),
				payload: *loginRequest,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error password incorrect",
			s: &Service{
				userRepository: mockRepo,
				hashing:        mockHash,
			},
			args: args{
				ctx:     context.Background(),
				payload: *loginRequest,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error password incorrect 3x",
			s: &Service{
				userRepository: mockRepo,
				hashing:        mockHash,
			},
			args: args{
				ctx:     context.Background(),
				payload: *loginRequest,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error password incorrect 3x after wrong pass login",
			s: &Service{
				userRepository: mockRepo,
				hashing:        mockHash,
			},
			args: args{
				ctx:     context.Background(),
				payload: *loginRequest,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error saving wrong password counter",
			s: &Service{
				userRepository: mockRepo,
				hashing:        mockHash,
			},
			args: args{
				ctx:     context.Background(),
				payload: *loginRequest,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error saving wrong password counter back to 0",
			s: &Service{
				userRepository: mockRepo,
				hashing:        mockHash,
			},
			args: args{
				ctx:     context.Background(),
				payload: *loginRequest,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "success",
			s: &Service{
				userRepository: mockRepo,
				hashing:        mockHash,
			},
			args: args{
				ctx:     context.Background(),
				payload: *loginRequest,
			},
			want:    verifiedUser,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.name {
			case "error record not found":
				mockLogin(nil, errors.New("record not found"))(mockRepo)
			case "unexpected error":
				mockLogin(nil, errors.New("other error"))(mockRepo)
			case "error password incorrect":
				mockLogin(unverifiedUser, nil)(mockRepo)
				mockCompareHash(errors.New("any error"))(mockHash)
				mockSaveUser(nil)(mockRepo)
			case "error unverified account":
				mockLogin(unverifiedUser, nil)(mockRepo)
				mockCompareHash(nil)(mockHash)
			case "error password incorrect 3x":
				mockLogin(verifiedUserWrongPass3x, nil)(mockRepo)
			case "error saving wrong password counter":
				mockLogin(unverifiedUser, nil)(mockRepo)
				mockCompareHash(errors.New("any error"))(mockHash)
				mockSaveUser(errors.New("any error"))(mockRepo)
			case "error password incorrect 3x after wrong pass login":
				mockLogin(verifiedUserWrongPass2x, nil)(mockRepo)
				mockCompareHash(errors.New("any error"))(mockHash)
				mockSaveUser(nil)(mockRepo)
			case "error saving wrong password counter back to 0":
				mockLogin(verifiedUser, nil)(mockRepo)
				mockCompareHash(nil)(mockHash)
				mockSaveUser(errors.New("any error"))(mockRepo)
			case "success":
				mockLogin(verifiedUser, nil)(mockRepo)
				mockCompareHash(nil)(mockHash)
				mockSaveUser(nil)(mockRepo)
			}
			got, err := tt.s.Login(tt.args.ctx, tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.Login() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mock.NewMockUser(ctrl)
	mockHash := mockTools.NewMockBcryptInterface(ctrl)
	redis := mockRedis.NewMockRedisInterface(ctrl)

	req := &pb.UserRegisterRequest{
		Email:       "mail@mail.com",
		Password:    "pass",
		Name:        "name",
		PhoneNumber: "0123",
		Address:     "address",
		RoleId:      []string{"1"},
	}

	dataInactive := &entity.User{
		Email:      "mail@mail.com",
		IsVerified: false,
	}

	dataActive := &entity.User{
		Email:      "mail@mail.com",
		IsVerified: true,
	}

	type args struct {
		ctx context.Context
		req *pb.UserRegisterRequest
	}
	tests := []struct {
		name    string
		s       *Service
		args    args
		wantErr bool
		mocks   []*gomock.Call
	}{
		{
			name: "error hashing password",
			s: &Service{
				hashing: mockHash,
			},
			args: args{
				ctx: context.Background(),
				req: req,
			},
			wantErr: true,
			mocks: []*gomock.Call{
				mockHash.EXPECT().GenerateFromPassword(gomock.Any(), gomock.Any()).Return(nil, errors.New("error")),
			},
		},
		{
			name: "internal error",
			s: &Service{
				hashing:        mockHash,
				userRepository: mockRepo,
			},
			args: args{
				ctx: context.Background(),
				req: req,
			},
			wantErr: true,
			mocks: []*gomock.Call{
				mockHash.EXPECT().GenerateFromPassword(gomock.Any(), gomock.Any()).Return([]byte{}, nil),
				mockRepo.EXPECT().GetByEmail(gomock.Any(), gomock.Any()).Return(nil, errors.New("other error")),
			},
		},
		{
			name: "data exist with inactive status",
			s: &Service{
				userRepository: mockRepo,
				hashing:        mockHash,
			},
			args: args{
				ctx: context.Background(),
				req: req,
			},
			wantErr: true,
			mocks: []*gomock.Call{
				mockHash.EXPECT().GenerateFromPassword(gomock.Any(), gomock.Any()).Return([]byte{}, nil),
				mockRepo.EXPECT().GetByEmail(gomock.Any(), gomock.Any()).Return(dataInactive, nil),
			},
		},
		{
			name: "data exist with active status",
			s: &Service{
				userRepository: mockRepo,
				hashing:        mockHash,
			},
			args: args{
				ctx: context.Background(),
				req: req,
			},
			wantErr: true,
			mocks: []*gomock.Call{
				mockHash.EXPECT().GenerateFromPassword(gomock.Any(), gomock.Any()).Return([]byte{}, nil),
				mockRepo.EXPECT().GetByEmail(gomock.Any(), gomock.Any()).Return(dataActive, nil),
			},
		},
		{
			name: "error register from create in repository layer",
			s: &Service{
				userRepository: mockRepo,
				hashing:        mockHash,
			},
			args: args{
				ctx: context.Background(),
				req: req,
			},
			wantErr: true,
			mocks: []*gomock.Call{
				mockHash.EXPECT().GenerateFromPassword(gomock.Any(), gomock.Any()).Return([]byte{}, nil),
				mockRepo.EXPECT().GetByEmail(gomock.Any(), gomock.Any()).Return(nil, errors.New("record not found")),
				mockRepo.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("error")),
			},
		},
		{
			name: "success register user error in redis",
			s: &Service{
				userRepository: mockRepo,
				hashing:        mockHash,
				redis:          redis,
			},
			args: args{
				ctx: context.Background(),
				req: req,
			},
			wantErr: true,
			mocks: []*gomock.Call{
				mockHash.EXPECT().GenerateFromPassword(gomock.Any(), gomock.Any()).Return([]byte{}, nil),
				mockRepo.EXPECT().GetByEmail(gomock.Any(), gomock.Any()).Return(nil, errors.New("record not found")),
				mockRepo.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil),
				redis.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("any error")),
			},
		},
		{
			name: "success",
			s: &Service{
				userRepository: mockRepo,
				hashing:        mockHash,
				redis:          redis,
			},
			args: args{
				ctx: context.Background(),
				req: req,
			},
			wantErr: false,
			mocks: []*gomock.Call{
				mockHash.EXPECT().GenerateFromPassword(gomock.Any(), gomock.Any()).Return([]byte{}, nil),
				mockRepo.EXPECT().GetByEmail(gomock.Any(), gomock.Any()).Return(nil, errors.New("record not found")),
				mockRepo.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil),
				redis.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.s.Register(tt.args.ctx, tt.args.req)
			if err != nil != tt.wantErr {
				t.Errorf("Service.Register() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestService_CreateRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mock.NewMockRole(ctrl)
	mockCreate := func(err error) func(m *mock.MockRole) {
		return func(m *mock.MockRole) {
			m.EXPECT().Create(gomock.Any(), gomock.Any()).Return(err)
		}
	}
	type args struct {
		ctx  context.Context
		role *entity.Role
	}
	tests := []struct {
		name    string
		s       *Service
		args    args
		wantErr bool
	}{
		{
			name: "create role unit test",
			s: &Service{
				roleRepo: mockRepo,
			},
			args: args{
				ctx:  context.Background(),
				role: &entity.Role{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCreate(nil)(mockRepo)
			if err := tt.s.CreateRole(tt.args.ctx, tt.args.role); (err != nil) != tt.wantErr {
				t.Errorf("Service.CreateRole() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestService_GetRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mock.NewMockRole(ctrl)
	mockGetRole := func(data []entity.Role, err error) func(m *mock.MockRole) {
		return func(m *mock.MockRole) {
			m.EXPECT().GetRole(gomock.Any()).Return(data, err)
		}
	}
	data := []entity.Role{
		{
			RoleName: "Merchant",
		},
		{
			RoleName: "Customer",
		},
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		s       *Service
		args    args
		want    []entity.Role
		wantErr bool
	}{
		{
			name: "error get role from repository",
			s: &Service{
				roleRepo: mockRepo,
			},
			args: args{
				ctx: context.Background(),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "success get role",
			s: &Service{
				roleRepo: mockRepo,
			},
			args: args{
				ctx: context.Background(),
			},
			want:    data,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.name {
			case "error get role from repository":
				mockGetRole(nil, errors.New("any error"))(mockRepo)
			case "success get role":
				mockGetRole(data, nil)(mockRepo)
			}
			got, err := tt.s.GetRole(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.GetRole() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.GetRole() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_generateRandom4DigitNumber(t *testing.T) {
	tests := []struct {
		name string
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := generateRandom4DigitNumber(); got != tt.want {
				t.Errorf("generateRandom4DigitNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_VerifyOTP(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockUser := mock.NewMockUser(ctrl)
	redis := mockRedis.NewMockRedisInterface(ctrl)
	mockGetUser := func(user *entity.User, err error) func(m *mock.MockUser) {
		return func(m *mock.MockUser) {
			m.EXPECT().GetByEmail(gomock.Any(), gomock.Any()).Return(user, err)
		}
	}
	mockUpdateUser := func(res bool, err error) func(m *mock.MockUser) {
		return func(m *mock.MockUser) {
			m.EXPECT().VerifyUserByEmail(gomock.Any(), gomock.Any()).Return(res, err)
		}
	}
	mockGetStringKey := func(value string, err error) func(m *mockRedis.MockRedisInterface) {
		return func(m *mockRedis.MockRedisInterface) {
			m.EXPECT().GetStringKey(gomock.Any(), gomock.Any()).Return(value, err)
		}
	}
	mockGetContext := func(ctx context.Context) func(m *mockRedis.MockRedisInterface) {
		return func(m *mockRedis.MockRedisInterface) {
			m.EXPECT().GetContext().Return(ctx)
		}
	}
	id := uuid.New()
	verifiedUser := &entity.User{
		Id:         id,
		Email:      "test@mail.com",
		IsVerified: true,
	}
	unverifiedUser := &entity.User{
		Id:         id,
		Email:      "test@mail.com",
		IsVerified: false,
	}

	failedStoredJSON := `{OTP:error`
	succcessStoredJSON := `{"OTP":"1234"}`
	redisKey := "otp:" + "test@mail.com"

	type args struct {
		ctx      context.Context
		otp      int
		redisKey string
	}
	tests := []struct {
		name     string
		s        *Service
		args     args
		wantUser *entity.User
		wantErr  bool
	}{
		{
			name: "error verify otp caused by verified user",
			s: &Service{
				userRepository: mockUser,
				redis:          redis,
			},
			args: args{
				ctx:      context.Background(),
				otp:      generateRandom4DigitNumber(),
				redisKey: redisKey,
			},
			wantUser: nil,
			wantErr:  true,
		},
		{
			name: "error verify otp caused by no record",
			s: &Service{
				userRepository: mockUser,
				redis:          redis,
			},
			args: args{
				ctx:      context.Background(),
				otp:      generateRandom4DigitNumber(),
				redisKey: redisKey,
			},
			wantUser: nil,
			wantErr:  true,
		},
		{
			name: "error verify otp caused by redis nil",
			s: &Service{
				userRepository: mockUser,
				redis:          redis,
			},
			args: args{
				ctx:      context.Background(),
				otp:      generateRandom4DigitNumber(),
				redisKey: redisKey,
			},
			wantUser: nil,
			wantErr:  true,
		},
		{
			name: "error verify otp caused by other redis error",
			s: &Service{
				userRepository: mockUser,
				redis:          redis,
			},
			args: args{
				ctx:      context.Background(),
				otp:      generateRandom4DigitNumber(),
				redisKey: redisKey,
			},
			wantUser: nil,
			wantErr:  true,
		},
		{
			name: "error verify otp caused by unmarshal stored json",
			s: &Service{
				userRepository: mockUser,
				redis:          redis,
			},
			args: args{
				ctx:      context.Background(),
				otp:      generateRandom4DigitNumber(),
				redisKey: redisKey,
			},
			wantUser: nil,
			wantErr:  true,
		},
		{
			name: "error verify otp caused by incorrect input otp",
			s: &Service{
				userRepository: mockUser,
				redis:          redis,
			},
			args: args{
				ctx:      context.Background(),
				otp:      generateRandom4DigitNumber(),
				redisKey: redisKey,
			},
			wantUser: nil,
			wantErr:  true,
		},
		{
			name: "error verify otp caused by error saving user",
			s: &Service{
				userRepository: mockUser,
				redis:          redis,
			},
			args: args{
				ctx:      context.Background(),
				otp:      1234,
				redisKey: redisKey,
			},
			wantUser: nil,
			wantErr:  true,
		},
		{
			name: "success",
			s: &Service{
				userRepository: mockUser,
				redis:          redis,
			},
			args: args{
				ctx:      context.Background(),
				otp:      1234,
				redisKey: redisKey,
			},
			wantUser: unverifiedUser,
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.name {
			case "error verify otp caused by verified user":
				mockGetUser(verifiedUser, nil)(mockUser)
			case "error verify otp caused by no record":
				mockGetUser(nil, errors.New("any error"))(mockUser)
			case "error verify otp caused by redis nil":
				mockGetUser(unverifiedUser, nil)(mockUser)
				mockGetContext(context.Background())(redis)
				mockGetStringKey("", r.Nil)(redis)
			case "error verify otp caused by other redis error":
				mockGetUser(unverifiedUser, nil)(mockUser)
				mockGetContext(context.Background())(redis)
				mockGetStringKey("", errors.New("other error"))(redis)
			case "error verify otp caused by unmarshal stored json":
				mockGetUser(unverifiedUser, nil)(mockUser)
				mockGetContext(context.Background())(redis)
				mockGetStringKey(failedStoredJSON, nil)(redis)
			case "error verify otp caused by incorrect input otp":
				mockGetUser(unverifiedUser, nil)(mockUser)
				mockGetContext(context.Background())(redis)
				mockGetStringKey(succcessStoredJSON, nil)(redis)
			case "error verify otp caused by error saving user":
				mockGetUser(unverifiedUser, nil)(mockUser)
				mockGetContext(context.Background())(redis)
				mockGetStringKey(succcessStoredJSON, nil)(redis)
				mockUpdateUser(false, errors.New("any error"))(mockUser)
			case "success":
				mockGetUser(unverifiedUser, nil)(mockUser)
				mockGetContext(context.Background())(redis)
				mockGetStringKey(succcessStoredJSON, nil)(redis)
				mockUpdateUser(true, nil)(mockUser)
			}
			gotUser, err := tt.s.VerifyOTP(tt.args.ctx, tt.args.otp, tt.args.redisKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.VerifyOTP() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotUser, tt.wantUser) {
				t.Errorf("Service.VerifyOTP() = %v, want %v", gotUser, tt.wantUser)
			}
		})
	}
}

func TestService_ChangePassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockUser := mock.NewMockUser(ctrl)
	mockGetUser := func(user *entity.User, err error) func(m *mock.MockUser) {
		return func(m *mock.MockUser) {
			m.EXPECT().GetByEmail(gomock.Any(), gomock.Any()).Return(user, err)
		}
	}
	mockSaveUser := func(err error) func(m *mock.MockUser) {
		return func(m *mock.MockUser) {
			m.EXPECT().Save(gomock.Any(), gomock.Any()).Return(err)
		}
	}
	redis := mockRedis.NewMockRedisInterface(ctrl)
	mockGetStringKey := func(value string, err error) func(m *mockRedis.MockRedisInterface) {
		return func(m *mockRedis.MockRedisInterface) {
			m.EXPECT().GetStringKey(gomock.Any(), gomock.Any()).Return(value, err)
		}
	}
	mockGetContext := func(ctx context.Context) func(m *mockRedis.MockRedisInterface) {
		return func(m *mockRedis.MockRedisInterface) {
			m.EXPECT().GetContext().Return(ctx)
		}
	}
	mockHash := mockTools.NewMockBcryptInterface(ctrl)
	mockHashing := func(hashedPassword []byte, err error) func(m *mockTools.MockBcryptInterface) {
		return func(m *mockTools.MockBcryptInterface) {
			m.EXPECT().GenerateFromPassword(gomock.Any(), gomock.Any()).Return(hashedPassword, err)
		}
	}

	req := &pb.ChangePasswordRequest{
		Email:    "test@mail.com",
		Password: "password",
		OtpCode:  1234,
	}
	succcessStoredJSON := `{"OTP":"1234"}`
	user := &entity.User{
		Email:    "test@mail.com",
		Password: string([]byte{'a'}),
	}
	type args struct {
		ctx context.Context
		req *pb.ChangePasswordRequest
	}
	tests := []struct {
		name    string
		s       *Service
		args    args
		want    *entity.User
		wantErr bool
	}{
		{
			name: "error unregistered email",
			s: &Service{
				userRepository: mockUser,
			},
			args: args{
				ctx: context.Background(),
				req: req,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error repository issue",
			s: &Service{
				userRepository: mockUser,
			},
			args: args{
				ctx: context.Background(),
				req: req,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error verify otp",
			s: &Service{
				userRepository: mockUser,
				redis:          redis,
			},
			args: args{
				ctx: context.Background(),
				req: req,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error hashing password",
			s: &Service{
				userRepository: mockUser,
				redis:          redis,
				hashing:        mockHash,
			},
			args: args{
				ctx: context.Background(),
				req: req,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error update user data",
			s: &Service{
				userRepository: mockUser,
				redis:          redis,
				hashing:        mockHash,
			},
			args: args{
				ctx: context.Background(),
				req: req,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "success",
			s: &Service{
				userRepository: mockUser,
				redis:          redis,
				hashing:        mockHash,
			},
			args: args{
				ctx: context.Background(),
				req: req,
			},
			want:    user,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		switch tt.name {
		case "error unregistered email":
			mockGetUser(nil, errors.New("record not found"))(mockUser)
		case "error repository issue":
			mockGetUser(nil, errors.New("any error"))(mockUser)
		case "error verify otp":
			mockGetUser(user, nil)(mockUser)
			mockGetContext(context.Background())(redis)
			mockGetStringKey("", errors.New("any error"))(redis)
		case "error hashing password":
			mockGetUser(user, nil)(mockUser)
			mockGetContext(context.Background())(redis)
			mockGetStringKey(succcessStoredJSON, nil)(redis)
			mockHashing(nil, errors.New("any error"))(mockHash)
		case "error update user data":
			mockGetUser(user, nil)(mockUser)
			mockGetContext(context.Background())(redis)
			mockGetStringKey(succcessStoredJSON, nil)(redis)
			mockHashing([]byte{'a'}, nil)(mockHash)
			mockSaveUser(errors.New("any error"))(mockUser)
		case "success":
			mockGetUser(user, nil)(mockUser)
			mockGetContext(context.Background())(redis)
			mockGetStringKey(succcessStoredJSON, nil)(redis)
			mockHashing([]byte{'a'}, nil)(mockHash)
			mockSaveUser(nil)(mockUser)
		}
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.ChangePassword(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.ChangePassword() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.ChangePassword() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_verifyOtpFromRedis(t *testing.T) {
	type args struct {
		s        *Service
		otp      int
		redisKey string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := verifyOtpFromRedis(tt.args.s, tt.args.otp, tt.args.redisKey); (err != nil) != tt.wantErr {
				t.Errorf("verifyOtpFromRedis() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestService_ResendOTP(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockUser := mock.NewMockUser(ctrl)
	mockUserRecord := mockUser.EXPECT()
	redis := mockRedis.NewMockRedisInterface(ctrl)
	redisRecord := redis.EXPECT()
	email := "test@mail.com"
	user := &entity.User{
		Name:  "test",
		Email: email,
	}
	type args struct {
		ctx   context.Context
		email string
	}
	tests := []struct {
		name    string
		s       *Service
		args    args
		want    *entity.OtpMailReq
		wantErr bool
		mocks   []*gomock.Call
	}{
		{
			name: "error get by email repo",
			s: &Service{
				userRepository: mockUser,
			},
			args: args{
				ctx:   context.Background(),
				email: email,
			},
			want:    nil,
			wantErr: true,
			mocks: []*gomock.Call{
				mockUserRecord.GetByEmail(gomock.Any(), gomock.Any()).Return(nil, errors.New("any error")),
			},
		},
		{
			name: "error set key value in redis",
			s: &Service{
				userRepository: mockUser,
				redis:          redis,
			},
			args: args{
				ctx:   context.Background(),
				email: email,
			},
			want:    nil,
			wantErr: true,
			mocks: []*gomock.Call{
				mockUserRecord.GetByEmail(gomock.Any(), gomock.Any()).Return(user, nil),
				redisRecord.Set(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("any errror")),
			},
		},
		{
			name: "success",
			s: &Service{
				userRepository: mockUser,
				redis:          redis,
			},
			args: args{
				ctx:   context.Background(),
				email: email,
			},
			want: &entity.OtpMailReq{
				Name:  user.Name,
				Email: user.Email,
			},
			wantErr: false,
			mocks: []*gomock.Call{
				mockUserRecord.GetByEmail(gomock.Any(), gomock.Any()).Return(user, nil),
				redisRecord.Set(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.ResendOTP(tt.args.ctx, tt.args.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.ResendOTP() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil {
				if !reflect.DeepEqual(got.Name, tt.want.Name) {
					t.Errorf("Service.ResendOTP() = %v, want %v", got, tt.want)
				}
				if !reflect.DeepEqual(got.Email, tt.want.Email) {
					t.Errorf("Service.ResendOTP() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestService_Logout(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockUser := mock.NewMockUser(ctrl)
	ctx := context.Background()
	s := &Service{
		userRepository: mockUser,
	}
	logoutRequest := pb.LogoutRequest{}
	user := &entity.User{}

	t.Run("Should return no error", func(t *testing.T) {
		mockUser.EXPECT().GetByEmail(gomock.Any(), gomock.Any()).Return(user, nil)
		mockUser.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil)
		err := s.Logout(ctx, &logoutRequest)
		require.NoError(t, err)
	})

	t.Run("Should return error when email not found", func(t *testing.T) {
		mockUser.EXPECT().GetByEmail(gomock.Any(), gomock.Any()).Return(user, errors.New("not found"))
		err := s.Logout(ctx, &logoutRequest)
		require.Error(t, err)
	})

	t.Run("Should return error when find email error", func(t *testing.T) {
		mockUser.EXPECT().GetByEmail(gomock.Any(), gomock.Any()).Return(user, errors.New("failed"))
		err := s.Logout(ctx, &logoutRequest)
		require.Error(t, err)
	})

	t.Run("Should return error when saving error", func(t *testing.T) {
		mockUser.EXPECT().GetByEmail(gomock.Any(), gomock.Any()).Return(user, nil)
		mockUser.EXPECT().Save(gomock.Any(), gomock.Any()).Return(errors.New("failed"))
		err := s.Logout(ctx, &logoutRequest)
		require.Error(t, err)
	})

}
