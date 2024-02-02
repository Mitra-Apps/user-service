package grpc

import (
	"context"
	"errors"
	"reflect"
	"testing"

	pb "github.com/Mitra-Apps/be-user-service/domain/proto/user"
	"github.com/Mitra-Apps/be-user-service/domain/user/entity"
	"github.com/Mitra-Apps/be-user-service/service"
	"github.com/Mitra-Apps/be-user-service/service/mock"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/structpb"
	"gorm.io/gorm"
)

func TestNew(t *testing.T) {
	type args struct {
		service service.ServiceInterface
		auth    service.Authentication
	}
	ctrl := gomock.NewController(t)
	mockSvc := mock.NewMockServiceInterface(ctrl)
	mockAuth := mock.NewMockAuthentication(ctrl)

	tests := []struct {
		name string
		args args
		want pb.UserServiceServer
	}{
		{
			name: "implemented",
			args: args{
				service: mockSvc,
				auth:    mockAuth,
			},
			want: &GrpcRoute{
				service: mockSvc,
				auth:    mockAuth,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.service, tt.args.auth); !reflect.DeepEqual(got, tt.want) {
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
		want    *pb.SuccessResponse
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
	mockRegister := func(otp string, err error) func(m *mock.MockServiceInterface) {
		return func(m *mock.MockServiceInterface) {
			m.EXPECT().Register(gomock.Any(), gomock.Any()).Return(otp, err)
		}
	}
	otpStruct := map[string]interface{}{
		"otp": "",
	}
	data, _ := structpb.NewStruct(otpStruct)

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
					PhoneNumber: "0123456789",
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
					PhoneNumber: "0123456789",
					Address:     "address",
					RoleId:      []string{"1", "2"},
				},
			},
			want: &pb.SuccessResponse{
				Data: data,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.name {
			case "test fail register user from service layer":
				mockRegister("", errors.New("error"))(mockSvc)
			case "test success register user":
				mockRegister("", nil)(mockSvc)
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

func TestGrpcRoute_CreateRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockSvc := mock.NewMockServiceInterface(ctrl)
	mockCreateRole := func(err error) func(m *mock.MockServiceInterface) {
		return func(m *mock.MockServiceInterface) {
			m.EXPECT().CreateRole(gomock.Any(), gomock.Any()).Return(err)
		}
	}
	type args struct {
		ctx context.Context
		req *pb.Role
	}
	tests := []struct {
		name    string
		g       *GrpcRoute
		args    args
		want    *pb.SuccessResponse
		wantErr bool
	}{
		{
			name: "error parsing from proto",
			g:    &GrpcRoute{},
			args: args{
				ctx: context.Background(),
				req: &pb.Role{
					Id: "a",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error create role from service layer",
			g: &GrpcRoute{
				service: mockSvc,
			},
			args: args{
				ctx: context.Background(),
				req: &pb.Role{
					Id:          "1",
					RoleName:    "merchant",
					Description: "",
					IsActive:    true,
					Permission:  &structpb.Struct{},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "success",
			g: &GrpcRoute{
				service: mockSvc,
			},
			args: args{
				ctx: context.Background(),
				req: &pb.Role{
					Id:          "1",
					RoleName:    "merchant",
					Description: "",
					IsActive:    true,
					Permission:  &structpb.Struct{},
				},
			},
			want:    &pb.SuccessResponse{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.name {
			case "error create role from service layer":
				mockCreateRole(errors.New("any error"))(mockSvc)
			case "success":
				mockCreateRole(nil)(mockSvc)
			}
			got, err := tt.g.CreateRole(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("GrpcRoute.CreateRole() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GrpcRoute.CreateRole() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGrpcRoute_GetRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockSvc := mock.NewMockServiceInterface(ctrl)
	mockGetRole := func(data []entity.Role, err error) func(m *mock.MockServiceInterface) {
		return func(m *mock.MockServiceInterface) {
			m.EXPECT().GetRole(gomock.Any()).Return(data, err)
		}
	}
	data := []entity.Role{
		{
			Model:    gorm.Model{ID: 1},
			RoleName: "merchant",
		},
		{
			Model:    gorm.Model{ID: 2},
			RoleName: "customer",
		},
	}

	listRole := &pb.ListRole{
		Roles: []*pb.Role{
			{
				Id:       "1",
				RoleName: "merchant",
			},
			{
				Id:       "2",
				RoleName: "customer",
			},
		},
	}
	listStruct := map[string]interface{}{
		"roles": listRole,
	}
	pbData, _ := structpb.NewStruct(listStruct)
	type args struct {
		ctx context.Context
		req *emptypb.Empty
	}
	tests := []struct {
		name    string
		g       *GrpcRoute
		args    args
		want    *pb.SuccessResponse
		wantErr bool
	}{
		{
			name: "error get role from service layer",
			g: &GrpcRoute{
				service: mockSvc,
			},
			args: args{
				ctx: context.Background(),
				req: &emptypb.Empty{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "success",
			g: &GrpcRoute{
				service: mockSvc,
			},
			args: args{
				ctx: context.Background(),
				req: &emptypb.Empty{},
			},
			want: &pb.SuccessResponse{
				Data: pbData,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.name {
			case "error get role from service layer":
				mockGetRole(nil, errors.New("any error"))(mockSvc)
			case "success":
				mockGetRole(data, nil)(mockSvc)
			}
			got, err := tt.g.GetRole(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("GrpcRoute.GetRole() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GrpcRoute.GetRole() = %v, want %v", got, tt.want)
			}
		})
	}
}
