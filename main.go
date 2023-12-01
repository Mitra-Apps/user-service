package main

import (
	"log"
	"net"
	"user-service/config/postgre"
	"user-service/domain/user/pb"
	userPostgreRepo "user-service/domain/user/repository/postgre"
	"user-service/lib"
	grpcRoute "user-service/route/grpc"
	"user-service/service"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedUserServiceServer
}

func main() {
	envInit()

	lis, err := net.Listen("tcp", lib.GetEnv("APP_PORT"))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	db := postgre.Connection()
	usrRepo := userPostgreRepo.NewPostgre(db)
	svc := service.New(usrRepo)
	grpcServer := grpc.NewServer()
	route := grpcRoute.New(svc)
	pb.RegisterUserServiceServer(grpcServer, route)

	log.Printf("GRPC Server listening at %v ", lis.Addr())
	if err = grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v \n", err)
	}
}

func envInit() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Fatalln("No .env file found")
	}
}
