package grpc

import (
	"context"
	"user-service/domain/user/pb"
)

type GrpcRoute struct {
}

func New() *GrpcRoute {
	return &GrpcRoute{}
}

func (g *GrpcRoute) GetUsers(ctx context.Context, req *pb.GetUserListRequest) (*pb.UserList, error) {
	list := pb.UserList{}
	list.List = make([]*pb.User, 0)
	list.List = append(list.List, &pb.User{
		Id:          1,
		Username:    "user",
		Email:       "user@mitrais.com",
		PhoneNumber: "0812933253243",
	})
	return &list, nil
}
