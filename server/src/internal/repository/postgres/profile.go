package postgres

import (
	"learn_bot_admin_panel/internal/entity/profile"
	"learn_bot_admin_panel/internal/repository"
	"learn_bot_admin_panel/internal/transaction"
	"learn_bot_admin_panel/tools/sql_gen"
	"time"
)

type profileRepo struct{}

func NewProfile() repository.Profile {
	return &profileRepo{}
}

func (r *profileRepo) CreateProfile(ts transaction.Session, param profile.CreateProfileParam) (int64, error) {
	sqlQuery := `
		INSERT INTO admin_panel_users
			(first_name, last_name, tg_user_name, ar_id)
		VALUES
			(:first_name, :last_name, :tg_user_name, :ar_id)
		RETURNING u_id
	`

	return sql_gen.ExecNamedReturnLastInsterted(SqlxTx(ts), sqlQuery, param)
}

func (r *profileRepo) RedactProfile(ts transaction.Session, param profile.RedactProfileParam) error {
	sqlQuery := `
		UPDATE admin_panel_users
		SET
			first_name = :first_name,
			last_name = :last_name,
			tg_user_name = :tg_user_name
		WHERE u_id = :u_id
	`

	return sql_gen.ExecNamed(SqlxTx(ts), sqlQuery, param)
}

func (r *profileRepo) DeleteProfile(ts transaction.Session, userID int) error {
	sqlQuery := `
		UPDATE admin_panel_users
		SET
			deleted = true
		WHERE u_id = $1
	`

	_, err := SqlxTx(ts).Exec(sqlQuery, userID)
	return err
}

func (r *profileRepo) FindProfileByTGUserNameOrID(ts transaction.Session, userName string, TGID int64) (profile.User, error) {
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
		OR u.tg_id = $2
	`

	return sql_gen.Get[profile.User](SqlxTx(ts), sqlQuery, userName, TGID)
}

func (r *profileRepo) LoadUsersProfile(ts transaction.Session) ([]profile.User, error) {
	sqlQuery := `
		SELECT
			u.u_id,
			ar.name as access_right,
			u.first_name,
			u.last_name,
			u.tg_user_name,
			u.last_login
		FROM admin_panel_users u
    	JOIN admin_panel_acces_rights ar ON (ar.ar_id = u.ar_id)
		WHERE u.deleted = false
		ORDER BY u.u_id DESC
	`

	return sql_gen.Select[profile.User](SqlxTx(ts), sqlQuery)
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
		AND u.deleted = false
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

func (r *profileRepo) SetProfileLastLoginDate(ts transaction.Session, userID int, date time.Time) error {
	sqlQuery := `
		UPDATE admin_panel_users
			SET last_login = $2
		WHERE u_id = $1
	`

	_, err := SqlxTx(ts).Exec(sqlQuery, userID, date)
	return err
}
