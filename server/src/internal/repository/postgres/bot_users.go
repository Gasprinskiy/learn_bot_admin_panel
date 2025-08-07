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

func (r *botUsers) FindBotRegisteredUsers(ts transaction.Session, param bot_users.FindBotRegisteredUsersInnerParam) ([]bot_users.BotUserProfile, error) {
	const (
		WhereMark      = "-- put_where"
		FiltersMark    = "-- put_filters"
		PaginationMark = "-- put_pagination"
	)

	sqlQuery := `
 		WITH filtered AS (
			SELECT *
			FROM bot_users_profile bu
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
			(SELECT COUNT(*) FROM filtered) AS total_count
		FROM (
			SELECT
				bu.u_id,
				bu.tg_id,
				bu.tg_user_name,
				bu.first_name,
				bu.last_name,
				bu.birth_date,
				bu.phone_number,
				bu.join_date,
				bu.register_date
			FROM filtered bu
			%s
			ORDER BY join_date DESC, u_id DESC
			LIMIT :limit
		) AS data
	`

	var filtersQueryMap = map[string]string{
		"query:true": `(
			bu.tg_user_name ILIKE '%' || :query || '%'
			OR
			bu.first_name ILIKE '%' || :query || '%'
			OR
			bu.last_name ILIKE '%' || :query || '%'
  	)`,

		"birth_date:true:false": ` AND bu.birth_date <= :birth_date_from`,
		"birth_date:false:true": ` AND bu.birth_date >= :birth_date_till`,
	}

	sqlQuery = fmt.Sprintf(sqlQuery, WhereMark, FiltersMark, PaginationMark)

	var filterQuery string

	if param.Query.Valid {
		filterQuery += `(
			bu.tg_user_name ILIKE '%' || :query || '%'
			OR
			bu.first_name ILIKE '%' || :query || '%'
			OR
			bu.last_name ILIKE '%' || :query || '%'
  	)`
	}

	if param.BirthDateFrom.Valid && !param.BirthDateTill.Valid {
		filterQuery += ` AND bu.birth_date <= :birth_date_from`
	} else {
		filterQuery += ` AND bu.birth_date >= :birth_date_till`
	}

	if param.BirthDateFrom.Valid && param.BirthDateTill.Valid {
		filterQuery += ` AND (bu.birth_date <= :birth_date_from AND bu.birth_date >= :birth_date_till)`
	}

	if filterQuery != "" {
		sqlQuery = strings.Replace(sqlQuery, WhereMark, "WHERE ", 1)
		sqlQuery = strings.Replace(sqlQuery, FiltersMark, filterQuery, 1)
	}

	if param.NextCursorDate.Valid && param.NextCursorID.Valid {
		paginationQuery := `WHERE	(bu.join_date, bu.u_id) < (:next_cursor_date, :next_cursor_id)`

		sqlQuery = strings.Replace(sqlQuery, PaginationMark, paginationQuery, 1)
	}

	return sql_gen.SelectNamedStruct[bot_users.BotUserProfile](SqlxTx(ts), sqlQuery, param)
}
