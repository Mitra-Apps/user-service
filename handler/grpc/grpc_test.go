package grpc

import (
	"context"
	"errors"
	"reflect"
	"testing"

	pb "github.com/Mitra-Apps/be-user-service/domain/proto/user"
	"github.com/Mitra-Apps/be-user-service/service"
	"github.com/Mitra-Apps/be-user-service/service/mock"
	"go.uber.org/mock/gomock"
)

func TestNew(t *testing.T) {
	type args struct {
		service service.ServiceInterface
	}
	tests := []struct {
		name string
		args args
		want pb.UserServiceServer
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.service); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGrpcRoute_GetUsers(t *testing.T) {
	type args struct {
		ctx context.Context
		req *pb.GetUsersRequest
	}
	tests := []struct {
		name    string
		g       *GrpcRoute
		args    args
		want    *pb.GetUsersResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.g.GetUsers(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("GrpcRoute.GetUsers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GrpcRoute.GetUsers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGrpcRoute_Login(t *testing.T) {
	type args struct {
		ctx context.Context
		req *pb.UserLoginRequest
	}
	tests := []struct {
		name    string
		g       *GrpcRoute
		args    args
		want    *pb.UserLoginResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.g.Login(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("GrpcRoute.Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GrpcRoute.Login() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGrpcRoute_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockSvc := mock.NewMockServiceInterface(ctrl)
	mockRegister := func(err error) func(m *mock.MockServiceInterface) {
		return func(m *mock.MockServiceInterface) {
			m.EXPECT().Register(gomock.Any(), gomock.Any()).Return(err)
		}
	}
	type args struct {
		ctx context.Context
		req *pb.UserRegisterRequest
	}
	tests := []struct {
		name    string
		g       *GrpcRoute
		args    args
		want    *pb.SuccessResponse
		wantErr bool
	}{
		{
			name: "test error validating user",
			g: &GrpcRoute{
				service: mockSvc,
			},
			args: args{
				ctx: context.Background(),
				req: &pb.UserRegisterRequest{
					Email:       "email test",
					Password:    "password",
					Name:        "name",
					PhoneNumber: "0123",
					Address:     "address",
					RoleId:      []string{"1", "2"},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "test fail register user from service layer",
			g: &GrpcRoute{
				service: mockSvc,
			},
			args: args{
				ctx: context.Background(),
				req: &pb.UserRegisterRequest{
					Email:       "email@mail.com",
					Password:    "password",
					Name:        "name",
					PhoneNumber: "0123",
					Address:     "address",
					RoleId:      []string{"1", "2"},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "test success register user",
			g: &GrpcRoute{
				service: mockSvc,
			},
			args: args{
				ctx: context.Background(),
				req: &pb.UserRegisterRequest{
					Email:       "email@mail.com",
					Password:    "password",
					Name:        "name",
					PhoneNumber: "0123",
					Address:     "address",
					RoleId:      []string{"1", "2"},
				},
			},
			want:    &pb.SuccessResponse{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.name {
			case "test fail register user from service layer":
				mockRegister(errors.New("error"))(mockSvc)
			case "test success register user":
				mockRegister(nil)(mockSvc)
			}

			got, err := tt.g.Register(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("GrpcRoute.Register() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GrpcRoute.Register() = %v, want %v", got, tt.want)
			}
		})
	}
}
