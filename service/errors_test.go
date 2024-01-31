package service

import (
	"testing"

	"github.com/Mitra-Apps/be-user-service/config"
	"google.golang.org/grpc/codes"
)

func TestNewError(t *testing.T) {
	type args struct {
		code   codes.Code
		newErr *config.ErrorResponse
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := NewError(tt.args.code, tt.args.newErr); (err != nil) != tt.wantErr {
				t.Errorf("NewError() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
