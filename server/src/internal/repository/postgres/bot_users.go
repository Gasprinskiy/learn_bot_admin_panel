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
	limitMark      = "-- limit"
)

var filtersQueryMap = map[string]string{
	bot_users.FilterKeyQuery(true): ` (
		bu.tg_user_name ILIKE '%' || :query || '%'
		OR
		bu.phone_number ILIKE '%' || :query || '%'
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
	bot_users.FilterKeySubscriptionStatus(true, bot_users.SubscriptionStatusActive): ` (
		DATE_TRUNC('day', CURRENT_TIMESTAMP) < DATE_TRUNC('day', bp.p_time + make_interval(months := bst.term_in_month))
	)`,
	bot_users.FilterKeySubscriptionStatus(true, bot_users.SubscriptionStatusNotExists): ` bp.sub_id IS NULL`,
	bot_users.FilterKeySubscriptionStatus(true, bot_users.SubscriptionStatusExpired): ` (
		DATE_TRUNC('day', CURRENT_TIMESTAMP) >= DATE_TRUNC('day', bp.p_time + make_interval(months := bst.term_in_month))
	)`,
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
				LEFT JOIN LATERAL (
					SELECT bp.sub_id, bp.p_time
					FROM bot_users_purchases bp
					WHERE bp.u_id = bu.u_id
					ORDER BY bp.p_id DESC
					LIMIT 1
				) bp ON TRUE
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
				f.u_id,
				f.tg_id,
				f.tg_user_name,
				f.first_name,
				f.last_name,
				f.birth_date,
				f.phone_number,
				f.join_date,
				f.register_date,
				f.sub_id,
				f.p_time,
				f.term_in_month
			FROM filtered f
			%s
			ORDER BY join_date DESC, u_id DESC
			%s
		) AS data
	`

	sqlQuery = fmt.Sprintf(sqlQuery, whereMark, filtersMark, paginationMark, limitMark)

	var filterQuery string

	filtersQueryKeys := [4]string{
		bot_users.FilterKeyQuery(param.Query.Valid),
		bot_users.FilterKeyBirthDate(param.BirthDateFrom.Valid, param.BirthDateTill.Valid),
		bot_users.FilterKeyJoinDate(param.JoinDateFrom.Valid, param.JoinDateTill.Valid),
		bot_users.FilterKeySubscriptionStatus(param.SubscriptionStatus.Valid, param.SubscriptionStatus.Value),
	}

	existsFilters := make([]string, 0, len(filtersQueryKeys))

	for _, key := range filtersQueryKeys {
		query, exists := filtersQueryMap[key]
		if !exists {
			continue
		}

		existsFilters = append(existsFilters, query)
	}

	for i, filter := range existsFilters {
		filterQuery += filter

		if i == len(existsFilters)-1 {
			continue
		}

		filterQuery += " AND"
	}

	if filterQuery != "" {
		sqlQuery = strings.Replace(sqlQuery, whereMark, "WHERE ", 1)
		sqlQuery = strings.Replace(sqlQuery, filtersMark, filterQuery, 1)
	}

	if param.NextCursorDate.Valid && param.NextCursorID.Valid {
		paginationQuery := `WHERE	(f.join_date, f.u_id) < (:next_cursor_date, :next_cursor_id)`

		sqlQuery = strings.Replace(sqlQuery, paginationMark, paginationQuery, 1)
	}

	if param.Limit > 0 {
		sqlQuery = strings.Replace(sqlQuery, limitMark, "LIMIT :limit", 1)
	}

	return sql_gen.SelectNamedStruct[bot_users.BotUserProfile](SqlxTx(ts), sqlQuery, param)
}

func (r *botUsers) FindUserByID(ts transaction.Session, id int) (bot_users.BotUserCommonData, error) {
	sqlQuery := `
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
		FROM bot_users_profile bu
		WHERE bu.u_id = $1
	`

	return sql_gen.Get[bot_users.BotUserCommonData](SqlxTx(ts), sqlQuery, id)
}

func (r *botUsers) FindUserPurchases(ts transaction.Session, userID int) ([]bot_users.BotUserPurchase, error) {
	sqlQuery := `
		SELECT
			up.p_id,
			up.sub_id,
			up.discount,
			up.p_time,
			up.manager_id,
			up.payment_type_id,
			up.receipt_file_name,
			st.term_in_month,
			st.price,
			ps.first_name as manager_first_name,
			ps.last_name as manager_last_name
		FROM bot_users_purchases up
			JOIN bot_subscription_types st ON (st.sub_id = up.sub_id)
			LEFT JOIN admin_panel_users ps ON (ps.u_id = up.manager_id)
		WHERE up.u_id = $1
		ORDER BY up.p_id DESC
	`

	return sql_gen.Select[bot_users.BotUserPurchase](SqlxTx(ts), sqlQuery, userID)
}

func (r *botUsers) LoadAllBotSubscriptionTypes(ts transaction.Session) ([]bot_users.BotSubscriptionType, error) {
	sqlQuery := `
		SELECT
			bst.sub_id,
			bst.term_in_month,
			bst.price
		FROM bot_subscription_types bst
		ORDER BY bst.price
	`

	return sql_gen.Select[bot_users.BotSubscriptionType](SqlxTx(ts), sqlQuery)
}

func (r *botUsers) CreateSubscriptionPurchase(ts transaction.Session, param bot_users.Purchase) (int64, error) {
	sqlQuery := `
	INSERT INTO bot_users_purchases
		(sub_id, u_id, p_time, discount, receipt_json, manager_id, payment_type_id)
	VALUES
		(:sub_id, :u_id, :p_time, :discount, :receipt_json, :manager_id, :payment_type_id)
	RETURNING p_id
	`

	return sql_gen.ExecNamedReturnLastInsterted(SqlxTx(ts), sqlQuery, param)
}

func (r *botUsers) SavePurchaseFileName(ts transaction.Session, pID int, fileName string) error {
	sqlQuery := `UPDATE bot_users_purchases SET receipt_file_name = $2 WHERE p_id = $1`

	_, err := SqlxTx(ts).Exec(sqlQuery, pID, fileName)
	return err
}
