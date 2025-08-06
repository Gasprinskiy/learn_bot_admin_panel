package bot_users

import (
	"time"
)

// BotUserProfile represents the bot_users_profile table in the database.
type BotUserProfile struct {
	UID              int64     `db:"u_id" json:"u_id"`
	TgID             int64     `db:"tg_id" json:"tg_id"`
	TgUserName       string    `db:"tg_user_name" json:"tg_user_name"`
	FirstName        string    `db:"first_name" json:"first_name"`
	LastName         string    `db:"last_name" json:"last_name"`
	BirthDate        time.Time `db:"birth_date" json:"birth_date"`
	PhoneNumber      string    `db:"phone_number" json:"phone_number"`
	JoinDate         time.Time `db:"join_date" json:"join_date"`
	RegisterDate     time.Time `db:"register_date" json:"register_date"`
	CommonTotalCount int       `db:"total_count" json:"-"`
}
