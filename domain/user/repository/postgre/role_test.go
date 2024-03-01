package postgre

import (
	"context"
	"log"
	"reflect"
	"testing"

	"github.com/Mitra-Apps/be-user-service/domain/user/entity"
	"github.com/Mitra-Apps/be-user-service/domain/user/repository"
	"gorm.io/gorm"
)

func TestNewRoleRepoImpl(t *testing.T) {
	db, err := DBConn()
	if err != nil {
		log.Fatal(err.Error())
	}
	type args struct {
		db *gorm.DB
	}
	tests := []struct {
		name string
		args args
		want repository.Role
	}{
		{
			name: "implemented",
			args: args{
				db: db,
			},
			want: &RoleRepoImpl{
				db: db,
			},
		},
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
	db, err := DBConn()
	if err != nil {
		log.Fatal(err.Error())
	}
	seedRole(db)
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
		{
			name: "success",
			r: &RoleRepoImpl{
				db: db,
			},
			args: args{
				ctx: context.Background(),
				role: &entity.Role{
					RoleName: "customer",
				},
			},
			wantErr: false,
		},
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
	db, err := DBConn()
	if err != nil {
		log.Fatal(err.Error())
	}
	role, err := seedRole(db)
	if err != nil {
		log.Fatal(err.Error())
	}
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
		{
			name: "success",
			r: &RoleRepoImpl{
				db: db,
			},
			args: args{
				ctx: context.Background(),
			},
			want:    *role,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.r.GetRole(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("RoleRepoImpl.GetRole() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil {
				for i, v := range got {
					if !reflect.DeepEqual(v.RoleName, tt.want[i].RoleName) {
						t.Errorf("RoleRepoImpl.GetRole() = %v, want %v", got, tt.want)
					}
				}
			}
		})
	}
}
