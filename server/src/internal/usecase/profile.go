package usecase

import (
	"context"
	"learn_bot_admin_panel/internal/entity/global"
	"learn_bot_admin_panel/internal/entity/profile"
	"learn_bot_admin_panel/internal/transaction"
	"learn_bot_admin_panel/rimport"
	"learn_bot_admin_panel/tools/logger"
)

type Profile struct {
	ri  *rimport.RepositoryImports
	log *logger.Logger
}

func NewProfile(ri *rimport.RepositoryImports, log *logger.Logger) *Profile {
	return &Profile{ri, log}
}

func (u *Profile) logPrefix() string {
	return "[profile]"
}

func (u *Profile) CreateProfile(ctx context.Context, param profile.CreateProfileParam) (int, error) {
	ts := transaction.MustGetSession(ctx)

	ID, err := u.ri.Repository.Profile.CreateProfile(ts, param)
	if err != nil {
		u.log.Db.Errorln(u.logPrefix(), "не удалось создать новый профиль: ", err)
		return ID, global.ErrInternalError
	}

	return ID, nil
}
