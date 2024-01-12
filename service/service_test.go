package service

import (
	"reflect"
	"testing"

	"github.com/Mitra-Apps/be-user-service/domain/user/repository"
)

func TestNew(t *testing.T) {
	type args struct {
		userRepository repository.User
		roleRepo       repository.Role
	}
	tests := []struct {
		name string
		args args
		want *Service
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.userRepository, tt.args.roleRepo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}
