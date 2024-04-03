package utilityservice

import (
	"context"
	"log"
	"os"

	pb "github.com/Mitra-Apps/be-utility-service/domain/proto/utility"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type serviceClient struct {
	client pb.UtilServiceClient
}

//go:generate mockgen -source=service.go -destination=mock/service.go -package=mock
type ServiceInterface interface {
	UpsertEnvVariable(ctx context.Context, req *pb.UpsertEnvVariableReq) (*pb.UtilSuccessResponse, error)
	SendOtpMail(ctx context.Context, req *pb.OtpMailReq) (*pb.UtilSuccessResponse, error)
	GetEnvVariable(ctx context.Context, req *pb.GetEnvVariableReq) (*pb.GetEnvVariableRes, error)
}

func NewClient(ctx context.Context) ServiceInterface {
	utilityGrpcConn, err := grpc.DialContext(ctx, os.Getenv("GRPC_UTILITY_HOST"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("Cannot connect to utility grpc server ", err)
	}
	defer func() {
		log.Println("Closing connection ...")
		utilityGrpcConn.Close()
	}()
	client := pb.NewUtilServiceClient(utilityGrpcConn)

	return &serviceClient{
		client: client,
	}
}

func (s *serviceClient) UpsertEnvVariable(ctx context.Context, req *pb.UpsertEnvVariableReq) (*pb.UtilSuccessResponse, error) {
	return s.client.UpsertEnvVariable(ctx, req)
}

func (s *serviceClient) SendOtpMail(ctx context.Context, req *pb.OtpMailReq) (*pb.UtilSuccessResponse, error) {
	return s.client.SendOtpMail(ctx, req)
}

func (s *serviceClient) GetEnvVariable(ctx context.Context, req *pb.GetEnvVariableReq) (*pb.GetEnvVariableRes, error) {
	return s.client.GetEnvVariable(ctx, req)
}
