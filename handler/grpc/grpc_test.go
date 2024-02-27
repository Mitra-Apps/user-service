package grpc

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/Mitra-Apps/be-user-service/config/postgre"
	pb "github.com/Mitra-Apps/be-user-service/domain/proto/user"
	"github.com/Mitra-Apps/be-user-service/domain/user/entity"
	userPostgreRepo "github.com/Mitra-Apps/be-user-service/domain/user/repository/postgre"
	"github.com/Mitra-Apps/be-user-service/service"
	"github.com/Mitra-Apps/be-user-service/service/mock"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/structpb"
	"gorm.io/gorm"
)

func init() {
	if _, err := os.Stat("./../../.env"); !os.IsNotExist(err) {
		err := godotenv.Load(os.ExpandEnv("./../../.env"))
		if err != nil {
			log.Fatalf("Error getting env %v\n", err)
		}
	}
}

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
	ctrl := gomock.NewController(t)
	mockSvc := mock.NewMockServiceInterface(ctrl)
	mockAuth := mock.NewMockAuthentication(ctrl)
	mockLogin := func(user *entity.User, err error) func(m *mock.MockServiceInterface) {
		return func(m *mock.MockServiceInterface) {
			m.EXPECT().Login(gomock.Any(), gomock.Any()).Return(user, err)
		}
	}
	mockGenerateToken := func(token string, err error) func(m *mock.MockAuthentication) {
		return func(m *mock.MockAuthentication) {
			m.EXPECT().GenerateToken(gomock.Any(), gomock.Any(), gomock.Any()).Return(token, err)
		}
	}
	req := &pb.UserLoginRequest{
		Email:    "test@mail.com",
		Password: "@Abc123",
	}
	user := &entity.User{
		Id:         uuid.New(),
		Email:      "test@mail.com",
		IsVerified: true,
	}
	accessToken := "accessToken"
	refreshToken := "refreshToken"
	token := map[string]interface{}{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}
	data, _ := structpb.NewStruct(token)

	res := &pb.SuccessResponse{
		Data: data,
	}

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
		{
			name: "error validation",
			g: &GrpcRoute{
				service: mockSvc,
				auth:    mockAuth,
			},
			args: args{
				ctx: context.Background(),
				req: &pb.UserLoginRequest{
					Email:    "111",
					Password: "",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error caused by login service",
			g: &GrpcRoute{
				service: mockSvc,
				auth:    mockAuth,
			},
			args: args{
				ctx: context.Background(),
				req: req,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error caused by generate access token",
			g: &GrpcRoute{
				service: mockSvc,
				auth:    mockAuth,
			},
			args: args{
				ctx: context.Background(),
				req: req,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error caused by generate refresh token",
			g: &GrpcRoute{
				service: mockSvc,
				auth:    mockAuth,
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
			g: &GrpcRoute{
				service: mockSvc,
				auth:    mockAuth,
			},
			args: args{
				ctx: context.Background(),
				req: req,
			},
			want:    res,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.name {
			case "error caused by login service":
				mockLogin(nil, errors.New("any error"))(mockSvc)
			case "error caused by generate access token":
				mockLogin(user, nil)(mockSvc)
				mockGenerateToken("", errors.New("any error"))(mockAuth)
			case "error caused by generate refresh token":
				mockLogin(user, nil)(mockSvc)
				mockGenerateToken(accessToken, nil)(mockAuth)
				mockGenerateToken("", errors.New("any error"))(mockAuth)
			case "success":
				mockLogin(user, nil)(mockSvc)
				mockGenerateToken(accessToken, nil)(mockAuth)
				mockGenerateToken(refreshToken, nil)(mockAuth)
			}
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
	roles := []entity.Role{
		{
			Model:    gorm.Model{ID: 1},
			RoleName: "merchant",
		},
		{
			Model:    gorm.Model{ID: 2},
			RoleName: "customer",
		},
	}

	Roles := []*pb.Role{
		{
			Id:       "1",
			RoleName: "merchant",
		},
		{
			Id:       "2",
			RoleName: "customer",
		},
	}

	listStruct := map[string]interface{}{
		"roles": Roles,
	}

	data, _ := json.Marshal(listStruct)
	json.Unmarshal(data, &listStruct)
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
				Code:    0,
				Message: "roles data",
				Data:    pbData,
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
				mockGetRole(roles, nil)(mockSvc)
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

func TestGrpcRoute_GetRole_E2E(t *testing.T) {
	db := postgre.Connection()
	usrRepo := userPostgreRepo.NewUserRepoImpl(db)
	roleRepo := userPostgreRepo.NewRoleRepoImpl(db)
	usrSvc := service.New(usrRepo, roleRepo, nil, nil, nil, nil)
	permission := map[string]interface{}{
		"store": "create store",
	}
	data, _ := json.Marshal(permission)
	json.Unmarshal(data, &permission)
	permissionData, _ := structpb.NewStruct(permission)

	Roles := []*pb.Role{
		{
			Id:          "1",
			RoleName:    "merchant",
			Description: "role for merchant",
			Permission:  permissionData,
			IsActive:    true,
		},
		{
			Id:          "2",
			RoleName:    "customer",
			Description: "role for customer",
			IsActive:    true,
		},
		{
			Id:          "3",
			RoleName:    "admin",
			Description: "role for merchant",
			Permission:  permissionData,
			IsActive:    true,
		},
	}

	listStruct := map[string]interface{}{
		"roles": Roles,
	}

	data, _ = json.Marshal(listStruct)
	json.Unmarshal(data, &listStruct)
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
			name: "success",
			g: &GrpcRoute{
				service: usrSvc,
			},
			args: args{
				ctx: context.Background(),
				req: &emptypb.Empty{},
			},
			want: &pb.SuccessResponse{
				Code:    0,
				Message: "roles data",
				Data:    pbData,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.g.GetRole(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("GrpcRoute.GetRole() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil {
				if !reflect.DeepEqual(got.Code, tt.want.Code) {
					t.Errorf("GrpcRoute.GetRole() = %v, want %v", got, tt.want)
				}
				if !reflect.DeepEqual(got.Message, tt.want.Message) {
					t.Errorf("GrpcRoute.GetRole() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestGrpcRoute_VerifyOtp(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockSvc := mock.NewMockServiceInterface(ctrl)
	mockAuth := mock.NewMockAuthentication(ctrl)
	mockVerifyOtp := func(user *entity.User, err error) func(m *mock.MockServiceInterface) {
		return func(m *mock.MockServiceInterface) {
			m.EXPECT().VerifyOTP(gomock.Any(), gomock.Any(), gomock.Any()).Return(user, err)
		}
	}
	mockGenerateToken := func(token string, err error) func(m *mock.MockAuthentication) {
		return func(m *mock.MockAuthentication) {
			m.EXPECT().GenerateToken(gomock.Any(), gomock.Any(), gomock.Any()).Return(token, err)
		}
	}

	user := &entity.User{
		Id:         uuid.New(),
		Email:      "test@mail.com",
		IsVerified: true,
	}
	req := &pb.VerifyOTPRequest{
		Email:   "test@mail.com",
		OtpCode: 1234,
	}
	accessToken := "accessToken"
	refreshToken := "refreshToken"
	token := map[string]interface{}{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}
	data, _ := structpb.NewStruct(token)

	res := &pb.SuccessResponse{
		Data: data,
	}

	type args struct {
		ctx context.Context
		req *pb.VerifyOTPRequest
	}
	tests := []struct {
		name    string
		g       *GrpcRoute
		args    args
		want    *pb.SuccessResponse
		wantErr bool
	}{
		{
			name: "error caused by service verify otp",
			g: &GrpcRoute{
				service: mockSvc,
				auth:    mockAuth,
			},
			args: args{
				ctx: context.Background(),
				req: req,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error caused by generate access token",
			g: &GrpcRoute{
				service: mockSvc,
				auth:    mockAuth,
			},
			args: args{
				ctx: context.Background(),
				req: req,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error caused by generate refresh token",
			g: &GrpcRoute{
				service: mockSvc,
				auth:    mockAuth,
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
			g: &GrpcRoute{
				service: mockSvc,
				auth:    mockAuth,
			},
			args: args{
				ctx: context.Background(),
				req: req,
			},
			want:    res,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.name {
			case "error caused by service verify otp":
				mockVerifyOtp(nil, errors.New("any error"))(mockSvc)
			case "error caused by generate access token":
				mockVerifyOtp(user, nil)(mockSvc)
				mockGenerateToken("", errors.New("any error"))(mockAuth)
			case "error caused by generate refresh token":
				mockVerifyOtp(user, nil)(mockSvc)
				mockGenerateToken(accessToken, nil)(mockAuth)
				mockGenerateToken("", errors.New("any error"))(mockAuth)
			case "success":
				mockVerifyOtp(user, nil)(mockSvc)
				mockGenerateToken(accessToken, nil)(mockAuth)
				mockGenerateToken(refreshToken, nil)(mockAuth)
			}
			got, err := tt.g.VerifyOtp(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("GrpcRoute.VerifyOtp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GrpcRoute.VerifyOtp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGrpcRoute_ResendOtp(t *testing.T) {
	type args struct {
		ctx context.Context
		req *pb.ResendOTPRequest
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
			got, err := tt.g.ResendOtp(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("GrpcRoute.ResendOtp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GrpcRoute.ResendOtp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGrpcRoute_ChangePassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockSvc := mock.NewMockServiceInterface(ctrl)
	mockChangePassword := func(user *entity.User, err error) func(m *mock.MockServiceInterface) {
		return func(m *mock.MockServiceInterface) {
			m.EXPECT().ChangePassword(gomock.Any(), gomock.Any()).Return(user, err)
		}
	}
	mockAuth := mock.NewMockAuthentication(ctrl)
	mockGenerateToken := func(token string, err error) func(m *mock.MockAuthentication) {
		return func(m *mock.MockAuthentication) {
			m.EXPECT().GenerateToken(gomock.Any(), gomock.Any(), gomock.Any()).Return(token, err)
		}
	}
	req := &pb.ChangePasswordRequest{
		Email:    "test@mail.com",
		Password: "@Abc123",
	}
	user := &entity.User{
		Id:         uuid.New(),
		Email:      "test@mail.com",
		IsVerified: true,
	}
	accessToken := "accessToken"
	refreshToken := "refreshToken"
	token := map[string]interface{}{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}
	data, _ := structpb.NewStruct(token)

	res := &pb.SuccessResponse{
		Code:    int32(codes.OK),
		Message: "Sandi berhasil diubah!",
		Data:    data,
	}

	type args struct {
		ctx context.Context
		req *pb.ChangePasswordRequest
	}
	tests := []struct {
		name    string
		g       *GrpcRoute
		args    args
		want    *pb.SuccessResponse
		wantErr bool
	}{
		{
			name: "error validate proto",
			g: &GrpcRoute{
				service: mockSvc,
			},
			args: args{
				ctx: context.Background(),
				req: &pb.ChangePasswordRequest{
					Email: "a",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error from change password service",
			g: &GrpcRoute{
				service: mockSvc,
			},
			args: args{
				ctx: context.Background(),
				req: req,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error from generate access token",
			g: &GrpcRoute{
				service: mockSvc,
				auth:    mockAuth,
			},
			args: args{
				ctx: context.Background(),
				req: req,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error from generate refresh token",
			g: &GrpcRoute{
				service: mockSvc,
				auth:    mockAuth,
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
			g: &GrpcRoute{
				service: mockSvc,
				auth:    mockAuth,
			},
			args: args{
				ctx: context.Background(),
				req: req,
			},
			want:    res,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.name {
			case "error validate proto":
			case "error from change password service":
				mockChangePassword(nil, errors.New("any error"))(mockSvc)
			case "error from generate access token":
				mockChangePassword(user, nil)(mockSvc)
				mockGenerateToken("", errors.New("any error"))(mockAuth)
			case "error from generate refresh token":
				mockChangePassword(user, nil)(mockSvc)
				mockGenerateToken(accessToken, nil)(mockAuth)
				mockGenerateToken("", errors.New("any error"))(mockAuth)
			case "success":
				mockChangePassword(user, nil)(mockSvc)
				mockGenerateToken(accessToken, nil)(mockAuth)
				mockGenerateToken(refreshToken, nil)(mockAuth)
			}
			got, err := tt.g.ChangePassword(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("GrpcRoute.ChangePassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GrpcRoute.ChangePassword() = %v, want %v", got, tt.want)
			}
		})
	}
}
