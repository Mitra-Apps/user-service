package service

import (
	"context"
	"errors"
	"reflect"
	"testing"

	mTools "github.com/Mitra-Apps/be-user-service/config/tools/mock"
	pb "github.com/Mitra-Apps/be-user-service/domain/proto/user"
	"github.com/Mitra-Apps/be-user-service/domain/user/entity"
	"github.com/Mitra-Apps/be-user-service/domain/user/repository/mock"
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
		Email:    "mail@mail.com",
		IsActive: false,
	}

	dataActive := &entity.User{
		Email:    "mail@mail.com",
		IsActive: true,
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
		// {
		// 	name: "success register user",
		// 	s: &Service{
		// 		userRepository: mockRepo,
		// 		hashing:        mockHash,
		// 	},
		// 	args: args{
		// 		ctx: context.Background(),
		// 		req: req,
		// 	},
		// 	wantErr: false,
		// },
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
			//TODO: check success flow after otp is done and update the unit test
			// case "success register user":
			// 	mockHashing([]byte{}, nil)(mockHash)
			// 	mockGetEmail(nil, errors.New("record not found"))(mockRepo)
			// 	mockRegister(nil)(mockRepo)
		}

		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.s.Register(tt.args.ctx, tt.args.req)
			if err != nil != tt.wantErr {
				t.Errorf("Service.Register() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_checkPassword(t *testing.T) {
	type args struct {
		password       string
		hashedPassword string
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
			if err := checkPassword(tt.args.password, tt.args.hashedPassword); (err != nil) != tt.wantErr {
				t.Errorf("checkPassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestService_CreateRole(t *testing.T) {
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.CreateRole(tt.args.ctx, tt.args.role); (err != nil) != tt.wantErr {
				t.Errorf("Service.CreateRole() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestService_GetRole(t *testing.T) {
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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

func TestService_generateUnique4DigitNumber(t *testing.T) {
	tests := []struct {
		name    string
		s       *Service
		want    int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.generateUnique4DigitNumber()
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.generateUnique4DigitNumber() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Service.generateUnique4DigitNumber() = %v, want %v", got, tt.want)
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
