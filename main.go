package main

import (
	"log"
	"net"
	"user-service/domain/user/pb"
	"user-service/lib"
	grpcRoute "user-service/route/grpc"

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
	s := grpc.NewServer()
	route := grpcRoute.New()
	pb.RegisterUserServiceServer(s, route)

	log.Printf("GRPC Server listening at %v ", lis.Addr())
	if err = s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v \n", err)
	}
}

func envInit() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Fatalln("No .env file found")
	}
}
