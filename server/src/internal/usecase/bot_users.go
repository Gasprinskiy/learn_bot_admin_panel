package usecase

import (
	"context"
	"learn_bot_admin_panel/config"
	"learn_bot_admin_panel/internal/entity/bot_users"
	"learn_bot_admin_panel/internal/entity/global"
	"learn_bot_admin_panel/internal/transaction"
	"learn_bot_admin_panel/rimport"
	"learn_bot_admin_panel/tools/logger"

	"github.com/go-telegram/bot"
)

type BotUsers struct {
	ri     *rimport.RepositoryImports
	log    *logger.Logger
	config *config.Config
	b      *bot.Bot
}

func NewBotUsers(
	ri *rimport.RepositoryImports,
	log *logger.Logger,
	config *config.Config,
	b *bot.Bot,
) *BotUsers {
	return &BotUsers{ri, log, config, b}
}

func (u *BotUsers) logPrefix() string {
	return "[bot_users]"
}

func (u *BotUsers) FindRegisteredUsers(
	ctx context.Context,
	param bot_users.FindBotRegisteredUsersParam,
) (global.CommotListSearchResponse[bot_users.BotUserProfile], error) {
	var zero global.CommotListSearchResponse[bot_users.BotUserProfile]

	ts := transaction.MustGetSession(ctx)

	data, err := u.ri.Repository.BotUsers.FindBotRegisteredUsers(ts, param)
	switch err {
	case nil:
	case global.ErrNoData:
		return zero, err

	default:
		u.log.Db.Errorln(u.logPrefix(), "не удалось найти зарегестрированных в боте пользователей:", err)
		return zero, global.ErrInternalError
	}

	result := global.NewCommotListSearchResponse(data, data[0].CommonTotalCount)

	return result, nil
}
