package service

import (
	"context"
	"reflect"
	"testing"

	"github.com/google/uuid"
)

func TestNewAuthClient(t *testing.T) {
	type args struct {
		secret string
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
			if got := NewAuthClient(tt.args.secret); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAuthClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_authClient_GenerateToken(t *testing.T) {
	type args struct {
		ctx           context.Context
		id            uuid.UUID
		expiredMinute int
	}
	tests := []struct {
		name      string
		c         *authClient
		args      args
		wantToken string
		wantErr   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotToken, err := tt.c.GenerateToken(tt.args.ctx, tt.args.id, tt.args.expiredMinute)
			if (err != nil) != tt.wantErr {
				t.Errorf("authClient.GenerateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotToken != tt.wantToken {
				t.Errorf("authClient.GenerateToken() = %v, want %v", gotToken, tt.wantToken)
			}
		})
	}
}

func Test_authClient_ValidateToken(t *testing.T) {
	type args struct {
		ctx          context.Context
		requestToken string
	}
	tests := []struct {
		name    string
		c       *authClient
		args    args
		want    uuid.UUID
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.c.ValidateToken(tt.args.ctx, tt.args.requestToken)
			if (err != nil) != tt.wantErr {
				t.Errorf("authClient.ValidateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("authClient.ValidateToken() = %v, want %v", got, tt.want)
			}
		})
	}
}
