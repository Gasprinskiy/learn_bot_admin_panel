package postgres

import (
	"learn_bot_admin_panel/internal/entity/bot_users"
	"learn_bot_admin_panel/internal/repository"
	"learn_bot_admin_panel/internal/transaction"
	"learn_bot_admin_panel/tools/sql_gen"
)

type botUsers struct{}

func NewBotUsers() repository.BotUsers {
	return &botUsers{}
}

func (r *botUsers) FindBotRegisteredUsers(ts transaction.Session, param bot_users.FindBotRegisteredUsersParam) ([]bot_users.BotUserProfile, error) {
	sqlQuery := `
 		WITH filtered AS (
			SELECT *
			FROM bot_users_profile bu
			WHERE
				(
					:query IS NULL
					OR bu.tg_user_name ILIKE '%' || :query || %
					OR bu.first_name ILIKE '%' || :query || '%'
					OR bu.last_name ILIKE '%' || :query || '%'
				)
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
				WHERE
					(
						(:join_date IS NULL AND :cursor_u_id IS NULL)
						OR
						(bu.join_date, bu.u_id) < (:next_cursor_date, :next_cursor_id)
					)
				ORDER BY join_date DESC, u_id DESC
				LIMIT :limit
		) AS data
	`

	return sql_gen.SelectNamedStruct[bot_users.BotUserProfile](SqlxTx(ts), sqlQuery, param)
}
