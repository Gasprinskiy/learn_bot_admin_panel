package bot_users

import (
	"learn_bot_admin_panel/tools/gennull"
	"learn_bot_admin_panel/tools/sql_null"
	"mime/multipart"
	"strings"
	"time"
)

type FindBotRegisteredUsersQuertParseParam struct {
	Limit              int       `form:"limit"`
	PageCount          int       `form:"page"`
	Query              string    `form:"query"`
	NextCursorID       int       `form:"next_cursor_id"`
	NextCursorDate     time.Time `form:"next_cursor_date"`
	JoinDateFrom       time.Time `form:"join_date_from"`
	JoinDateTill       time.Time `form:"join_date_till"`
	AgeFrom            int       `form:"age_from"`
	AgeTill            int       `form:"age_till"`
	SubscriptionStatus string    `form:"subscription_status"`
}

func (fp FindBotRegisteredUsersQuertParseParam) InnerParam() FindBotRegisteredUsersInnerParam {
	var (
		query          sql_null.NullString
		nextCursorDate sql_null.NullTime
		nextCursorID   sql_null.NullInt64
		birthDateFrom  sql_null.NullTime
		birthDateTill  sql_null.NullTime
		joinDateFrom   sql_null.NullTime
		joinDateTill   sql_null.NullTime
		suStatus       gennull.GenericNull[SubscriptionStatus]
	)

	trimQuery := strings.TrimSpace(fp.Query)

	if len([]rune(trimQuery)) > 0 {
		query = sql_null.NewString(trimQuery)
	}

	if !fp.NextCursorDate.IsZero() {
		nextCursorDate = sql_null.NewNullTime(fp.NextCursorDate)
	}

	if fp.NextCursorID > 0 {
		nextCursorID = sql_null.NewInt64(fp.NextCursorID)
	}

	if fp.AgeFrom > 0 {
		date := time.Now().AddDate(-fp.AgeFrom, 0, 0)
		birthDateFrom = sql_null.NewNullTime(time.Date(date.Year(), 1, 1, 0, 0, 0, 0, time.Local))
	}

	if fp.AgeTill > 0 {
		date := time.Now().AddDate(-fp.AgeTill, 0, 0)
		birthDateTill = sql_null.NewNullTime(time.Date(date.Year(), 1, 1, 0, 0, 0, 0, time.Local))
	}

	if !fp.JoinDateFrom.IsZero() {
		joinDateFrom = sql_null.NewNullTime(fp.JoinDateFrom)
	}

	if !fp.JoinDateTill.IsZero() {
		joinDateTill = sql_null.NewNullTime(fp.JoinDateTill)
	}

	if fp.SubscriptionStatus != "" {
		suStatus = gennull.NewGenericNull(SubscriptionStatus(fp.SubscriptionStatus))
	}

	return FindBotRegisteredUsersInnerParam{
		Limit:              fp.Limit,
		PageCount:          fp.PageCount,
		Query:              query,
		NextCursorDate:     nextCursorDate,
		NextCursorID:       nextCursorID,
		JoinDateFrom:       joinDateFrom,
		JoinDateTill:       joinDateTill,
		BirthDateFrom:      birthDateFrom,
		BirthDateTill:      birthDateTill,
		SubscriptionStatus: suStatus,
	}
}

type FindBotRegisteredUsersInnerParam struct {
	Limit              int `db:"limit"`
	PageCount          int
	Query              sql_null.NullString `db:"query"`
	NextCursorDate     sql_null.NullTime   `db:"next_cursor_date"`
	NextCursorID       sql_null.NullInt64  `db:"next_cursor_id"`
	BirthDateFrom      sql_null.NullTime   `db:"birth_date_from"`
	BirthDateTill      sql_null.NullTime   `db:"birth_date_till"`
	JoinDateFrom       sql_null.NullTime   `db:"join_date_from"`
	JoinDateTill       sql_null.NullTime   `db:"join_date_till"`
	SubscriptionStatus gennull.GenericNull[SubscriptionStatus]
}

type PurchaseSubscriptionParam struct {
	BotUserID int
	ManagerID int
	SubID     int
	File      multipart.File
}

type PurchaseSubscriptionDbParam struct {
	BotUserID    int       `db:"u_id"`
	ManagerID    int       `db:"manager_id"`
	SubID        int       `db:"sub_id"`
	PurchaseTime time.Time `db:"p_time"`
}

func (p PurchaseSubscriptionParam) NewPurchaseSubscriptionDbParam() PurchaseSubscriptionDbParam {
	return PurchaseSubscriptionDbParam{
		BotUserID:    p.BotUserID,
		ManagerID:    p.ManagerID,
		SubID:        p.SubID,
		PurchaseTime: time.Now(),
	}
}
