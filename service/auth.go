package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Mitra-Apps/be-user-service/domain/user/entity"
	util "github.com/Mitra-Apps/be-utility-service/service"
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

type JwtCustomClaim struct {
	Roles []string `json:"roles"`
	jwt.RegisteredClaims
}

type authClient struct {
	secret string
}

//go:generate mockgen -source=auth.go -destination=mock/auth.go -package=mock
type Authentication interface {
	GenerateToken(ctx context.Context, user *entity.User, expiredMinute int) (token string, err error)
	ValidateToken(ctx context.Context, requestToken string) (*JwtCustomClaim, error)
}

// Authentication client constructor
func NewAuthClient(secret string) *authClient {
	return &authClient{
		secret: secret,
	}
}

// CreateAccessToken will create access token that will be used for user authentication.
// access token will be needed in API that needs user to be authorized
func (c *authClient) GenerateToken(ctx context.Context, user *entity.User, expiredMinute int) (token string, err error) {
	if user.Id == uuid.Nil {
		return "", errUserIDRequired
	}

	var roles []string
	for _, role := range user.Roles {
		roles = append(roles, role.RoleName)
	}

	currentTime := time.Now().UTC()
	//set token with criteria below and input userID into subject
	//this will be needed to check which user is this token for
	registeredClaims := jwt.RegisteredClaims{
		Subject:   user.Id.String(),
		ExpiresAt: jwt.NewNumericDate(currentTime.Add(time.Minute * time.Duration(expiredMinute))),
		IssuedAt:  jwt.NewNumericDate(currentTime),
	}

	claims := &JwtCustomClaim{
		Roles:            roles,
		RegisteredClaims: registeredClaims,
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
func (c *authClient) ValidateToken(ctx context.Context, requestToken string) (*JwtCustomClaim, error) {
	token, err := jwt.Parse(requestToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(c.secret), nil
	})

	if err != nil {
		return nil, util.NewError(codes.Unauthenticated, codes.Unauthenticated.String(), err.Error())
	}
	// assert jwt.MapClaims type
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, util.NewError(codes.Unauthenticated, codes.Unauthenticated.String(), errInvalidToken.Error())
	}

	currentTime := time.Now().UTC()
	expTime, err := claims.GetExpirationTime()
	if err != nil {
		return nil, util.NewError(codes.Unauthenticated, codes.Unauthenticated.String(), errClaimingToken.Error())
	}

	if expTime.Before(currentTime) {
		return nil, util.NewError(codes.Unauthenticated, codes.Unauthenticated.String(), errTokenExpired.Error())
	}

	sub, err := claims.GetSubject()
	if err != nil {
		return nil, util.NewError(codes.Unauthenticated, codes.Unauthenticated.String(), errClaimingToken.Error())
	}
	iat, err := claims.GetIssuedAt()
	if err != nil {
		return nil, util.NewError(codes.Unauthenticated, codes.Unauthenticated.String(), errClaimingToken.Error())
	}
	var roles []string
	claimRoles := claims["roles"].([]interface{})
	for _, v := range claimRoles {
		roles = append(roles, v.(string))
	}

	res := &JwtCustomClaim{
		Roles: roles,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   sub,
			ExpiresAt: expTime,
			IssuedAt:  iat,
		},
	}

	return res, nil
}
