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
			u.tg_id,
			u.password
		FROM admin_panel_users u
    	JOIN admin_panel_acces_rights ar ON (ar.ar_id = u.ar_id)
		WHERE u.tg_user_name = $1
	`

	return sql_gen.Get[profile.User](SqlxTx(ts), sqlQuery, userName)
}

func (r *profileRepo) FindProfileByID(ts transaction.Session, userID int) (profile.User, error) {
	sqlQuery := `
		SELECT
			u.u_id,
			ar.name as access_right,
			u.first_name,
			u.last_name,
			u.tg_user_name,
			u.tg_id,
			u.password
		FROM admin_panel_users u
    	JOIN admin_panel_acces_rights ar ON (ar.ar_id = u.ar_id)
		WHERE u.u_id = $1
	`

	return sql_gen.Get[profile.User](SqlxTx(ts), sqlQuery, userID)
}

func (r *profileRepo) FindUserDeviceIDList(ts transaction.Session, userID int) ([]string, error) {
	sqlQuery := `
		SELECT
			ud.device_id
		FROM admin_panel_user_devices ud
		WHERE ud.u_id = $1
	`

	return sql_gen.Select[string](SqlxTx(ts), sqlQuery, userID)
}

func (r *profileRepo) CreateUserDeviceID(ts transaction.Session, userID int, deviceID string) error {
	sqlQuery := `INSERT INTO admin_panel_user_devices (u_id, device_id) VALUES ($1, $2)`

	_, err := SqlxTx(ts).Exec(sqlQuery, userID, deviceID)
	return err
}

func (r *profileRepo) SetProfilePassword(ts transaction.Session, userID int, password string) error {
	sqlQuery := `
		UPDATE admin_panel_users
			SET password = $1
		WHERE u_id = $2
	`

	_, err := SqlxTx(ts).Exec(sqlQuery, password, userID)
	return err
}

func (r *profileRepo) SetProfileTGID(ts transaction.Session, userID int, TGID int64) error {
	sqlQuery := `
		UPDATE admin_panel_users
			SET tg_id = $1
		WHERE u_id = $2
	`

	_, err := SqlxTx(ts).Exec(sqlQuery, TGID, userID)
	return err
}
