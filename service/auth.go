package service

// import (
// 	"context"
// 	"fmt"
// 	"net/http"
// 	"strconv"
// 	"time"

// 	"github.com/Mitra-Apps/be-user-service/config/tools"
// 	"github.com/Mitra-Apps/be-user-service/domain/user/repository"
// 	"github.com/Mitra-Apps/be-user-service/lib"
// 	jwt "github.com/golang-jwt/jwt/v5"
// 	"github.com/google/uuid"
// 	"github.com/labstack/echo"
// 	"google.golang.org/grpc/codes"
// )

// type authClient struct {
// 	secret   string
// 	userRepo repository.User
// }

// type Authentication interface {
// 	GenerateToken(ctx context.Context, userId uuid.UUID) (string, error)
// 	VerifyToken(tokenString string) (uuid.UUID, error)
// }

// func NewAuthClient(secret string, userRepo repository.User) Authentication {
// 	return &authClient{
// 		secret:   secret,
// 		userRepo: userRepo,
// 	}
// }

// // GenerateJWT generates a JWT token with a specific payload
// func (a *authClient) GenerateToken(ctx context.Context, userId uuid.UUID) (string, error) {

// 	if userId == uuid.Nil {
// 		return "", NewError(codes.InvalidArgument, &tools.ErrorResponse{
// 			Code:       codes.InvalidArgument.String(),
// 			CodeDetail: codes.Unknown.String(),
// 			Message:    "need user id",
// 		})
// 	}

// 	currentTime := time.Now().UTC()
// 	//set token with criteria below and input userID into subject
// 	//this will be needed to check which user is this token for
// 	claims := &jwt.RegisteredClaims{
// 		Subject:   strconv.Itoa(int(userID)),
// 		ExpiresAt: jwt.NewNumericDate(currentTime.Add(time.Hour * 12)),
// 		IssuedAt:  jwt.NewNumericDate(currentTime),
// 	}
// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 	t, err := token.SignedString([]byte(c.secret))
// 	if err != nil {
// 		return "", err
// 	}
// 	return t, err
// 	expireTime, err := time.ParseDuration(lib.GetEnv("JWT_EXPIRED_TIME"))
// 	var secretKey = []byte(lib.GetEnv("JWT_SECRET"))
// 	if err != nil {
// 		return "", echo.NewHTTPError(http.StatusBadRequest, "Invalid JWT expired time")
// 	}

// 	var roleNames []string
// 	for _, role := range user.Roles {
// 		roleNames = append(roleNames, role.RoleName)
// 	}

// 	// Define the token payload
// 	claims := &JwtCustomClaim{
// 		UserId:    user.Id.String(),
// 		RoleNames: roleNames,
// 		RegisteredClaims: jwt.RegisteredClaims{
// 			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expireTime)),
// 		},
// 	}

// 	// Create the token
// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

// 	// Sign the token with the secret key
// 	tokenString, err := token.SignedString(secretKey)
// 	if err != nil {
// 		return "", fmt.Errorf("error signing token: %v", err)
// 	}

// 	return tokenString, nil
// }

// func (a *authClient) VerifyToken(tokenString string) (uuid.UUID, error) {
// 	var secretKey = []byte(lib.GetEnv("JWT_SECRET"))
// 	// Parse the token
// 	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 		// Check signing method
// 		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
// 		}

// 		// Provide the key for validation
// 		return secretKey, nil
// 	})

// 	if err != nil {
// 		return nil, fmt.Errorf("error parsing token: %v", err)
// 	}

// 	// Validate the token
// 	if !token.Valid {
// 		return nil, fmt.Errorf("token is not valid")
// 	}

// 	return token, nil
// }
