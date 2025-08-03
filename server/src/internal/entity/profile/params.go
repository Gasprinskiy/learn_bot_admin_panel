package profile

import (
	"learn_bot_admin_panel/tools/sql_null"
)

type CreateProfileParam struct {
	FirstName  string              `json:"first_name" db:"first_name"`
	LastName   string              `json:"last_name" db:"last_name"`
	TgUserName string              `json:"tg_user_name" db:"tg_user_name"`
	TgID       sql_null.NullString `json:"tg_id" db:"tg_id"`
	Access     AccessRight         `json:"access_right" db:"access_right"`
}

type PasswordLoginParam struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

type SetPasswordParam struct {
	Password string `json:"password"`
}
