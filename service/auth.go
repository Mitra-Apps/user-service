package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
)

var (
	errClaimingToken  = errors.New("claim token error")
	errInvalidToken   = errors.New("invalid token error")
	errTokenExpired   = errors.New("token expired error")
	errUserIDRequired = errors.New("user id is required")
)

type authClient struct {
	secret string
}

//go:generate mockgen -source=auth.go -destination=mock/auth.go -package=mock
type Authentication interface {
	GenerateToken(ctx context.Context, id uuid.UUID, expiredMinute int) (token string, err error)
	ValidateToken(ctx context.Context, requestToken string) (uuid.UUID, error)
}

// Authentication client constructor
func NewAuthClient(secret string) *authClient {
	return &authClient{
		secret: secret,
	}
}

// CreateAccessToken will create access token that will be used for user authentication.
// access token will be needed in API that needs user to be authorized
func (c *authClient) GenerateToken(ctx context.Context, id uuid.UUID, expiredMinute int) (token string, err error) {
	if len(id) == 0 {
		return "", errUserIDRequired
	}

	currentTime := time.Now().UTC()
	//set token with criteria below and input userID into subject
	//this will be needed to check which user is this token for
	claims := &jwt.RegisteredClaims{
		Subject:   id.String(),
		ExpiresAt: jwt.NewNumericDate(currentTime.Add(time.Minute * time.Duration(expiredMinute))),
		IssuedAt:  jwt.NewNumericDate(currentTime),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err = t.SignedString([]byte(c.secret))
	if err != nil {
		return "", err
	}
	return token, err
}

// ValidateTokens will validate whether the token is valid
// and will return user id if the user is exist in our database
// otherwise it will return error
func (c *authClient) ValidateToken(ctx context.Context, requestToken string) (uuid.UUID, error) {
	token, err := jwt.Parse(requestToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(c.secret), nil
	})

	if err != nil {
		return uuid.Nil, NewError(codes.Unauthenticated, codes.Unauthenticated.String(), err.Error())
	}
	// assert jwt.MapClaims type
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return uuid.Nil, NewError(codes.Unauthenticated, codes.Unauthenticated.String(), errInvalidToken.Error())
	}

	currentTime := time.Now().UTC()
	expTime, err := claims.GetExpirationTime()
	if err != nil {
		return uuid.Nil, NewError(codes.Unauthenticated, codes.Unauthenticated.String(), errClaimingToken.Error())
	}

	if expTime.Before(currentTime) {
		return uuid.Nil, NewError(codes.Unauthenticated, codes.Unauthenticated.String(), errTokenExpired.Error())
	}

	//claim our user id input in subject from token
	id := claims["sub"].(string)
	var userID uuid.UUID
	userID, err = uuid.Parse(id)
	if err != nil {
		return uuid.Nil, NewError(codes.Unauthenticated, codes.Unauthenticated.String(), err.Error())
	}

	return userID, nil
}
