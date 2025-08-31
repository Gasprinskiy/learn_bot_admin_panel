package profile

import (
	"learn_bot_admin_panel/tools/sql_null"
)

type CreateProfileParam struct {
	FirstName  string              `form:"first_name" db:"first_name"`
	LastName   string              `form:"last_name" db:"last_name"`
	TgUserName string              `form:"tg_user_name" db:"tg_user_name"`
	TgID       sql_null.NullString `form:"tg_id" db:"tg_id"`
	Access     AccessRight         `form:"access_right" db:"access_right"`
	//
	ArID int `form:"-" db:"ar_id"`
}

func (p *CreateProfileParam) SetAccesRightID() {
	arID, exists := AccessRightIDMap[p.Access]
	if !exists {
		arID = AccessRightIDMap[AccessRightManager]
	}

	p.ArID = arID
}

type RedactProfileParam struct {
	ID         int    `form:"id" db:"u_id"`
	FirstName  string `form:"first_name" db:"first_name"`
	LastName   string `form:"last_name" db:"last_name"`
	TgUserName string `form:"tg_user_name" db:"tg_user_name"`
}

type PasswordLoginParam struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

type SetPasswordParam struct {
	Password string `json:"password"`
}
