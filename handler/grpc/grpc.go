package grpc

import (
	"context"
	"fmt"

	pb "github.com/Mitra-Apps/be-user-service/domain/proto/user"
	"github.com/Mitra-Apps/be-user-service/domain/user/entity"
	"github.com/Mitra-Apps/be-user-service/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/emptypb"
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

func (g *GrpcRoute) Register(ctx context.Context, req *pb.UserRegisterRequest) (*pb.SuccessResponse, error) {
	fmt.Println("test register grpc", req)
	if err := req.ValidateAll(); err != nil {
		return nil, err
	}
	if err := g.service.Register(ctx, req); err != nil {
		return nil, err
	}
	return &pb.SuccessResponse{}, nil
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
	roles, err := g.service.GetRole(ctx)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	protoRoles := []*pb.Role{}
	for _, r := range roles {
		protoRoles = append(protoRoles, r.ToProto())
	}
	roleData := &pb.ListRole{
		Roles: protoRoles,
	}
	data, err := anypb.New(roleData)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	return &pb.SuccessResponse{
		Data: data,
	}, nil
}
