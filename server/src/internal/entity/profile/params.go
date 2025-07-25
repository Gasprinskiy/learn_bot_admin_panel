package profile

import "lear_bot_admin_panel/src/tools/sql_null"

type AddProfileParam struct {
	FirstName   string              `json:"first_name" db:"first_name"`
	LastName    string              `json:"last_name" db:"last_name"`
	TgID        sql_null.NullString `json:"tg_id" db:"tg_id"`
	TgUserName  sql_null.NullString `json:"tg_user_name" db:"tg_user_name"`
	AccessRight AccessRight         `json:"access_right" db:"access_right"`
}
