package service

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/Mitra-Apps/be-user-service/domain/user/entity"
	"github.com/Mitra-Apps/be-user-service/domain/user/repository"
	mockRepo "github.com/Mitra-Apps/be-user-service/domain/user/repository/mock"
	"github.com/Mitra-Apps/be-user-service/external/redis"
	mockRedis "github.com/Mitra-Apps/be-user-service/external/redis/mock"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
)

func TestNewAuthClient(t *testing.T) {
	type args struct {
		secret   string
		redis    redis.RedisInterface
		userRepo repository.User
	}
	tests := []struct {
		name string
		args args
		want Authentication
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAuthClient(tt.args.secret, tt.args.redis, tt.args.userRepo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAuthClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_authClient_ValidateToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	redis := mockRedis.NewMockRedisInterface(ctrl)
	userRepo := mockRepo.NewMockUser(ctrl)
	auth := NewAuthClient("secret", redis, userRepo)
	user := &entity.User{
		Id: uuid.MustParse("b70a2a5e-bbd2-4000-96c0-aaa533b8236f"),
		Roles: []entity.Role{
			{
				RoleName: "merchant",
			},
			{
				RoleName: "customer",
			},
			{
				RoleName: "admin",
			},
		},
	}
	redis.EXPECT().GetStringKey(gomock.Any(), gomock.Any()).Return("60", nil)
	redis.EXPECT().GetStringKey(gomock.Any(), gomock.Any()).Return("43200", nil)
	token, err := auth.GenerateToken(context.Background(), user)
	if err != nil {
		panic(err.Error())
	}
	type args struct {
		ctx          context.Context
		requestToken string
	}
	tests := []struct {
		name    string
		c       *AuthClient
		args    args
		want    *JwtCustomClaim
		wantErr bool
	}{
		{
			name: "success",
			c: &AuthClient{
				secret: "secret",
				redis:  redis,
			},
			args: args{
				ctx:          context.Background(),
				requestToken: token.AccessToken,
			},
			want: &JwtCustomClaim{
				Roles: []string{
					"merchant",
					"customer",
					"admin",
				},
				RegisteredClaims: jwt.RegisteredClaims{
					Subject: "b70a2a5e-bbd2-4000-96c0-aaa533b8236f",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.c.ValidateToken(tt.args.ctx, tt.args.requestToken)
			if (err != nil) != tt.wantErr {
				t.Errorf("authClient.ValidateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil {
				if !reflect.DeepEqual(got.Roles, tt.want.Roles) {
					t.Errorf("authClient.ValidateToken() = %v, want %v", got, tt.want)
				}
				if !reflect.DeepEqual(got.Subject, tt.want.Subject) {
					t.Errorf("authClient.ValidateToken() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func Test_authClient_ValidateBlacklistToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockUserRepo := mockRepo.NewMockUser(ctrl)
	authSvc := NewAuthClient("secret", nil, mockUserRepo)

	user := entity.User{
		Id: uuid.New(),
	}
	params := entity.GetByTokensRequest{}

	type args struct {
		ctx    context.Context
		params *entity.GetByTokensRequest
	}
	tests := []struct {
		name     string
		mocks    []*gomock.Call
		args     args
		wantErr  bool
		wantResp bool
	}{
		{
			name: "should be return error because data not found",
			args: args{
				ctx:    context.Background(),
				params: &params,
			},
			mocks: []*gomock.Call{
				mockUserRepo.EXPECT().
					GetByTokens(gomock.Any(), &params).
					Times(1).
					Return(nil, nil),
			},
			wantErr:  true,
			wantResp: false,
		},
		{
			name: "should be return error because got error on repository",
			args: args{
				ctx:    context.Background(),
				params: &params,
			},
			mocks: []*gomock.Call{
				mockUserRepo.EXPECT().
					GetByTokens(gomock.Any(), &params).
					Times(1).
					Return(nil, errors.New("error")),
			},
			wantErr:  true,
			wantResp: false,
		},
		{
			name: "should be return nil (success)",
			args: args{
				ctx:    context.Background(),
				params: &params,
			},
			mocks: []*gomock.Call{
				mockUserRepo.EXPECT().
					GetByTokens(gomock.Any(), &params).
					Times(1).
					Return(&user, nil),
			},
			wantErr:  false,
			wantResp: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isValid, err := authSvc.IsTokenValid(tt.args.ctx, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("auth.ValidateBlacklistToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if isValid != tt.wantResp {
				t.Errorf("auth.ValidateBlacklistToken() wantResp %v", tt.wantResp)
				return
			}
		})
	}
}
