package profile

import (
	"learn_bot_admin_panel/tools/sql_null"

	"github.com/google/uuid"
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

func (u User) IsActivated() bool {
	return u.TgID.Valid
}

func (u User) IsPasswordSet() bool {
	return u.Password.Valid
}

type UUIDResponse struct {
	AuthUrl string `json:"auth_url"`
	UUID    string `json:"uu_id"`
}

func NewUUIDResponse() UUIDResponse {
	return UUIDResponse{
		UUID: uuid.NewString(),
	}
}
