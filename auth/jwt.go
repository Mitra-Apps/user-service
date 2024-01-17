package auth

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Mitra-Apps/be-user-service/domain/user/entity"
	"github.com/Mitra-Apps/be-user-service/lib"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo"
)

type JwtCustomClaim struct {
	UserId    string `json:"userId"`
	RoleNames []string
	jwt.RegisteredClaims
}

var secretKey = []byte(lib.GetEnv("JWT_SECRET"))

// GenerateJWT generates a JWT token with a specific payload
func GenerateToken(ctx context.Context, user *entity.User) (string, error) {
	expireTime, err := time.ParseDuration(lib.GetEnv("JWT_EXPIRED_TIME"))
	if err != nil {
		return "", echo.NewHTTPError(http.StatusBadRequest, "Invalid JWT expired time")
	}

	var roleNames []string
	for _, role := range user.Roles {
		roleNames = append(roleNames, role.RoleName)
	}

	// Define the token payload
	claims := &JwtCustomClaim{
		UserId:    user.Id.String(),
		RoleNames: roleNames,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expireTime)),
		},
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", fmt.Errorf("error signing token: %v", err)
	}

	return tokenString, nil
}

func VerifyToken(tokenString string) (*jwt.Token, error) {
	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Check signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Provide the key for validation
		return secretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("error parsing token: %v", err)
	}

	// Validate the token
	if !token.Valid {
		return nil, fmt.Errorf("token is not valid")
	}

	return token, nil
}
