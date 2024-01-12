package grpc

import (
	"context"

	pb "github.com/Mitra-Apps/be-user-service/domain/proto/user"
	"github.com/Mitra-Apps/be-user-service/domain/user/entity"
	"github.com/Mitra-Apps/be-user-service/service"
)

type GrpcRoute struct {
	service service.ServiceInterface
	pb.UnimplementedUserServiceServer
}

func New(service service.ServiceInterface) pb.UserServiceServer {
	return &GrpcRoute{
		service: service,
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

func (g *GrpcRoute) Login(ctx context.Context, req *pb.UserLoginRequest) (*pb.UserLoginResponse, error) {
	loginRequest := entity.LoginRequest{
		Username: req.Username,
		Password: req.Password,
	}
	user, err := g.service.Login(ctx, loginRequest)
	if err != nil {
		return nil, err
	}
	protoUser := user.ToProto()
	return &pb.UserLoginResponse{
		User: protoUser,
	}, nil
}

func (g *GrpcRoute) Register(ctx context.Context, req *pb.UserRegisterRequest) (*pb.SuccessMessage, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, err
	}
	if err := g.service.Register(ctx, req); err != nil {
		return nil, err
	}
	res := &pb.SuccessMessage{
		Message: "Akun berhasil di daftarkan",
	}

	return res, nil
}

func (g *GrpcRoute) CreateRole(ctx context.Context, req *pb.Role) (*pb.SuccessMessage, error) {
	role := &entity.Role{}
	if err := role.FromProto(req); err != nil {
		return nil, err
	}
	if err := g.service.CreateRole(ctx, role); err != nil {
		return nil, err
	}
	return &pb.SuccessMessage{
		Message: "Role successfully created",
	}, nil
}
