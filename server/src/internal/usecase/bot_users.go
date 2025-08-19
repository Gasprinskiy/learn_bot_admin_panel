package usecase

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"learn_bot_admin_panel/config"
	"learn_bot_admin_panel/internal/entity/bot_users"
	"learn_bot_admin_panel/internal/entity/global"
	"learn_bot_admin_panel/internal/transaction"
	"learn_bot_admin_panel/rimport"
	"learn_bot_admin_panel/tools/chronos"
	"learn_bot_admin_panel/tools/excel"
	"learn_bot_admin_panel/tools/logger"
	"learn_bot_admin_panel/tools/sql_null"
	"os"
	"path/filepath"

	"github.com/go-telegram/bot"
	"github.com/sirupsen/logrus"
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

	for i, user := range result.Data {
		user = u.setUserSubscriptionStatus(user)
		result.Data[i] = user
	}

	return result, nil
}

func (u *BotUsers) FindUserByID(ctx context.Context, id int) (bot_users.BotUserProfile, error) {
	var zero bot_users.BotUserProfile

	ts := transaction.MustGetSession(ctx)

	user, err := u.ri.Repository.BotUsers.FindUserByID(ts, id)
	switch err {
	case nil:
		user = u.setUserSubscriptionStatus(user)
	case global.ErrNoData:
		return zero, err

	default:
		u.log.Db.Errorln(u.logPrefix(), "не удалось найти пользователя:", err)
		return zero, global.ErrInternalError
	}

	return user, nil
}

func (u *BotUsers) setUserSubscriptionStatus(user bot_users.BotUserProfile) bot_users.BotUserProfile {
	now := chronos.BeginingOfNow()

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

	return user
}

func (u *BotUsers) LoadAllBotSubscriptionTypes(ctx context.Context) ([]bot_users.BotSubscriptionType, error) {
	ts := transaction.MustGetSession(ctx)

	user, err := u.ri.Repository.BotUsers.LoadAllBotSubscriptionTypes(ts)
	if err != nil {
		u.log.Db.Errorln(u.logPrefix(), "не удалось загрузить список подписок:", err)
		return nil, global.ErrInternalError
	}

	return user, nil
}

func (u *BotUsers) PurchaseSubscription(ctx context.Context, param bot_users.PurchaseSubscriptionParam) error {
	lf := logrus.Fields{
		"u_id":     param.ManagerID,
		"bot_u_id": param.BotUserID,
	}

	ts := transaction.MustGetSession(ctx)

	purchase := bot_users.NewPurchase(
		param.SubID,
		param.BotUserID,
		chronos.BeginingOfNow(),
		sql_null.NullFloat64{},
		sql_null.NullString{},
		sql_null.NewInt64(param.ManagerID),
	)

	purchaseID, err := u.ri.BotUsers.CreateSubscriptionPurchase(ts, purchase)
	if err != nil {
		u.log.Db.WithFields(lf).Errorln(u.logPrefix(), "не создать запись о покупке подписки:", err)
		return global.ErrInternalError
	}

	fileName := param.FileData.CreateFileName(fmt.Sprintf(bot_users.BilFileNameTemplate, purchaseID))
	fullPath := filepath.Join("./uploads", fileName)

	dst, err := os.Create(fullPath)
	if err != nil {
		u.log.Db.WithFields(lf).Errorln(u.logPrefix(), "не удалось создать файл:", err)
		return global.ErrInternalError
	}
	defer dst.Close()

	_, err = io.Copy(dst, param.FileData.File)
	if err != nil {
		u.log.Db.WithFields(lf).Errorln(u.logPrefix(), "не удалось скопировать файл:", err)
		return global.ErrInternalError
	}

	user, err := u.ri.BotUsers.FindUserByID(ts, param.BotUserID)
	if err != nil {
		u.log.Db.WithFields(lf).Errorln(u.logPrefix(), "не удалось найти пользователя:", err)
		return global.ErrInternalError
	}

	_, err = u.ri.NotifyMessage.SendInviteLink(ctx, user.TgID)
	if err != nil {
		u.log.Db.WithFields(lf).Errorln("не удалось отправить ссылку на канал:", err)
		return global.ErrInternalError
	}

	return nil
}
