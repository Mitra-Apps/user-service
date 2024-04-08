package grpc

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/Mitra-Apps/be-user-service/config/postgre"
	pb "github.com/Mitra-Apps/be-user-service/domain/proto/user"
	"github.com/Mitra-Apps/be-user-service/domain/user/entity"
	userPostgreRepo "github.com/Mitra-Apps/be-user-service/domain/user/repository/postgre"
	"github.com/Mitra-Apps/be-user-service/external/redis"
	utilityservice "github.com/Mitra-Apps/be-user-service/external/utility_service"
	mockUtil "github.com/Mitra-Apps/be-user-service/external/utility_service/mock"

	mockRedis "github.com/Mitra-Apps/be-user-service/external/redis/mock"
	"github.com/Mitra-Apps/be-user-service/service"
	"github.com/Mitra-Apps/be-user-service/service/mock"
	utilPb "github.com/Mitra-Apps/be-utility-service/domain/proto/utility"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
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

func contextWithBearerToken(ctx context.Context, token string) context.Context {
	md := metadata.New(map[string]string{
		"authorization": "Bearer " + token,
	})
	return metadata.NewIncomingContext(ctx, md)

}

func TestNew(t *testing.T) {
	type args struct {
		service     service.ServiceInterface
		auth        service.Authentication
		utilService utilityservice.ServiceInterface
		redis       redis.RedisInterface
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
			if got := New(tt.args.service, tt.args.auth, tt.args.utilService, tt.args.redis); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGrpcRoute_GetUsers(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockSvc := mock.NewMockServiceInterface(ctrl)
	mockSvcRec := mockSvc.EXPECT()
	users := []*entity.User{
		{
			Name: "test",
		},
	}

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
		mock    *gomock.Call
	}{
		{
			name: "error get all users from repo",
			g: &GrpcRoute{
				service: mockSvc,
			},
			args: args{
				ctx: context.Background(),
				req: &pb.GetUsersRequest{},
			},
			want:    nil,
			wantErr: true,
			mock:    mockSvcRec.GetAll(gomock.Any()).Return(nil, errors.New("any error")),
		},
		{
			name: "success",
			g: &GrpcRoute{
				service: mockSvc,
			},
			args: args{
				ctx: context.Background(),
				req: &pb.GetUsersRequest{},
			},
			want: &pb.GetUsersResponse{
				Users: []*pb.User{
					{
						Id:   uuid.Nil.String(),
						Name: "test",
					},
				},
			},
			wantErr: false,
			mock:    mockSvcRec.GetAll(gomock.Any()).Return(users, nil),
		},
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
	mockSvcRecord := mockSvc.EXPECT()
	mockAuthRecord := mockAuth.EXPECT()
	req := &pb.UserLoginRequest{
		Email:    "test@mail.com",
		Password: "@Abc123",
	}
	userWithToken := &entity.User{
		Id:           uuid.New(),
		Email:        "test@mail.com",
		IsVerified:   true,
		AccessToken:  "AccessToken",
		RefreshToken: "RefreshToken",
	}
	userWithoutToken := &entity.User{
		Id:          uuid.New(),
		Email:       "test@mail.com",
		IsVerified:  true,
		AccessToken: "",
	}
	genToken := &entity.Token{
		AccessToken:  "AccessToken",
		RefreshToken: "RefreshToken",
	}

	token := map[string]interface{}{
		"access_token":  "AccessToken",
		"refresh_token": "RefreshToken",
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
		mocks   []*gomock.Call
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
			mocks: []*gomock.Call{
				mockSvc.EXPECT().Login(gomock.Any(), gomock.Any()).Return(nil, errors.New("any error")),
			},
		},
		{
			name: "user already have token in his db",
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
			mocks: []*gomock.Call{
				mockSvcRecord.Login(gomock.Any(), gomock.Any()).Return(userWithToken, nil),
				mockSvcRecord.Save(gomock.Any(), gomock.Any()).Return(errors.New("any error")),
			},
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
			mocks: []*gomock.Call{
				mockSvcRecord.Login(gomock.Any(), gomock.Any()).Return(userWithoutToken, nil),
				mockAuthRecord.GenerateToken(gomock.Any(), gomock.Any()).Return(nil, errors.New("any errors")),
			},
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
			mocks: []*gomock.Call{
				mockSvcRecord.Login(gomock.Any(), gomock.Any()).Return(userWithoutToken, nil),
				mockAuthRecord.GenerateToken(gomock.Any(), gomock.Any()).Return(genToken, nil),
				mockSvcRecord.Save(gomock.Any(), gomock.Any()).Return(nil),
			},
		},
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
	mockSvcRec := mockSvc.EXPECT()
	mockUtilSvc := mockUtil.NewMockServiceInterface(ctrl)
	mockUtilRec := mockUtilSvc.EXPECT()

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
		mocks   []*gomock.Call
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
					Email:       "test@mail.com",
					Password:    "password",
					Name:        "test",
					PhoneNumber: "0123456789",
					Address:     "address",
					RoleId:      []string{"1", "2"},
				},
			},
			want:    nil,
			wantErr: true,
			mocks: []*gomock.Call{
				mockSvcRec.Register(gomock.Any(), gomock.Any()).Return(nil, errors.New("any error")),
			},
		},
		{
			name: "error send otp",
			g: &GrpcRoute{
				service:     mockSvc,
				utilService: mockUtilSvc,
			},
			args: args{
				ctx: context.Background(),
				req: &pb.UserRegisterRequest{
					Email:       "test@mail.com",
					Password:    "password",
					Name:        "test",
					PhoneNumber: "0123456789",
					Address:     "address",
					RoleId:      []string{"1", "2"},
				},
			},
			want:    nil,
			wantErr: true,
			mocks: []*gomock.Call{
				mockSvcRec.Register(gomock.Any(), gomock.Any()).Return(&entity.OtpMailReq{}, nil),
				mockUtilRec.SendOtpMail(gomock.Any(), gomock.Any()).Return(nil, errors.New("any error")),
			},
		},
		{
			name: "success",
			g: &GrpcRoute{
				service:     mockSvc,
				utilService: mockUtilSvc,
			},
			args: args{
				ctx: context.Background(),
				req: &pb.UserRegisterRequest{
					Email:       "test@mail.com",
					Password:    "password",
					Name:        "test",
					PhoneNumber: "0123456789",
					Address:     "address",
					RoleId:      []string{"1", "2"},
				},
			},
			want: &pb.SuccessResponse{
				Code:    int32(codes.OK),
				Message: "success",
			},
			wantErr: false,
			mocks: []*gomock.Call{
				mockSvcRec.Register(gomock.Any(), gomock.Any()).Return(&entity.OtpMailReq{}, nil),
				mockUtilRec.SendOtpMail(gomock.Any(), gomock.Any()).Return(&utilPb.UtilSuccessResponse{Code: int32(codes.OK), Message: "success"}, nil),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
	usrSvc := service.New(usrRepo, roleRepo, nil, nil, nil)
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
	mockSvcRec := mockSvc.EXPECT()
	mockAuthRec := mockAuth.EXPECT()

	user := &entity.User{
		Id:         uuid.New(),
		Email:      "test@mail.com",
		IsVerified: true,
	}
	req := &pb.VerifyOTPRequest{
		Email:   "test@mail.com",
		OtpCode: 1234,
	}
	genToken := &entity.Token{
		AccessToken:  "AccessToken",
		RefreshToken: "RefreshToken",
	}

	token := map[string]interface{}{
		"access_token":  "AccessToken",
		"refresh_token": "RefreshToken",
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
		mocks   []*gomock.Call
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
			mocks: []*gomock.Call{
				mockSvcRec.VerifyOTP(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("any error")),
			},
		},
		{
			name: "error caused by generate token",
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
			mocks: []*gomock.Call{
				mockSvcRec.VerifyOTP(gomock.Any(), gomock.Any(), gomock.Any()).Return(user, nil),
				mockAuthRec.GenerateToken(gomock.Any(), gomock.Any()).Return(nil, errors.New("any error")),
			},
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
			mocks: []*gomock.Call{
				mockSvcRec.VerifyOTP(gomock.Any(), gomock.Any(), gomock.Any()).Return(user, nil),
				mockAuthRec.GenerateToken(gomock.Any(), gomock.Any()).Return(genToken, nil),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
	ctrl := gomock.NewController(t)
	mockSvc := mock.NewMockServiceInterface(ctrl)
	mockSvcRec := mockSvc.EXPECT()
	mockUtilSvc := mockUtil.NewMockServiceInterface(ctrl)
	mockUtilRec := mockUtilSvc.EXPECT()

	utilRes := &utilPb.UtilSuccessResponse{
		Code:    int32(codes.OK),
		Message: "success",
	}

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
		mocks   []*gomock.Call
	}{
		{
			name: "fail resend otp",
			g: &GrpcRoute{
				service: mockSvc,
			},
			args: args{
				ctx: context.Background(),
				req: &pb.ResendOTPRequest{
					Email: "test@mail.com",
				},
			},
			want:    nil,
			wantErr: true,
			mocks: []*gomock.Call{
				mockSvcRec.ResendOTP(gomock.Any(), gomock.Any()).Return(nil, errors.New("any error")),
			},
		},
		{
			name: "error send otp mail",
			g: &GrpcRoute{
				service:     mockSvc,
				utilService: mockUtilSvc,
			},
			args: args{
				ctx: context.Background(),
				req: &pb.ResendOTPRequest{
					Email: "test@mail.com",
				},
			},
			want:    nil,
			wantErr: true,
			mocks: []*gomock.Call{
				mockSvcRec.ResendOTP(gomock.Any(), gomock.Any()).Return(&entity.OtpMailReq{}, nil),
				mockUtilRec.SendOtpMail(gomock.Any(), gomock.Any()).Return(nil, errors.New("any error")),
			},
		},
		{
			name: "success",
			g: &GrpcRoute{
				service:     mockSvc,
				utilService: mockUtilSvc,
			},
			args: args{
				ctx: context.Background(),
				req: &pb.ResendOTPRequest{
					Email: "test@mail.com",
				},
			},
			want: &pb.SuccessResponse{
				Code:    int32(codes.OK),
				Message: "success",
			},
			wantErr: false,
			mocks: []*gomock.Call{
				mockSvcRec.ResendOTP(gomock.Any(), gomock.Any()).Return(&entity.OtpMailReq{}, nil),
				mockUtilRec.SendOtpMail(gomock.Any(), gomock.Any()).Return(utilRes, nil),
			},
		},
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
	mockGenerateToken := func(token *entity.Token, err error) func(m *mock.MockAuthentication) {
		return func(m *mock.MockAuthentication) {
			m.EXPECT().GenerateToken(gomock.Any(), gomock.Any()).Return(token, err)
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
	genToken := &entity.Token{
		AccessToken:  "accessToken",
		RefreshToken: "refreshToken",
	}

	token := map[string]interface{}{
		"access_token":  "accessToken",
		"refresh_token": "refreshToken",
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
			name: "error from generate token",
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
			case "error from generate token":
				mockChangePassword(user, nil)(mockSvc)
				mockGenerateToken(nil, errors.New("any error"))(mockAuth)
			case "success":
				mockChangePassword(user, nil)(mockSvc)
				mockGenerateToken(genToken, nil)(mockAuth)
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

func TestLogout(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockSvc := mock.NewMockServiceInterface(mockCtrl)
	mockSvc.EXPECT().Logout(gomock.Any(), gomock.Any()).Return(nil)
	reqBody := &emptypb.Empty{}
	server := &GrpcRoute{
		service: mockSvc,
	}

	t.Run("Should return 200", func(t *testing.T) {
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/users/logout", strings.NewReader(string(body)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		payload := &emptypb.Empty{}
		c := context.Background()
		_, err := server.Logout(c, payload)
		require.Nil(t, err)
		require.Equal(t, http.StatusOK, rec.Code)
	})
}

func TestGrpcRoute_RefreshToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockSvc := mock.NewMockServiceInterface(ctrl)
	mockAuth := mock.NewMockAuthentication(ctrl)
	mockAuthRec := mockAuth.EXPECT()
	mockSvcRec := mockSvc.EXPECT()

	userID := uuid.New()
	user := &entity.User{
		Id:   userID,
		Name: "test",
	}
	ctx := contextWithBearerToken(context.Background(), "any_token")

	genToken := &entity.Token{
		AccessToken:  "new_access_token",
		RefreshToken: "new_refresh_token",
	}
	token := map[string]interface{}{
		"access_token":  "new_access_token",
		"refresh_token": "new_refresh_token",
	}
	data, _ := structpb.NewStruct(token)

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
		mocks   []*gomock.Call
	}{
		{
			name: "error get user",
			g: &GrpcRoute{
				auth:    mockAuth,
				service: mockSvc,
			},
			args: args{
				ctx: ctx,
			},
			want:    nil,
			wantErr: true,
			mocks: []*gomock.Call{
				mockSvcRec.GetByID(ctx, gomock.Any()).Return(nil, errors.New("any error")),
			},
		},
		{
			name: "error generate new token",
			g: &GrpcRoute{
				auth:    mockAuth,
				service: mockSvc,
			},
			args: args{
				ctx: ctx,
			},
			want:    nil,
			wantErr: true,
			mocks: []*gomock.Call{
				mockSvcRec.GetByID(ctx, gomock.Any()).Return(user, nil),
				mockAuthRec.GenerateToken(ctx, gomock.Any()).Return(nil, errors.New("any error")),
			},
		},
		{
			name: "error update user",
			g: &GrpcRoute{
				auth:    mockAuth,
				service: mockSvc,
			},
			args: args{
				ctx: ctx,
			},
			want:    nil,
			wantErr: true,
			mocks: []*gomock.Call{
				mockSvcRec.GetByID(ctx, gomock.Any()).Return(user, nil),
				mockAuthRec.GenerateToken(ctx, gomock.Any()).Return(genToken, nil),
				mockSvcRec.Save(ctx, gomock.Any()).Return(errors.New("any error")),
			},
		},
		{
			name: "success",
			g: &GrpcRoute{
				auth:    mockAuth,
				service: mockSvc,
			},
			args: args{
				ctx: ctx,
			},
			want: &pb.SuccessResponse{
				Code:    int32(codes.OK),
				Message: "success",
				Data:    data,
			},
			wantErr: false,
			mocks: []*gomock.Call{
				mockSvcRec.GetByID(ctx, gomock.Any()).Return(user, nil),
				mockAuthRec.GenerateToken(ctx, gomock.Any()).Return(genToken, nil),
				mockSvcRec.Save(ctx, gomock.Any()).Return(nil),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.g.RefreshToken(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("GrpcRoute.RefreshToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GrpcRoute.RefreshToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGrpcRoute_SetEnvVariable(t *testing.T) {
	ctrl := gomock.NewController(t)
	redis := mockRedis.NewMockRedisInterface(ctrl)
	// mockSvc := mock.NewMockServiceInterface(ctrl)
	mockUtilSvc := mockUtil.NewMockServiceInterface(ctrl)

	// mockSvcRec := mockSvc.EXPECT()
	mockUtilSvcRec := mockUtilSvc.EXPECT()
	redisRec := redis.EXPECT()

	ctx := context.Background()
	req := &pb.EnvRequest{
		Variable: "var",
		Value:    "val",
	}
	type args struct {
		ctx context.Context
		req *pb.EnvRequest
	}
	tests := []struct {
		name    string
		g       *GrpcRoute
		args    args
		want    *pb.SuccessResponse
		wantErr bool
		mocks   []*gomock.Call
	}{
		{
			name: "error set variable",
			g: &GrpcRoute{
				utilService: mockUtilSvc,
			},
			args: args{
				ctx: ctx,
				req: req,
			},
			want:    nil,
			wantErr: true,
			mocks: []*gomock.Call{
				mockUtilSvcRec.UpsertEnvVariable(ctx, gomock.Any()).Return(nil, errors.New("any error")),
			},
		},
		{
			name: "error set value to redis",
			g: &GrpcRoute{
				utilService: mockUtilSvc,
				redis:       redis,
			},
			args: args{
				ctx: ctx,
				req: req,
			},
			want:    nil,
			wantErr: true,
			mocks: []*gomock.Call{
				mockUtilSvcRec.UpsertEnvVariable(ctx, gomock.Any()).Return(&utilPb.UtilSuccessResponse{Code: int32(codes.OK), Message: "success"}, nil),
				redisRec.Set(ctx, gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("any error")),
			},
		},
		{
			name: "success",
			g: &GrpcRoute{
				utilService: mockUtilSvc,
				redis:       redis,
			},
			args: args{
				ctx: ctx,
				req: req,
			},
			want: &pb.SuccessResponse{
				Code:    int32(codes.OK),
				Message: "success",
			},
			wantErr: false,
			mocks: []*gomock.Call{
				mockUtilSvcRec.UpsertEnvVariable(ctx, gomock.Any()).Return(&utilPb.UtilSuccessResponse{Code: int32(codes.OK), Message: "success"}, nil),
				redisRec.Set(ctx, gomock.Any(), gomock.Any(), gomock.Any()).Return(nil),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.g.SetEnvVariable(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("GrpcRoute.SetEnvVariable() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GrpcRoute.SetEnvVariable() = %v, want %v", got, tt.want)
			}
		})
	}
}
