package service

import (
	"reflect"
	"testing"

	"github.com/Mitra-Apps/be-user-service/domain/user/repository"
	"github.com/Mitra-Apps/be-user-service/domain/user/repository/mock"
	"github.com/Mitra-Apps/be-user-service/external"
	"github.com/Mitra-Apps/be-user-service/external/redis"

	"go.uber.org/mock/gomock"
)

func TestNew(t *testing.T) {
	type args struct {
		userRepository repository.User
		roleRepo       repository.Role
		hashing        external.BcryptInterface
		redis          redis.RedisInterface
		auth           Authentication
	}
	ctrl := gomock.NewController(t)
	mockUser := mock.NewMockUser(ctrl)
	mockRole := mock.NewMockRole(ctrl)
	tests := []struct {
		name string
		args args
		want *Service
	}{
		{
			name: "construct interface",
			args: args{
				userRepository: mockUser,
				roleRepo:       mockRole,
			},
			want: &Service{
				userRepository: mockUser,
				roleRepo:       mockRole,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.userRepository, tt.args.roleRepo, tt.args.hashing, tt.args.redis, tt.args.auth); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}
