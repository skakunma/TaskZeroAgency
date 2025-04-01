package jwtAuth

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/skakunma/TaskZeroAgency/internal/config"
	"time"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID int
}

var TokenEXP = time.Hour * 3

func BuildJWTString(cfg *config.Config, userID int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenEXP)),
		},
		UserID: userID,
	})

	tokenString, err := token.SignedString([]byte(cfg.SecretKEY))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
