package postgre

import (
	"context"
	"log"
	"reflect"
	"testing"
	"time"

	"github.com/Mitra-Apps/be-user-service/domain/user/entity"
	"github.com/Mitra-Apps/be-user-service/domain/user/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func TestNewUserRepoImpl(t *testing.T) {
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
		want repository.User
	}{
		{
			name: "implemented",
			args: args{
				db: db,
			},
			want: &userRepoImpl{
				db: db,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUserRepoImpl(tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUserRepoImpl() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_userRepoImpl_GetAll(t *testing.T) {
	db, err := DBConn()
	if err != nil {
		log.Fatal(err.Error())
	}
	user, err := seedUser(db)
	if err != nil {
		log.Fatal(err.Error())
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		p       *userRepoImpl
		args    args
		want    []*entity.User
		wantErr bool
	}{
		{
			name: "success",
			p: &userRepoImpl{
				db: db,
			},
			args: args{
				ctx: context.Background(),
			},
			want:    user,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.p.GetAll(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("userRepoImpl.GetAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("userRepoImpl.GetAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_userRepoImpl_GetByEmail(t *testing.T) {
	db, err := DBConn()
	if err != nil {
		log.Fatal(err.Error())
	}
	user, err := seedUser(db)
	if err != nil {
		log.Fatal(err.Error())
	}
	type args struct {
		ctx   context.Context
		email string
	}
	tests := []struct {
		name    string
		p       *userRepoImpl
		args    args
		want    *entity.User
		wantErr bool
	}{
		{
			name: "error get email",
			p: &userRepoImpl{
				db: db,
			},
			args: args{
				ctx:   context.Background(),
				email: "fail@mail.com",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "success",
			p: &userRepoImpl{
				db: db,
			},
			args: args{
				ctx:   context.Background(),
				email: "test1@mail.com",
			},
			want:    user[0],
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.p.GetByEmail(tt.args.ctx, tt.args.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("userRepoImpl.GetByEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil {
				if !reflect.DeepEqual(got.Id, tt.want.Id) {
					t.Errorf("userRepoImpl.GetByEmail() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func Test_userRepoImpl_GetByID(t *testing.T) {
	db, err := DBConn()
	if err != nil {
		log.Fatal(err.Error())
	}
	users, err := seedUser(db)
	if err != nil {
		log.Fatal(err.Error())
	}
	type args struct {
		ctx context.Context
		id  uuid.UUID
	}
	tests := []struct {
		name    string
		p       *userRepoImpl
		args    args
		want    *entity.User
		wantErr bool
	}{
		{
			name: "error record not found",
			p: &userRepoImpl{
				db: db,
			},
			args: args{
				ctx: context.Background(),
				id:  uuid.Nil,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "success",
			p: &userRepoImpl{
				db: db,
			},
			args: args{
				ctx: context.Background(),
				id:  users[0].Id,
			},
			want:    users[0],
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.p.GetByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("userRepoImpl.GetByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil {
				if !reflect.DeepEqual(got.Id, tt.want.Id) {
					t.Errorf("userRepoImpl.GetByID() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func Test_userRepoImpl_Create(t *testing.T) {
	db, err := DBConn()
	if err != nil {
		log.Fatal(err.Error())
	}
	seedUser(db)
	seedRole(db)
	type args struct {
		ctx     context.Context
		user    *entity.User
		roleIds []string
	}
	tests := []struct {
		name    string
		p       *userRepoImpl
		args    args
		wantErr bool
	}{
		{
			name: "error create in db",
			p: &userRepoImpl{
				db: db,
			},
			args: args{
				ctx: context.Background(),
				user: &entity.User{
					Name:     "name1",
					Email:    "test1@mail.com",
					Password: "password1",
				},
			},
			wantErr: true,
		},
		{
			name: "error insert role ids in db",
			p: &userRepoImpl{
				db: db,
			},
			args: args{
				ctx: context.Background(),
				user: &entity.User{
					Name:        "name1",
					Username:    "test2@mail.com",
					Email:       "test2@mail.com",
					Password:    "password1",
					PhoneNumber: "0123",
				},
				roleIds: []string{"3", "4"},
			},
			wantErr: true,
		},
		{
			name: "success",
			p: &userRepoImpl{
				db: db,
			},
			args: args{
				ctx: context.Background(),
				user: &entity.User{
					Name:        "name1",
					Username:    "test2@mail.com",
					Email:       "test2@mail.com",
					Password:    "password1",
					PhoneNumber: "0123",
				},
				roleIds: []string{"1", "2"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.p.Create(tt.args.ctx, tt.args.user, tt.args.roleIds); (err != nil) != tt.wantErr {
				t.Errorf("userRepoImpl.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_userRepoImpl_Save(t *testing.T) {
	db, err := DBConn()
	if err != nil {
		log.Fatal(err.Error())
	}
	user, err := seedUser(db)
	if err != nil {
		log.Fatal(err.Error())
	}
	user[0].UpdatedAt = time.Now()
	type args struct {
		ctx  context.Context
		user *entity.User
	}
	tests := []struct {
		name    string
		p       *userRepoImpl
		args    args
		wantErr bool
	}{
		{
			name: "success",
			p: &userRepoImpl{
				db: db,
			},
			args: args{
				ctx:  context.Background(),
				user: user[0],
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.p.Save(tt.args.ctx, tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("userRepoImpl.Save() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_userRepoImpl_VerifyUserByEmail(t *testing.T) {
	db, err := DBConn()
	if err != nil {
		log.Fatal(err.Error())
	}
	user, err := seedUser(db)
	if err != nil {
		log.Fatal(err.Error())
	}
	type args struct {
		ctx   context.Context
		email string
	}
	tests := []struct {
		name    string
		p       *userRepoImpl
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "success",
			p: &userRepoImpl{
				db: db,
			},
			args: args{
				ctx:   context.Background(),
				email: user[0].Email,
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.p.VerifyUserByEmail(tt.args.ctx, tt.args.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("userRepoImpl.VerifyUserByEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("userRepoImpl.VerifyUserByEmail() = %v, want %v", got, tt.want)
			}
		})
	}
}
