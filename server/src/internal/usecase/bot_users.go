package usecase

import (
	"bufio"
	"bytes"
	"context"
	"learn_bot_admin_panel/config"
	"learn_bot_admin_panel/internal/entity/bot_users"
	"learn_bot_admin_panel/internal/entity/global"
	"learn_bot_admin_panel/internal/transaction"
	"learn_bot_admin_panel/rimport"
	"learn_bot_admin_panel/tools/chronos"
	"learn_bot_admin_panel/tools/excel"
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

func (u *BotUsers) PrintFindRegisteredUsers(
	ctx context.Context,
	param bot_users.FindBotRegisteredUsersInnerParam,
) ([]byte, error) {
	data, err := u.FindRegisteredUsers(ctx, param)
	if err != nil {
		return nil, err
	}

	file, err := excel.BuildExcelFileFromStruct(data.Data, "Пользователи бота")
	if err != nil {
		u.log.Db.Errorln(u.logPrefix(), "не удалось создать файл пользователей бота:", err)
		return nil, global.ErrInternalError
	}

	var b bytes.Buffer
	writer := bufio.NewWriter(&b)
	if err := file.Write(writer); err != nil {
		u.log.Db.Errorln("не удалось записать файл в буффер, ошибка: ", err)
		return nil, global.ErrInternalError
	}
	return b.Bytes(), nil
}

func (u *BotUsers) FindRegisteredUsers(
	ctx context.Context,
	param bot_users.FindBotRegisteredUsersInnerParam,
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

	result := global.NewCommotListSearchResponse(data, data[0].CommonTotalCount, param.Limit, param.PageCount)

	now := chronos.BeginingOfNow()

	for i, user := range result.Data {
		if user.SubscrPurchaseDate.Valid {
			subDate := chronos.BeginingOfDate(user.SubscrPurchaseDate.Time)
			expireDate := subDate.AddDate(0, user.SubscrTerm.GetInt(), 0)

			if now.After(expireDate) {
				user.SetSubscriptionStatus(bot_users.SubscriptionStatusExpired)
			} else {
				user.SetSubscriptionStatus(bot_users.SubscriptionStatusActive)
			}
		} else {
			user.SetSubscriptionStatus(bot_users.SubscriptionStatusNotExists)
		}

		result.Data[i] = user
	}

	return result, nil
}
