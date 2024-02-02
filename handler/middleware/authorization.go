package middleware

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"
)

type contextKey uint

const (
	userIdKey contextKey = 1
)

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
