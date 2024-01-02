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
