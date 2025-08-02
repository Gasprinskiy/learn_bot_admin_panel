package middleware

import (
	"learn_bot_admin_panel/internal/entity/global"
	"learn_bot_admin_panel/internal/usecase"
	"learn_bot_admin_panel/tools/gin_gen"

	"github.com/gin-gonic/gin"
)

const (
	UserDataKey = "user_data"
)

type AuthMiddleware struct {
	jwtUsecase *usecase.Jwt
}

func NewAuthMiddleware(jwtUsecase *usecase.Jwt) *AuthMiddleware {
	return &AuthMiddleware{jwtUsecase}
}

func (m *AuthMiddleware) CheckAccesToken() gin.HandlerFunc {
	return func(gctx *gin.Context) {
		token, err := gctx.Cookie("access_token")
		if err != nil {
			gin_gen.HandleError(gctx, global.ErrPermissionDenied)
			gctx.Abort()
			return
		}

		claims, err := m.jwtUsecase.ParseToken(token)
		if err != nil {
			gin_gen.HandleError(gctx, global.ErrExpired)
			gctx.Abort()
			return
		}

		gctx.Set(UserDataKey, claims)
		gctx.Next()
	}
}
