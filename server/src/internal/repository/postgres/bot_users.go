package postgres

import (
	"fmt"
	"learn_bot_admin_panel/internal/entity/bot_users"
	"learn_bot_admin_panel/internal/repository"
	"learn_bot_admin_panel/internal/transaction"
	"learn_bot_admin_panel/tools/sql_gen"
	"strings"
)

type botUsers struct{}

func NewBotUsers() repository.BotUsers {
	return &botUsers{}
}

const (
	whereMark      = "-- put_where"
	filtersMark    = "-- put_filters"
	paginationMark = "-- put_pagination"
)

var filtersQueryMap = map[string]string{
	bot_users.FilterKeyQuery(true): `(
		bu.tg_user_name ILIKE '%' || :query || '%'
		OR
		bu.first_name ILIKE '%' || :query || '%'
		OR
		bu.last_name ILIKE '%' || :query || '%'
  )`,
	bot_users.FilterKeyBirthDate(true, false): ` bu.birth_date <= :birth_date_from`,
	bot_users.FilterKeyBirthDate(false, true): ` bu.birth_date >= :birth_date_till`,
	bot_users.FilterKeyBirthDate(true, true):  ` (bu.birth_date <= :birth_date_from AND bu.birth_date >= :birth_date_till)`,
	bot_users.FilterKeyJoinDate(true, false):  ` bu.join_date >= :join_date_from`,
	bot_users.FilterKeyJoinDate(false, true):  ` bu.join_date <= :join_date_till`,
	bot_users.FilterKeyJoinDate(true, true):   ` (bu.join_date >= :join_date_from AND bu.join_date <= :join_date_till)`,
	bot_users.FilterKeyPurchases(true):        ` bp.sub_id IS NOT NULL`,
}

func (r *botUsers) FindBotRegisteredUsers(ts transaction.Session, param bot_users.FindBotRegisteredUsersInnerParam) ([]bot_users.BotUserProfile, error) {
	sqlQuery := `
 		WITH filtered AS (
			SELECT
				bu.u_id,
				bu.tg_id,
				bu.tg_user_name,
				bu.first_name,
				bu.last_name,
				bu.birth_date,
				bu.phone_number,
				bu.join_date,
				bu.register_date,
				bp.sub_id,
				bp.p_time,
				bst.term_in_month
			FROM bot_users_profile bu
				LEFT JOIN bot_users_purchases bp ON (bp.u_id = bu.u_id)
				LEFT JOIN bot_subscription_types bst ON (bst.sub_id = bp.sub_id)
			%s
			%s
		)
		SELECT
			data.u_id,
			data.tg_id,
			data.tg_user_name,
			data.first_name,
			data.last_name,
			data.birth_date,
			data.phone_number,
			data.join_date,
			data.register_date,
			data.sub_id,
			data.p_time,
			data.term_in_month,
			(SELECT COUNT(*) FROM filtered) AS total_count
		FROM (
			SELECT
				*
			FROM filtered bu
			%s
			ORDER BY join_date DESC, u_id DESC
			LIMIT :limit
		) AS data
	`

	sqlQuery = fmt.Sprintf(sqlQuery, whereMark, filtersMark, paginationMark)

	var filterQuery string

	filtersQueryKeys := [4]string{
		bot_users.FilterKeyQuery(param.Query.Valid),
		bot_users.FilterKeyBirthDate(param.BirthDateFrom.Valid, param.BirthDateTill.Valid),
		bot_users.FilterKeyJoinDate(param.BirthDateFrom.Valid, param.BirthDateTill.Valid),
		bot_users.FilterKeyPurchases(param.SubscriptionIsActive),
	}

	for i, key := range filtersQueryKeys {
		query, exists := filtersQueryMap[key]
		if !exists {
			continue
		}

		filterQuery += query
		if i > 0 && (i+1 < len(filtersQueryKeys)) {
			filterQuery += " AND"
		}
	}

	if filterQuery != "" {
		sqlQuery = strings.Replace(sqlQuery, whereMark, "WHERE ", 1)
		sqlQuery = strings.Replace(sqlQuery, filtersMark, filterQuery, 1)
	}

	if param.NextCursorDate.Valid && param.NextCursorID.Valid {
		paginationQuery := `WHERE	(bu.join_date, bu.u_id) < (:next_cursor_date, :next_cursor_id)`

		sqlQuery = strings.Replace(sqlQuery, paginationMark, paginationQuery, 1)
	}

	fmt.Println("sqlQuery: ", sqlQuery)

	return sql_gen.SelectNamedStruct[bot_users.BotUserProfile](SqlxTx(ts), sqlQuery, param)
}
