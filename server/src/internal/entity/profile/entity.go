package profile

import (
	"learn_bot_admin_panel/tools/sql_null"
)

type User struct {
	ID         int                 `json:"id" db:"u_id"`
	FirstName  string              `json:"first_name" db:"first_name"`
	LastName   string              `json:"last_name" db:"last_name"`
	TgUserName string              `json:"tg_user_name" db:"tg_user_name"`
	TgID       sql_null.NullString `json:"tg_id" db:"tg_id"`
	Password   sql_null.NullString `json:"-" db:"password"`
	Access     AccessRight         `json:"access_right" db:"access_right"`
}

type UserFirstLoginAnswer struct {
	IsPasswordSet bool `json:"is_password_set"`
}

type UserCommonInfo struct {
	ID            int         `json:"id" db:"u_id"`
	FirstName     string      `json:"first_name" db:"first_name"`
	LastName      string      `json:"last_name" db:"last_name"`
	TgUserName    string      `json:"tg_user_name" db:"tg_user_name"`
	IsPasswordSet bool        `json:"is_password_set"`
	Access        AccessRight `json:"access_right" db:"access_right"`
}

func (u User) NewUserCommonInfo() UserCommonInfo {
	return UserCommonInfo{
		ID:            u.ID,
		FirstName:     u.FirstName,
		LastName:      u.LastName,
		TgUserName:    u.TgUserName,
		IsPasswordSet: u.IsPasswordSet(),
		Access:        u.Access,
	}
}

func (u User) NewUserFirstLoginAnswer() UserFirstLoginAnswer {
	return UserFirstLoginAnswer{
		IsPasswordSet: u.IsPasswordSet(),
	}
}

func (u User) IsActivated() bool {
	return u.TgID.Valid
}

func (u User) IsPasswordSet() bool {
	return u.Password.Valid
}

type AuthUrlResponse struct {
	AuthUrl string `json:"auth_url"`
	UUID    string `json:"uu_id"`
}

func NewAuthUrlResponse(uuID, authUrl string) AuthUrlResponse {
	return AuthUrlResponse{
		AuthUrl: authUrl,
		UUID:    uuID,
	}
}
