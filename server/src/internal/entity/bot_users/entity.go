package bot_users

import (
	"learn_bot_admin_panel/tools/sql_null"
	"time"
)

type BotUserProfile struct {
	UID                int64              `db:"u_id" json:"u_id"`
	TgID               int64              `db:"tg_id" json:"tg_id"`
	TgUserName         string             `db:"tg_user_name" json:"tg_user_name" excel_head:"Юзернейм" excel_cell:"string"`
	FirstName          string             `db:"first_name" json:"first_name" excel_head:"Имя" excel_cell:"string"`
	LastName           string             `db:"last_name" json:"last_name" excel_head:"Фамилия" excel_cell:"string"`
	BirthDate          time.Time          `db:"birth_date" json:"birth_date" excel_head:"Дата рождения" excel_cell:"date"`
	PhoneNumber        string             `db:"phone_number" json:"phone_number" excel_head:"Номер телефона" excel_cell:"string"`
	JoinDate           time.Time          `db:"join_date" json:"join_date" excel_head:"Дата вступления" excel_cell:"date"`
	RegisterDate       time.Time          `db:"register_date" json:"register_date"`
	CommonTotalCount   int                `db:"total_count" json:"-"`
	SubscrID           sql_null.NullInt64 `db:"sub_id" json:"subscription_id"`
	SubscrPurchaseDate sql_null.NullTime  `db:"p_time" json:"-"`
	SubscrTerm         sql_null.NullInt64 `db:"term_in_month" json:"subscription_term"`
	//
	SubscrStatus SubscriptionStatus `json:"subscription_status" excel_head:"Подписка" excel_cell:"string"`
}

func (b *BotUserProfile) SetSubscriptionStatus(value SubscriptionStatus) {
	b.SubscrStatus = value
}
