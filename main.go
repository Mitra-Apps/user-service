package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/Mitra-Apps/be-user-service/config/postgre"
	pb "github.com/Mitra-Apps/be-user-service/domain/proto/user"
	userPostgreRepo "github.com/Mitra-Apps/be-user-service/domain/user/repository/postgre"
	"github.com/Mitra-Apps/be-user-service/external"
	"github.com/Mitra-Apps/be-user-service/external/redis"
	utilityservice "github.com/Mitra-Apps/be-user-service/external/utility_service"
	grpcRoute "github.com/Mitra-Apps/be-user-service/handler/grpc"
	"github.com/Mitra-Apps/be-user-service/handler/middleware"
	"github.com/Mitra-Apps/be-user-service/service"
	util "github.com/Mitra-Apps/be-utility-service/config/tools"
	"github.com/Mitra-Apps/be-utility-service/domain/proto/utility"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"go.elastic.co/apm/module/apmgrpc"
	"gorm.io/gorm"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

func main() {
	ctx := context.Background()

	godotenv.Load()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", os.Getenv("GRPC_PORT")))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	db := postgre.Connection()
	redis := redis.Connection()
	utilSvc := utilityservice.NewClient(ctx)
	defer utilSvc.Close()

	SetRedisEnv(ctx, utilSvc, redis)

	grpcServer, route := SetgRPCRoute(ctx, utilSvc, redis, db)
	pb.RegisterUserServiceServer(grpcServer, route)

	go func() {
		<-ctx.Done()
		grpcServer.GracefulStop()
	}()

	go HttpNewServer(ctx, os.Getenv("GRPC_PORT"), os.Getenv("HTTP_PORT"))

	grpcServer.Serve(lis)
}

func GrpcNewServer(ctx context.Context, auth service.Authentication, opts []grpc.ServerOption) *grpc.Server {
	logrusEntry := logrus.NewEntry(logrus.StandardLogger())
	logrusOpts := []grpc_logrus.Option{
		grpc_logrus.WithLevels(grpc_logrus.DefaultCodeToLevel),
	}
	grpc_logrus.ReplaceGrpcLogger(logrusEntry)

	jwtMiddleware := middleware.JwtMiddlewareInterceptor(auth)

	opts = append(opts, grpc.StreamInterceptor(
		grpc_middleware.ChainStreamServer(
			grpc_ctxtags.StreamServerInterceptor(),
			grpc_logrus.StreamServerInterceptor(logrusEntry, logrusOpts...),
			grpc_recovery.StreamServerInterceptor(),
			apmgrpc.NewStreamServerInterceptor(apmgrpc.WithRecovery()),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_logrus.UnaryServerInterceptor(logrusEntry, logrusOpts...),
			grpc_recovery.UnaryServerInterceptor(),
			apmgrpc.NewUnaryServerInterceptor(apmgrpc.WithRecovery()),
			jwtMiddleware,
		)),
	)

	myServer := grpc.NewServer(opts...)

	reflection.Register(myServer)
	return myServer
}

func HttpNewServer(ctx context.Context, grpcPort, httpPort string) error {
	mux := runtime.NewServeMux(runtime.WithErrorHandler(util.CustomErrorHandler))

	mux.HandlePath("GET", "/docs/v1/users/openapi.yaml", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		http.ServeFile(w, r, "docs/openapi.yaml")
	})

	mux.HandlePath("GET", "/docs/v1/users", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		http.ServeFile(w, r, "docs/index.html")
	})

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	if err := pb.RegisterUserServiceHandlerFromEndpoint(ctx, mux, fmt.Sprintf("localhost:%s", grpcPort), opts); err != nil {
		return err
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", httpPort),
		Handler: mux,
	}

	go func() {
		<-ctx.Done()
		if err := srv.Shutdown(ctx); err != nil {
			logrus.Panicf("failed to shutdown server: %v", err)
		}
	}()

	return srv.ListenAndServe()
}

func SetRedisEnv(ctx context.Context, utilSvc utilityservice.ServiceInterface, redis redis.RedisInterface) {
	//setup access token exp time
	accessTokenExpTimeVal, err := utilSvc.GetEnvVariable(ctx, &utility.GetEnvVariableReq{Variable: service.AccessTokenExpTime})
	if err != nil {
		accessTokenExpTimeVal = &utility.GetEnvVariableRes{
			Value: "60",
		}
	}
	redis.Set(ctx, service.AccessTokenExpTime, accessTokenExpTimeVal.Value, time.Hour*time.Duration(720))

	// setup refresh token exp time
	refreshTokenExpTimeVal, err := utilSvc.GetEnvVariable(ctx, &utility.GetEnvVariableReq{Variable: service.RefreshTokenExpTime})
	if err != nil {
		refreshTokenExpTimeVal = &utility.GetEnvVariableRes{
			Value: "43200",
		}
	}
	redis.Set(ctx, service.RefreshTokenExpTime, refreshTokenExpTimeVal.Value, time.Hour*time.Duration(720))
}

func SetgRPCRoute(ctx context.Context, utilSvc utilityservice.ServiceInterface, redis redis.RedisInterface, db *gorm.DB) (*grpc.Server, pb.UserServiceServer) {

	usrRepo := userPostgreRepo.NewUserRepoImpl(db)
	roleRepo := userPostgreRepo.NewRoleRepoImpl(db)
	bcrypt := external.New(&external.Bcrypt{})
	auth := service.NewAuthClient(os.Getenv("JWT_SECRET"), redis, usrRepo)
	svc := service.New(usrRepo, roleRepo, bcrypt, redis, auth)

	grpcServer := GrpcNewServer(ctx, auth, []grpc.ServerOption{})
	route := grpcRoute.New(svc, auth, utilSvc, redis)

	return grpcServer, route
}
