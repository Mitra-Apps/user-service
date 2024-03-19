package grpc

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	pb "github.com/Mitra-Apps/be-user-service/domain/proto/user"
	"github.com/Mitra-Apps/be-user-service/domain/user/entity"
	"github.com/Mitra-Apps/be-user-service/handler/middleware"
	"github.com/Mitra-Apps/be-user-service/service"
	utilPb "github.com/Mitra-Apps/be-utility-service/domain/proto/utility"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/structpb"
)

type GrpcRoute struct {
	service     service.ServiceInterface
	auth        service.Authentication
	utilService utilPb.MailServiceClient
	pb.UnimplementedUserServiceServer
}

func New(service service.ServiceInterface, auth service.Authentication, utilService utilPb.MailServiceClient) pb.UserServiceServer {
	return &GrpcRoute{
		service:     service,
		auth:        auth,
		utilService: utilService,
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
	user, err := g.service.Login(ctx, loginRequest)
	if err != nil {
		return nil, err
	}

	accessToken, err := g.auth.GenerateToken(ctx, user, 60)
	if err != nil {
		return nil, err
	}
	refreshToken, err := g.auth.GenerateToken(ctx, user, 43200)
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
	otpReq, err := g.service.Register(ctx, req)
	if err != nil {
		return nil, err
	}

	sendOtpReq := &utilPb.OtpMailReq{
		Name:    otpReq.Name,
		Email:   otpReq.Email,
		OtpCode: int32(otpReq.OtpCode),
	}

	//send otp to email
	res, err := g.utilService.SendOtpMail(ctx, sendOtpReq)
	if err != nil {
		return nil, err
	}

	return &pb.SuccessResponse{
		Code:    res.Code,
		Message: res.Message,
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
		Code:    int32(codes.OK),
		Message: "roles data",
		Data:    dataStruct,
	}, nil
}

func (g *GrpcRoute) VerifyOtp(ctx context.Context, req *pb.VerifyOTPRequest) (*pb.SuccessResponse, error) {
	redisKey := "otp:" + req.Email
	user, err := g.service.VerifyOTP(ctx, int(req.OtpCode), redisKey)
	if err != nil {
		return nil, err
	}
	accessToken, err := g.auth.GenerateToken(ctx, user, 60)
	if err != nil {
		return nil, err
	}
	refreshToken, err := g.auth.GenerateToken(ctx, user, 43200)
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

func (g *GrpcRoute) ResendOtp(ctx context.Context, req *pb.ResendOTPRequest) (*pb.SuccessResponse, error) {
	otpReq, err := g.service.ResendOTP(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	sendOtpReq := &utilPb.OtpMailReq{
		Name:    otpReq.Name,
		Email:   otpReq.Email,
		OtpCode: int32(otpReq.OtpCode),
	}

	//send otp to email
	res, err := g.utilService.SendOtpMail(ctx, sendOtpReq)
	if err != nil {
		return nil, err
	}

	return &pb.SuccessResponse{
		Code:    res.Code,
		Message: res.Message,
	}, nil
}

func (g *GrpcRoute) ChangePassword(ctx context.Context, req *pb.ChangePasswordRequest) (*pb.SuccessResponse, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, err
	}
	user, err := g.service.ChangePassword(ctx, req)
	if err != nil {
		return nil, err
	}
	accessToken, err := g.auth.GenerateToken(ctx, user, 60)
	if err != nil {
		return nil, err
	}
	refreshToken, err := g.auth.GenerateToken(ctx, user, 43200)
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
	res := &pb.SuccessResponse{
		Code:    int32(codes.OK),
		Message: "Sandi berhasil diubah!",
		Data:    data,
	}
	return res, nil
}

func (g *GrpcRoute) Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.SuccessResponse, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, err
	}
	go func() {
		ctx = context.Background()
		err := g.service.Logout(ctx, req)
		if err != nil {
			log.Print("Logout Error")
		}
	}()

	res := &pb.SuccessResponse{
		Code:    int32(codes.OK),
		Message: "Anda Berhasil Logout!",
	}
	return res, nil
}
