package grpc

import (
	"context"
	"user-service/domain/user/pb"
	"user-service/service"
)

type GrpcRoute struct {
	service service.ServiceInterface
}

func New(service service.ServiceInterface) *GrpcRoute {
	return &GrpcRoute{service}
}

func (g *GrpcRoute) GetUsers(ctx context.Context, req *pb.GetUserListRequest) (*pb.UserList, error) {
	users, err := g.service.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	list := pb.UserList{}
	list.List = make([]*pb.User, 0)
	for _, u := range users {
		list.List = append(list.List, &pb.User{
			Id:          u.Id,
			Username:    u.Username,
			Email:       u.Email,
			PhoneNumber: u.PhoneNumber,
			IsActive:    u.IsActive,
		})
	}
	return &list, nil
}
