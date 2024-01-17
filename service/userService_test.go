package service

import (
	"context"
	"errors"
	"reflect"
	"testing"

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

func TestService_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mock.NewMockUser(ctrl)
	mockRegister := func(err error) func(m *mock.MockUser) {
		return func(m *mock.MockUser) {
			m.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Return(err)
		}
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
			name: "error register from repository layer",
			s:    &Service{userRepository: mockRepo},
			args: args{
				ctx: context.Background(),
				req: &pb.UserRegisterRequest{
					Email:       "mail@mail.com",
					Password:    "pass",
					Name:        "name",
					PhoneNumber: "0123",
					Address:     "address",
					RoleId:      []string{"1"},
				},
			},
			wantErr: true,
		},
		{
			name: "success register user",
			s:    &Service{userRepository: mockRepo},
			args: args{
				ctx: context.Background(),
				req: &pb.UserRegisterRequest{
					Email:       "mail@mail.com",
					Password:    "pass",
					Name:        "name",
					PhoneNumber: "0123",
					Address:     "address",
					RoleId:      []string{"1"},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		switch tt.name {
		case "error register from repository layer":
			mockRegister(errors.New("error"))(mockRepo)
		case "success register user":
			mockRegister(nil)(mockRepo)
		}
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.Register(tt.args.ctx, tt.args.req); (err != nil) != tt.wantErr {
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

func TestService_Login(t *testing.T) {
	type args struct {
		ctx     context.Context
		payload entity.LoginRequest
	}
	tests := []struct {
		name    string
		args    args
		s       *Service
		want    *entity.User
		want1   []*entity.Role
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := tt.s.Login(tt.args.ctx, tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.Login() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Service.Login() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
