package middleware

import (
	"context"
	"errors"
	"log"
	"strings"

	"github.com/Mitra-Apps/be-user-service/domain/user/entity"
	"github.com/Mitra-Apps/be-user-service/service"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type contextKey uint

const (
	userIdKey contextKey = 1
)

// Middleware interceptor
func JwtMiddlewareInterceptor(auth service.Authentication) grpc.UnaryServerInterceptor {

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Check if the method should be included from the middleware
		log.Print(info.FullMethod)

		addMiddleware := true
		isLogout := false

		// Add the method that will be included for middleware
		switch info.FullMethod {
		case "/proto.UserService/GetUsers":
			// Middleware logic for specific route
		case "/proto.UserService/GetOwnData":
			// Middleware logic for specific route
		case "/proto.UserService/Logout":
			// Middleware logic for specific route
			isLogout = true
		case "/proto.UserService/RefreshToken":
			// Middleware logic for specific route
		default:
			addMiddleware = false
		}

		if !addMiddleware {
			// Call the actual handler to process the request
			return handler(ctx, req)
		}

		// Validate and parse the JWT token
		token, err := GetToken(ctx)
		if err != nil {
			log.Println("error get token from middleware")
			return nil, err
		}

		// Validate token
		claims, err := auth.ValidateToken(ctx, token)
		if err != nil && !isLogout {
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

		// Validate access token is belong to user and still valid
		isValid, err := auth.IsTokenValid(ctx, &entity.GetByTokensRequest{
			Token:  token,
			UserId: userId,
		})

		if err != nil || !isValid {
			return nil, err
		}

		ctx = SetUserIDKey(ctx, userId)

		// Call the actual handler to process the request
		return handler(ctx, req)
	}
}

func GetToken(ctx context.Context) (token string, err error) {
	// Extract JWT token from metadata
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", errors.New("metadata is not provided")
	}

	authHeader, ok := md["authorization"]
	if !ok || len(authHeader) == 0 {
		return "", errors.New("authorization token is not provided")
	}

	token = strings.TrimPrefix(authHeader[0], "Bearer ")
	if token == authHeader[0] {
		return "", errors.New("invalid authorization format")
	}

	return token, nil
}

func SetUserIDKey(ctx context.Context, userId uuid.UUID) context.Context {
	return context.WithValue(ctx, userIdKey, userId)
}

func GetUserIDValue(ctx context.Context) uuid.UUID {
	if ctx == nil {
		return uuid.Nil
	}
	if userId, ok := ctx.Value(userIdKey).(uuid.UUID); ok {
		return userId
	}
	return uuid.Nil
}
