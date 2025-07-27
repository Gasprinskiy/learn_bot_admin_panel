package postgres

import (
	"learn_bot_admin_panel/internal/entity/profile"
	"learn_bot_admin_panel/internal/repository"
	"learn_bot_admin_panel/internal/transaction"
	"learn_bot_admin_panel/tools/sql_gen"
)

type profileRepo struct{}

func NewProfile() repository.Profile {
	return &profileRepo{}
}

func (r *profileRepo) CreateProfile(ts transaction.Session, param profile.CreateProfileParam) (int, error) {
	return 1, nil
}

func (r *profileRepo) FindProfileByTGUserName(ts transaction.Session, userName string) (profile.User, error) {
	sqlQuery := `
		SELECT
			u.u_id,
			ar.name as access_right,
			u.first_name,
			u.last_name,
			u.tg_user_name,
			u.tg_id
		FROM admin_panel_users u
    	JOIN admin_panel_acces_rights ar ON (ar.ar_id = u.ar_id)
		WHERE u.tg_user_name = $1
	`

	return sql_gen.Get[profile.User](SqlxTx(ts), sqlQuery, userName)
}
