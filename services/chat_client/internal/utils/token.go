package utils

import (
	"context"
	"fmt"

	"github.com/golang-jwt/jwt"
	"google.golang.org/grpc/metadata"
)

const (
	authorisationHeader = "authorization"
	authPrefix          = "Bearer "
)

func PutAccessTokenToCtx(ctx context.Context, accessToken string) context.Context {
	md := metadata.New(map[string]string{authorisationHeader: authPrefix + accessToken})
	return metadata.NewOutgoingContext(ctx, md)
}

func GetUserIdFromToken(tokenString string) (int64, error) {
	var userID int64
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if iUserID, found := claims["id"]; found {
			userID = int64(iUserID.(float64))
		} else {
			return 0, fmt.Errorf("user id not found in token payload")
		}
	}

	return userID, nil
}
