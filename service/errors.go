package service

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Mitra-Apps/be-user-service/config/tools"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrWrongPassword = errors.New("password is wrong")
)

func NewError(code codes.Code, newErr *tools.ErrorResponse) error {
	// Marshal the ErrorResponse struct to JSON
	errJSON, marshalErr := json.Marshal(newErr)
	if marshalErr != nil {
		fmt.Println("Error marshaling ErrorResponse:", marshalErr)
		return nil
	}

	// Convert the JSON byte slice to a string
	errString := string(errJSON)

	return status.Errorf(code, errString)
}
