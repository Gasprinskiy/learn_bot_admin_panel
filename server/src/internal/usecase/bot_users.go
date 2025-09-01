package usecase

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"learn_bot_admin_panel/config"
	"learn_bot_admin_panel/internal/entity/bot_users"
	"learn_bot_admin_panel/internal/entity/chanel_kicker"
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

	file, err := excel.BuildExcelFileFromStruct(data.Data, "Зарегестрированные")
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
		status := bot_users.SubscriptionStatusNotExists
		if user.SubscrPurchaseDate.Valid {
			if user.SubscrKickTime.Valid {
				status = bot_users.SubscriptionStatusExpired
			} else {
				status = bot_users.SubscriptionStatusActive
			}
		}
		user.SetSubscriptionStatus(status)
		result.Data[i] = user
	}

	return result, nil
}

func (u *BotUsers) PrintFindUnregisteredUsers(
	ctx context.Context,
	param bot_users.FindBotUnregisteredUsersInnerParam,
) ([]byte, error) {
	data, err := u.FindUnregisteredUsers(ctx, param)
	if err != nil {
		return nil, err
	}

	file, err := excel.BuildExcelFileFromStruct(data.Data, "Не зарегестрированные")
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

func (u *BotUsers) FindUnregisteredUsers(
	ctx context.Context,
	param bot_users.FindBotUnregisteredUsersInnerParam,
) (global.CommotListSearchResponse[bot_users.BotUnregistredUserProfile], error) {
	var zero global.CommotListSearchResponse[bot_users.BotUnregistredUserProfile]

	ts := transaction.MustGetSession(ctx)

	data, err := u.ri.Repository.BotUsers.FindBotUnregisteredUsers(ts, param)
	switch err {
	case nil:
	case global.ErrNoData:
		return zero, err

	default:
		u.log.Db.Errorln(u.logPrefix(), "не удалось найти не зарегестрированных в боте пользователей:", err)
		return zero, global.ErrInternalError
	}

	result := global.NewCommotListSearchResponse(data, data[0].CommonTotalCount, param.Limit, param.PageCount)

	return result, nil
}

func (u *BotUsers) FindUserByID(ctx context.Context, id int) (bot_users.BotUserDetailData, error) {
	lf := logrus.Fields{
		"u_id": id,
	}

	var zero bot_users.BotUserDetailData

	ts := transaction.MustGetSession(ctx)

	commonData, err := u.ri.Repository.BotUsers.FindUserByID(ts, id)
	switch err {
	case nil:
	case global.ErrNoData:
		return zero, err

	default:
		u.log.Db.WithFields(lf).Errorln(u.logPrefix(), "не удалось найти пользователя:", err)
		return zero, global.ErrInternalError
	}

	purchaseData, err := u.ri.Repository.BotUsers.FindUserPurchases(ts, id)
	if err != nil && err != global.ErrNoData {
		u.log.Db.WithFields(lf).Errorln(u.logPrefix(), "не удалось найти историю покупок пользователя:", err)
		return zero, global.ErrInternalError
	}

	for i, purchase := range purchaseData {
		status := bot_users.SubscriptionStatusActive
		if purchase.KickTime.Valid {
			status = bot_users.SubscriptionStatusExpired
		}

		purchase.SetSubscriptionStatus(status)
		purchaseData[i] = purchase
	}

	return bot_users.NewBotUserDetailData(commonData, purchaseData), nil
}

func (u *BotUsers) getUserSubscriptionStatus(purchaseDate sql_null.NullTime, term int) bot_users.SubscriptionStatus {
	now := chronos.BeginingOfNow()

	if purchaseDate.Valid {
		subDate := chronos.BeginingOfDate(purchaseDate.Time)
		expireDate := subDate.AddDate(0, term, 0)

		if now.After(expireDate) {
			return bot_users.SubscriptionStatusExpired
		}

		return bot_users.SubscriptionStatusActive
	}

	return bot_users.SubscriptionStatusNotExists
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
	if param.FileData.Header.Size > global.MaxFileSize {
		return global.ErrFileSize
	}

	if _, exists := global.AbleFileExtMap[param.FileData.Ext()]; !exists {
		return global.ErrInvalidParam
	}

	lf := logrus.Fields{
		"u_id":     param.ManagerID,
		"bot_u_id": param.BotUserID,
	}

	ts := transaction.MustGetSession(ctx)
	purchase := bot_users.NewPurchase(
		param.SubID,
		param.BotUserID,
		chronos.BeginingOfNowLocal(),
		sql_null.NullFloat64{},
		sql_null.NewInt64(param.ManagerID),
		bot_users.PaymentTypeIDP2P,
	)

	purchaseID, err := u.ri.BotUsers.CreateSubscriptionPurchase(ts, purchase)
	if err != nil {
		u.log.Db.WithFields(lf).Errorln(u.logPrefix(), "не создать запись о покупке подписки:", err)
		return global.ErrInternalError
	}

	fileName := param.FileData.CreateFileName(fmt.Sprintf(bot_users.BilFileNameTemplate, purchaseID))
	fullPath := filepath.Join("../uploads", fileName)

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

	err = u.ri.BotUsers.SavePurchaseFileName(ts, int(purchaseID), fileName)
	if err != nil {
		u.log.Db.WithFields(lf).Errorln(u.logPrefix(), "не удалось сохранить название файла:", err)
		return global.ErrInternalError
	}

	user, err := u.ri.BotUsers.FindUserByID(ts, param.BotUserID)
	if err != nil {
		u.log.Db.WithFields(lf).Errorln(u.logPrefix(), "не удалось найти пользователя:", err)
		return global.ErrInternalError
	}

	sent, err := u.ri.NotifyMessage.SendInviteLink(ctx, user.TgID)
	if !sent || err != nil {
		u.log.Db.WithFields(lf).Errorln("не удалось отправить ссылку на канал", err)
		return global.ErrInternalError
	}

	return nil
}

func (u *BotUsers) CancelSubscrition(ctx context.Context, param chanel_kicker.KickUserParam) error {
	return u.ri.Kicker.KickUser(param)
}
