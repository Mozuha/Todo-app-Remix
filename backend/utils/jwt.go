package utils

import (
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

/*
include session id in the token so that we can invalidate the token by deleting the session
https://zenn.dev/ritou/articles/4a5d6597a5f250#%E3%80%8C%E3%82%BB%E3%83%83%E3%82%B7%E3%83%A7%E3%83%B3id%E3%82%92jwt%E3%81%AB%E5%86%85%E5%8C%85%E3%80%8D%E3%81%A8%E3%81%84%E3%81%86%E8%80%83%E3%81%88%E6%96%B9
*/

type jwtCustomClaims struct {
	UserID    string `json:"user_id"`
	SessionID string `json:"session_id"`
	jwt.RegisteredClaims
}

const issuer = "example_issuer"

func GenerateJWT(userID string, sessionID string) (string, error) {
	tokenLifeSpanHour, err := strconv.Atoi(os.Getenv("JWT_ACCESS_TOKEN_EXP_HOUR"))
	if err != nil {
		return "", err
	}

	claims := &jwtCustomClaims{
		userID,
		sessionID,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(tokenLifeSpanHour))),
			Issuer:    issuer,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func ValidateJWT(tokenString string) (*jwtCustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if claims, ok := token.Claims.(*jwtCustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, err
}
