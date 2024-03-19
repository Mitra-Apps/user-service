package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/Mitra-Apps/be-user-service/config/postgre"
	"github.com/Mitra-Apps/be-user-service/config/tools"
	"github.com/Mitra-Apps/be-user-service/config/tools/redis"
	pb "github.com/Mitra-Apps/be-user-service/domain/proto/user"
	"github.com/Mitra-Apps/be-user-service/domain/user/entity"
	userPostgreRepo "github.com/Mitra-Apps/be-user-service/domain/user/repository/postgre"
	grpcRoute "github.com/Mitra-Apps/be-user-service/handler/grpc"
	"github.com/Mitra-Apps/be-user-service/handler/middleware"
	"github.com/Mitra-Apps/be-user-service/service"
	util "github.com/Mitra-Apps/be-utility-service/config/tools"
	utilPb "github.com/Mitra-Apps/be-utility-service/domain/proto/utility"
	"github.com/google/uuid"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"go.elastic.co/apm/module/apmgrpc"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

// Middleware interceptor
func middlewareInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// Check if the method should be included from the middleware
	log.Print(info.FullMethod)
	addMiddleware := true
	// Add the method that will be included for middleware
	switch info.FullMethod {
	case "/proto.UserService/GetUsers":
		// Middleware logic for specific route
	case "/proto.UserService/GetOwnData":
		// Middleware logic for specific route
	default:
		addMiddleware = false
	}
	if addMiddleware {
		// Validate and parse the JWT token
		token, err := middleware.GetToken(ctx)
		if err != nil {
			return nil, err
		}

		auth := service.NewAuthClient(os.Getenv("JWT_SECRET"))
		claims, err := auth.ValidateToken(ctx, token)
		if err != nil {
			return nil, err
		}

		//claim our user id input in subject from token
		id, err := claims.GetSubject()
		if err != nil {
			return nil, err
		}
		var userId uuid.UUID
		userId, err = uuid.Parse(id)
		if err != nil {
			return nil, err
		}

		ctx = middleware.SetUserIDKey(ctx, userId)
		// Call the actual handler to process the request
		return handler(ctx, req)
	}
	// Call the actual handler to process the request
	return handler(ctx, req)
}

func main() {
	ctx := context.Background()

	godotenv.Load()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", os.Getenv("GRPC_PORT")))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	utilityGrpcConn, err := grpc.DialContext(ctx, os.Getenv("GRPC_UTILITY_HOST"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("Cannot connect to utility grpc server ", err)
	}
	defer func() {
		log.Println("Closing connection ...")
		utilityGrpcConn.Close()
	}()
	mailSvcClient := utilPb.NewMailServiceClient(utilityGrpcConn)

	db := postgre.Connection()
	user := &entity.User{}
	err = db.Where("name = '1'").First(user).Error
	fmt.Println("not found", err)

	fmt.Println(user)
	redis := redis.Connection()
	usrRepo := userPostgreRepo.NewUserRepoImpl(db)
	roleRepo := userPostgreRepo.NewRoleRepoImpl(db)
	bcrypt := tools.New(&tools.Bcrypt{})
	auth := service.NewAuthClient(os.Getenv("JWT_SECRET"))
	svc := service.New(usrRepo, roleRepo, bcrypt, redis, auth)
	grpcServer := GrpcNewServer(ctx, []grpc.ServerOption{})
	route := grpcRoute.New(svc, auth, mailSvcClient)
	pb.RegisterUserServiceServer(grpcServer, route)

	go func() {
		<-ctx.Done()
		grpcServer.GracefulStop()
	}()

	go HttpNewServer(ctx, os.Getenv("GRPC_PORT"), os.Getenv("HTTP_PORT"))

	grpcServer.Serve(lis)
}

func GrpcNewServer(ctx context.Context, opts []grpc.ServerOption) *grpc.Server {
	logrusEntry := logrus.NewEntry(logrus.StandardLogger())
	logrusOpts := []grpc_logrus.Option{
		grpc_logrus.WithLevels(grpc_logrus.DefaultCodeToLevel),
	}
	grpc_logrus.ReplaceGrpcLogger(logrusEntry)

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
			middlewareInterceptor,
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
