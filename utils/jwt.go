package utils

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

var (
	KEY = []byte("im")
)

type UserClaim struct {
	Identity string `json:"identity"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}

func GenerateJWT(identity, email string) (string, error) {
	userClaim := UserClaim{
		Identity: identity,
		Email:    email,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "im",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 12)), // 30 days
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaim)

	signingString, err := token.SignedString(KEY)
	if err != nil {
		return "", err
	}
	return signingString, nil
}

func ParseJWT(tokenString string) (*UserClaim, error) {
	userClaim := new(UserClaim)
	claims, err := jwt.ParseWithClaims(tokenString, userClaim, func(token *jwt.Token) (interface{}, error) {
		return KEY, nil
	})
	if err != nil {
		return nil, err
	}

	if !claims.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return userClaim, nil
}
