package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/Mitra-Apps/be-user-service/config/postgre"
	pb "github.com/Mitra-Apps/be-user-service/domain/proto/user"
	userPostgreRepo "github.com/Mitra-Apps/be-user-service/domain/user/repository/postgre"
	grpcRoute "github.com/Mitra-Apps/be-user-service/handler/grpc"
	"github.com/Mitra-Apps/be-user-service/service"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedUserServiceServer
}

func main() {
	lis, err := net.Listen("tcp", os.Getenv("APP_PORT"))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	db := postgre.Connection()
	usrRepo := userPostgreRepo.NewPostgre(db)
	svc := service.New(usrRepo)
	grpcServer := grpc.NewServer()
	route := grpcRoute.New(svc)
	pb.RegisterUserServiceServer(grpcServer, route)

	fmt.Printf("GRPC Server listening at %v ", lis.Addr())
	if err = grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v \n", err)
	}
}
