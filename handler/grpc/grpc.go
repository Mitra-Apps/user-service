package grpc

import (
	"context"
	"encoding/json"
	"fmt"

	pb "github.com/Mitra-Apps/be-user-service/domain/proto/user"
	"github.com/Mitra-Apps/be-user-service/domain/user/entity"
	"github.com/Mitra-Apps/be-user-service/handler/middleware"
	"github.com/Mitra-Apps/be-user-service/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/structpb"
)

type GrpcRoute struct {
	service service.ServiceInterface
	auth    service.Authentication
	pb.UnimplementedUserServiceServer
}

func New(service service.ServiceInterface, auth service.Authentication) pb.UserServiceServer {
	return &GrpcRoute{
		service: service,
		auth:    auth,
	}
}

func (g *GrpcRoute) GetUsers(ctx context.Context, req *pb.GetUsersRequest) (*pb.GetUsersResponse, error) {
	users, err := g.service.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	protoUsers := []*pb.User{}
	for _, user := range users {
		protoUsers = append(protoUsers, user.ToProto())
	}

	return &pb.GetUsersResponse{
		Users: protoUsers,
	}, nil
}

func (g *GrpcRoute) Login(ctx context.Context, req *pb.UserLoginRequest) (*pb.SuccessResponse, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, err
	}
	loginRequest := entity.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}
	userId, err := g.service.Login(ctx, loginRequest)
	if err != nil {
		return nil, err
	}

	accessToken, err := g.auth.GenerateToken(ctx, userId, 60)
	if err != nil {
		return nil, err
	}
	refreshToken, err := g.auth.GenerateToken(ctx, userId, 43200)
	if err != nil {
		return nil, err
	}

	token := map[string]interface{}{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}

	data, err := structpb.NewStruct(token)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	return &pb.SuccessResponse{
		Data: data,
	}, nil
}

func (g *GrpcRoute) Register(ctx context.Context, req *pb.UserRegisterRequest) (*pb.SuccessResponse, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, err
	}
	otp, err := g.service.Register(ctx, req)
	if err != nil {
		return nil, err
	}
	otpStruct := map[string]interface{}{
		"otp": otp,
	}
	data, err := structpb.NewStruct(otpStruct)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	return &pb.SuccessResponse{
		Data: data,
	}, nil
}

func (g *GrpcRoute) CreateRole(ctx context.Context, req *pb.Role) (*pb.SuccessResponse, error) {
	role := &entity.Role{}
	if err := role.FromProto(req); err != nil {
		return nil, err
	}
	if err := g.service.CreateRole(ctx, role); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}
	return &pb.SuccessResponse{}, nil
}

func (g *GrpcRoute) GetRole(ctx context.Context, req *emptypb.Empty) (*pb.SuccessResponse, error) {
	fmt.Println("get role handler", middleware.GetUserIDValue(ctx))
	roles, err := g.service.GetRole(ctx)
	if err != nil {
		fmt.Println("error 1", err.Error())
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	roleData := []*pb.Role{}
	for _, r := range roles {
		roleData = append(roleData, r.ToProto())
	}

	rolesStruct := map[string]interface{}{
		"roles": roleData,
	}

	data, err := json.Marshal(rolesStruct)
	if err != nil {
		fmt.Println("error 3", err.Error())
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	if err := json.Unmarshal(data, &rolesStruct); err != nil {
		fmt.Println("error 4", err.Error())
		return nil, err
	}

	dataStruct, err := structpb.NewStruct(rolesStruct)
	if err != nil {
		fmt.Println("error 5", err.Error())
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	return &pb.SuccessResponse{
		Data: dataStruct,
	}, nil
}

func (g *GrpcRoute) VerifyOtp(ctx context.Context, req *pb.VerifyOTPRequest) (*pb.SuccessResponse, error) {
	redisKey := "otp:" + req.Email
	_, err := g.service.VerifyOTP(ctx, int(req.OtpCode), redisKey)
	if err != nil {
		return nil, err
	}
	return &pb.SuccessResponse{}, nil
}

func (g *GrpcRoute) ResendOtp(ctx context.Context, req *pb.ResendOTPRequest) (*pb.SuccessResponse, error) {
	otp, err := g.service.ResendOTP(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	resendOtpStruct := map[string]interface{}{
		"otp": otp,
	}

	data, err := json.Marshal(resendOtpStruct)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, &resendOtpStruct); err != nil {
		return nil, err
	}
	dataStruct, err := structpb.NewStruct(resendOtpStruct)
	if err != nil {
		return nil, err
	}

	return &pb.SuccessResponse{
		Data: dataStruct,
	}, nil
}
