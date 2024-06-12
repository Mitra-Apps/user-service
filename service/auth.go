package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	pbErr "github.com/Mitra-Apps/be-user-service/domain/proto"
	"github.com/Mitra-Apps/be-user-service/domain/user/entity"
	"github.com/Mitra-Apps/be-user-service/domain/user/repository"
	"github.com/Mitra-Apps/be-user-service/external/redis"
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
	secret         string
	redis          redis.RedisInterface
	userRepository repository.User
}

//go:generate mockgen -source=auth.go -destination=mock/auth.go -package=mock
type Authentication interface {
	GenerateToken(ctx context.Context, user *entity.User) (*entity.Token, error)
	ValidateToken(ctx context.Context, requestToken string) (*JwtCustomClaim, error)
	IsTokenValid(ctx context.Context, params *entity.GetByTokensRequest) (isValid bool, err error)
}

// Authentication client constructor
func NewAuthClient(secret string, redis redis.RedisInterface, userRepository repository.User) *authClient {
	return &authClient{
		secret:         secret,
		redis:          redis,
		userRepository: userRepository,
	}
}

// CreateAccessToken will create access token that will be used for user authentication.
// access token will be needed in API that needs user to be authorized
func (c *authClient) GenerateToken(ctx context.Context, user *entity.User) (*entity.Token, error) {
	if user.Id == uuid.Nil {
		return nil, errUserIDRequired
	}

	var roles []string
	for _, role := range user.Roles {
		roles = append(roles, role.RoleName)
	}

	currentTime := time.Now().UTC()

	//get access token exp time env var from redis
	redisAccessTokenValue, err := c.redis.GetStringKey(ctx, AccessTokenExpTime)
	if err != nil {
		redisAccessTokenValue = "60"
	}
	accessTokenExpTime, err := strconv.ParseFloat(redisAccessTokenValue, 64)
	if err != nil {
		return nil, util.NewError(codes.Internal, codes.Unknown.String(), err.Error())
	}

	//get access token exp time env var from redis
	redisRefreshTokenValue, err := c.redis.GetStringKey(ctx, RefreshTokenExpTime)
	if err != nil {
		redisRefreshTokenValue = "43200"
	}
	refreshTokenExpTime, err := strconv.ParseFloat(redisRefreshTokenValue, 64)
	if err != nil {
		return nil, util.NewError(codes.Internal, codes.Unknown.String(), err.Error())
	}

	accessTokenClaims := jwt.RegisteredClaims{
		Subject:   user.Id.String(),
		ExpiresAt: jwt.NewNumericDate(currentTime.Add(time.Minute * time.Duration(accessTokenExpTime))),
		IssuedAt:  jwt.NewNumericDate(currentTime),
		Issuer:    AccessToken,
	}

	refreshTokenClaims := jwt.RegisteredClaims{
		Subject:   user.Id.String(),
		ExpiresAt: jwt.NewNumericDate(currentTime.Add(time.Minute * time.Duration(refreshTokenExpTime))),
		IssuedAt:  jwt.NewNumericDate(currentTime),
		Issuer:    RefreshToken,
	}

	claims := &JwtCustomClaim{
		Roles:            roles,
		RegisteredClaims: accessTokenClaims,
	}

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)

	accessToken, err := at.SignedString([]byte(c.secret))
	if err != nil {
		return nil, err
	}
	refreshToken, err := rt.SignedString([]byte(c.secret))
	if err != nil {
		return nil, err
	}
	return &entity.Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, err
}

// ValidateTokens will validate whether the token is valid
// and will return user id if the user is exist in our database
// otherwise it will return error
func (c *authClient) ValidateToken(ctx context.Context, requestToken string) (*JwtCustomClaim, error) {
	token, _ := jwt.Parse(requestToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Println("error parse jwt")
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(c.secret), nil
	})

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		log.Println("token invalid, fail to claim token")
		return nil, util.NewError(codes.Unauthenticated, pbErr.ErrorCode_AUTH_JWT_ERR_GET_CLAIMS.String(), errInvalidToken.Error())
	}

	currentTime := time.Now().UTC()
	expTime, err := claims.GetExpirationTime()
	if err != nil {
		log.Println("error get exp time")
		return nil, util.NewError(codes.Unauthenticated, pbErr.ErrorCode_AUTH_JWT_ERR_GET_CLAIMS.String(), errClaimingToken.Error())
	}

	sub, err := claims.GetSubject()
	if err != nil {
		log.Println("error get subject")
		return nil, util.NewError(codes.Unauthenticated, pbErr.ErrorCode_AUTH_JWT_ERR_GET_CLAIMS.String(), errClaimingToken.Error())
	}

	issuer, err := claims.GetIssuer()
	if err != nil {
		log.Println("error get issuer")
		return nil, util.NewError(codes.Unauthenticated, pbErr.ErrorCode_AUTH_JWT_ERR_GET_CLAIMS.String(), errClaimingToken.Error())
	}

	iat, err := claims.GetIssuedAt()
	if err != nil {
		log.Println("error get issued at")
		return nil, util.NewError(codes.Unauthenticated, pbErr.ErrorCode_AUTH_JWT_ERR_GET_CLAIMS.String(), errClaimingToken.Error())
	}

	if expTime.Before(currentTime) {
		if issuer == AccessToken {
			ErrorCodeDetail = pbErr.ErrorCode_AUTH_ACCESS_TOKEN_EXPIRED.String()
		} else {
			ErrorCodeDetail = pbErr.ErrorCode_AUTH_REFRESH_TOKEN_EXPIRED.String()
		}
		return &JwtCustomClaim{
			RegisteredClaims: jwt.RegisteredClaims{
				Subject:   sub,
				ExpiresAt: expTime,
				Issuer:    issuer,
				IssuedAt:  iat,
			},
		}, util.NewError(codes.Unauthenticated, ErrorCodeDetail, errTokenExpired.Error())
	}

	var roles []string
	if claims["roles"] != nil {
		claimRoles := claims["roles"].([]interface{})
		for _, v := range claimRoles {
			roles = append(roles, v.(string))
		}
	}

	res := &JwtCustomClaim{
		Roles: roles,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   sub,
			ExpiresAt: expTime,
			IssuedAt:  iat,
			Issuer:    issuer,
		},
	}

	return res, nil
}

// To check token still on db or no
// if no, it means user already logout and token is not allowed to access
func (s *authClient) IsTokenValid(ctx context.Context, params *entity.GetByTokensRequest) (isValid bool, err error) {

	user, err := s.userRepository.GetByTokens(ctx, params)
	if err != nil {
		return
	}

	if user == nil {
		err = errTokenExpired
		return
	}

	isValid = true

	return
}
