package usecase

import (
	"context"
	"learn_bot_admin_panel/config"
	"learn_bot_admin_panel/internal/entity/app_jwt"
	"learn_bot_admin_panel/internal/entity/global"
	"learn_bot_admin_panel/internal/entity/profile"
	"learn_bot_admin_panel/rimport"
	"learn_bot_admin_panel/tools/logger"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Jwt struct {
	ri     *rimport.RepositoryImports
	log    *logger.Logger
	config *config.Config
}

func NewJwt(ri *rimport.RepositoryImports, log *logger.Logger, config *config.Config) *Jwt {
	return &Jwt{ri, log, config}
}

func (u *Jwt) GenerateToken(userID int, accRight profile.AccessRight) (string, error) {
	expirationTime := time.Now().Add(u.config.JwtSecretTTL)

	claims := app_jwt.Claims{
		UserID: userID,
		Access: accRight,
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

func (u *Jwt) ParseToken(tokenString string) (*app_jwt.Claims, error) {
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

func (u *Jwt) GenerateTokenByTempAuthData(ctx context.Context, authKey string) (app_jwt.TokenWithUserData, error) {
	var zero app_jwt.TokenWithUserData

	userData, err := u.ri.Repository.AuthCache.GetTempUserData(ctx, authKey)
	if err != nil {
		u.log.Db.Errorln("не удалось найти временные данные пользователя в кеше:", err)
		return zero, global.ErrInternalError
	}

	token, err := u.GenerateToken(userData.ID, userData.Access)
	if err != nil {
		u.log.Db.Errorln("не удалось сгенерировать токен:", err)
		return zero, global.ErrInternalError
	}

	return app_jwt.NewTokenWithUserData(token, userData), nil
}
