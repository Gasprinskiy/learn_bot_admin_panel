package app_jwt

import (
	"learn_bot_admin_panel/internal/entity/profile"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID int                 `json:"user_id"`
	Access profile.AccessRight `json:"access_right"`
	jwt.RegisteredClaims
}

type TokenWithUserData struct {
	Token    string
	UserData profile.User
}

func NewTokenWithUserData(token string, userData profile.User) TokenWithUserData {
	return TokenWithUserData{
		Token:    token,
		UserData: userData,
	}
}
