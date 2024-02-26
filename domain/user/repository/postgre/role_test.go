package postgre

import (
	"context"
	"reflect"
	"testing"

	"github.com/Mitra-Apps/be-user-service/domain/user/entity"
	"github.com/Mitra-Apps/be-user-service/domain/user/repository"
	"gorm.io/gorm"
)

func TestNewRoleRepoImpl(t *testing.T) {
	type args struct {
		db *gorm.DB
	}
	tests := []struct {
		name string
		args args
		want repository.Role
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewRoleRepoImpl(tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRoleRepoImpl() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRoleRepoImpl_Create(t *testing.T) {
	type args struct {
		ctx  context.Context
		role *entity.Role
	}
	tests := []struct {
		name    string
		r       *RoleRepoImpl
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.r.Create(tt.args.ctx, tt.args.role); (err != nil) != tt.wantErr {
				t.Errorf("RoleRepoImpl.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRoleRepoImpl_GetRole(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		r       *RoleRepoImpl
		args    args
		want    []entity.Role
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.r.GetRole(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("RoleRepoImpl.GetRole() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RoleRepoImpl.GetRole() = %v, want %v", got, tt.want)
			}
		})
	}
}
