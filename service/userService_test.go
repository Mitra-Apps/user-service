package service

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/Mitra-Apps/be-user-service/config/redis"
	mTools "github.com/Mitra-Apps/be-user-service/config/tools/mock"
	pb "github.com/Mitra-Apps/be-user-service/domain/proto/user"
	"github.com/Mitra-Apps/be-user-service/domain/user/entity"
	"github.com/Mitra-Apps/be-user-service/domain/user/repository/mock"
	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
)

func TestService_GetAll(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		s       *Service
		args    args
		want    []*entity.User
		wantErr bool
	}{
		// TODO: Add test cases.
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
	mockHash := mTools.NewMockBcryptInterface(ctrl)
	mockCompareHash := func(err error) func(m *mTools.MockBcryptInterface) {
		return func(m *mTools.MockBcryptInterface) {
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
		Id:         userId,
		Email:      "test@email.com",
		Password:   "test@123",
		IsVerified: true,
	}
	type args struct {
		ctx     context.Context
		payload entity.LoginRequest
	}
	tests := []struct {
		name    string
		s       *Service
		args    args
		want    uuid.UUID
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
			want:    uuid.Nil,
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
			want:    uuid.Nil,
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
			want:    uuid.Nil,
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
			want:    uuid.Nil,
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
			want:    userId,
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
				mockLogin(verifiedUser, nil)(mockRepo)
				mockCompareHash(errors.New("any error"))(mockHash)
			case "error unverified account":
				mockLogin(unverifiedUser, nil)(mockRepo)
				mockCompareHash(nil)(mockHash)
			case "success":
				mockLogin(verifiedUser, nil)(mockRepo)
				mockCompareHash(nil)(mockHash)
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
	mockRegister := func(err error) func(m *mock.MockUser) {
		return func(m *mock.MockUser) {
			m.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Return(err)
		}
	}
	mockGetEmail := func(data *entity.User, err error) func(m *mock.MockUser) {
		return func(m *mock.MockUser) {
			m.EXPECT().GetByEmail(gomock.Any(), gomock.Any()).Return(data, err)
		}
	}
	mockHash := mTools.NewMockBcryptInterface(ctrl)
	mockHashing := func(hashedPassword []byte, err error) func(m *mTools.MockBcryptInterface) {
		return func(m *mTools.MockBcryptInterface) {
			m.EXPECT().GenerateFromPassword(gomock.Any(), gomock.Any()).Return(hashedPassword, err)
		}
	}

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
		},
		//TODO: check success flow after otp is done and update the unit test
		{
			name: "success register user",
			s: &Service{
				userRepository: mockRepo,
				hashing:        mockHash,
				redis:          redis.Connection(),
			},
			args: args{
				ctx: context.Background(),
				req: req,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		switch tt.name {
		case "error hashing password":
			mockHashing(nil, errors.New("error"))(mockHash)
		case "internal error":
			mockHashing([]byte{}, nil)(mockHash)
			mockGetEmail(nil, errors.New("other error"))(mockRepo)
		case "data exist with inactive status":
			mockHashing([]byte{}, nil)(mockHash)
			mockGetEmail(dataInactive, nil)(mockRepo)
		case "data exist with active status":
			mockHashing([]byte{}, nil)(mockHash)
			mockGetEmail(dataActive, nil)(mockRepo)
		case "error register from create in repository layer":
			mockHashing([]byte{}, nil)(mockHash)
			mockGetEmail(nil, errors.New("record not found"))(mockRepo)
			mockRegister(errors.New("error"))(mockRepo)
		case "success register user":
			mockHashing([]byte{}, nil)(mockHash)
			mockGetEmail(nil, errors.New("record not found"))(mockRepo)
			mockRegister(nil)(mockRepo)
		}

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
