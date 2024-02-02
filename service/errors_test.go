package service

import (
	"testing"

	"google.golang.org/grpc/codes"
)

func TestNewError(t *testing.T) {
	type args struct {
		code       codes.Code
		codeDetail string
		message    string
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
			if err := NewError(tt.args.code, tt.args.codeDetail, tt.args.message); (err != nil) != tt.wantErr {
				t.Errorf("NewError() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
