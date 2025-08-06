package bot_users

import "learn_bot_admin_panel/tools/sql_null"

type FindBotRegisteredUsersParam struct {
	Limit          int                 `json:"limit" db:"limit"`
	PageCount      int                 `json:"page" db:"-"`
	Query          sql_null.NullString `json:"query" db:"query"`
	NextCursorDate sql_null.NullTime   `json:"next_cursor_date" db:"next_cursor_date"`
	NextCursorID   sql_null.NullInt64  `json:"next_cursor_id" db:"next_cursor_id"`
}
