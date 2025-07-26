package usecase

import (
	"learn_bot_admin_panel/config"
	"learn_bot_admin_panel/internal/entity/app_jwt"
	"learn_bot_admin_panel/internal/entity/global"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtUsecase struct {
	config *config.Config
}

func NewJwtUsecase(config *config.Config) *JwtUsecase {
	return &JwtUsecase{config}
}

func (u *JwtUsecase) GenerateToken(userID int) (string, error) {
	expirationTime := time.Now().Add(u.config.JwtSecretTTL)

	claims := app_jwt.Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secretKey := []byte(u.config.JwtSecret)

	result, err := token.SignedString(secretKey)
	if err != nil {
		err = global.ErrInternalError
	}

	return result, err
}

func (u *JwtUsecase) ParseToken(tokenString string) (*app_jwt.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &app_jwt.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(u.config.JwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*app_jwt.Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}
