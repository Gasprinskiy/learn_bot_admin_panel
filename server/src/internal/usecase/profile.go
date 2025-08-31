package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"learn_bot_admin_panel/config"
	"learn_bot_admin_panel/internal/chanel_bus"
	"learn_bot_admin_panel/internal/entity/global"
	"learn_bot_admin_panel/internal/entity/profile"
	"learn_bot_admin_panel/internal/entity/telegram"
	"learn_bot_admin_panel/internal/transaction"
	"learn_bot_admin_panel/rimport"
	"learn_bot_admin_panel/tools/logger"
	"learn_bot_admin_panel/tools/passencoder"
	"learn_bot_admin_panel/tools/str"
	"strings"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type Profile struct {
	ri       *rimport.RepositoryImports
	log      *logger.Logger
	authChan *chanel_bus.BusChanel[profile.User]
	config   *config.Config
	b        *bot.Bot
}

func NewProfile(
	ri *rimport.RepositoryImports,
	log *logger.Logger,
	authChan *chanel_bus.BusChanel[profile.User],
	config *config.Config,
	b *bot.Bot,
) *Profile {
	return &Profile{ri, log, authChan, config, b}
}

func (u *Profile) logPrefix() string {
	return "[profile]"
}

func (u *Profile) CreateProfile(ctx context.Context, param profile.CreateProfileParam) (int64, error) {
	ts := transaction.MustGetSession(ctx)

	param.SetAccesRightID()
	ID, err := u.ri.Repository.Profile.CreateProfile(ts, param)
	if err != nil {
		u.log.Db.Errorln(u.logPrefix(), "не удалось создать новый профиль: ", err)
		return ID, global.ErrInternalError
	}

	return ID, nil
}

func (u *Profile) RedactProfile(ctx context.Context, param profile.RedactProfileParam) error {
	ts := transaction.MustGetSession(ctx)

	if err := u.ri.Repository.Profile.RedactProfile(ts, param); err != nil {
		u.log.Db.Errorln(u.logPrefix(), "не удалось отредактировать профиль: ", err)
		return global.ErrInternalError
	}

	return nil
}

func (u *Profile) CreateAuthUrlResponse() (profile.AuthUrlResponse, error) {
	var result profile.AuthUrlResponse

	botInfo, err := u.ri.Repository.TgBot.GetBotInfo()
	if err != nil {
		u.log.Db.Errorln(u.logPrefix(), "не удалось найти информацию о боте авторизации: ", err)
		return result, global.ErrInternalError
	}

	uuID := uuid.NewString()

	u.authChan.Create(uuID, u.config.SSETTL)

	result = profile.NewAuthUrlResponse(
		uuID,
		botInfo.Result.BotStartUrlWithQuery(uuID),
	)

	return result, nil
}

func (u *Profile) TwoStepTgAuthVerify(ctx context.Context, userName, callBackData string, TGID int64) (message string, err error) {
	splitted := strings.Split(callBackData, ":")
	if len(splitted) < 3 {
		return message, global.ErrInvalidParam
	}

	answer := splitted[1]
	tempID := splitted[2]
	if answer == telegram.TwoStepAuthDecline {
		var chanel chanel_bus.Chanel[profile.User]

		chanel.Error = global.ErrPermissionDenied
		done := u.authChan.Write(tempID, chanel)
		if !done {
			return message, global.ErrExpired
		}

		return telegram.TwoStepAuthDeclineMessage, nil
	}

	return u.TgAuthVerifyAnswer(ctx, userName, tempID, TGID)
}

func (u *Profile) TgAuthVerify(ctx context.Context, userName, text string, TGID int64) (message string, err error) {
	splitted := str.SplitStringByEmptySpace(text)
	if len(splitted) < 2 {
		return message, global.ErrInvalidParam
	}

	return u.TgAuthVerifyAnswer(ctx, userName, splitted[1], TGID)
}

func (u *Profile) TgAuthVerifyAnswer(ctx context.Context, userName, tempID string, TGID int64) (message string, err error) {
	var chanel chanel_bus.Chanel[profile.User]

	ts := transaction.MustGetSession(ctx)

	user, err := u.ri.Repository.Profile.FindProfileByTGUserNameOrID(ts, userName, TGID)
	switch err {
	case nil:
		if !user.IsActivated() {
			if err = u.ri.Repository.Profile.SetProfileTGID(ts, user.ID, TGID); err != nil {
				u.log.Db.Errorln(u.logPrefix(), "не удалось обновить telegram_id пользователя: ", err)
				return message, global.ErrInternalError
			}
		}

	case global.ErrNoData:
		chanel.Error = err
		err = global.ErrNoData

	default:
		u.log.Db.Errorln(u.logPrefix(), "не удалось найти пользователя по юзернейму в телеграм: ", err)
		chanel.Error = global.ErrInternalError
		err = global.ErrInternalError
	}

	chanel.Data = user
	done := u.authChan.Write(tempID, chanel)
	if !done {
		return message, global.ErrExpired
	}

	if err != nil {
		return message, err
	}

	return fmt.Sprintf(telegram.AuthSuccessfulyMessage, user.FirstName), nil
}

func (u *Profile) WaitTgAuthVerify(ctx context.Context, authKey string) ([]byte, error) {
	var data []byte

	authSession, exists := u.authChan.Read(authKey)
	if !exists {
		return data, global.ErrNoData
	}

	defer u.authChan.CleanUp(authKey)

	select {
	case <-ctx.Done():
		return data, global.ErrExpired

	case authChanel := <-authSession.Chan:
		userData := authChanel.Data

		if authChanel.Error != nil {
			return data, authChanel.Error
		}

		err := u.ri.Repository.AuthCache.SetTempUserData(ctx, authKey, userData)
		if err != nil {
			u.log.Db.Errorln("не удалось записать временные данные пользователя в кеш:", err)
			return data, global.ErrInternalError
		}

		return json.Marshal(userData.NewUserFirstLoginAnswer())

	case <-time.After(u.config.SSETTL):
		return data, global.ErrExpired
	}
}

func (u *Profile) CreateUserDeviceIDIfNotExists(ctx context.Context, userID int, deviceID string) error {
	lf := logrus.Fields{
		"u_id": userID,
	}

	ts := transaction.MustGetSession(ctx)

	deviceIDList, err := u.ri.Repository.Profile.FindUserDeviceIDList(ts, userID)
	if err != nil && err != global.ErrNoData {
		u.log.Db.WithFields(lf).Errorln(u.logPrefix(), "не удалось найти id устройств пользователя:", err)
		return global.ErrInternalError
	}

	if err := u.ri.Profile.SetProfileLastLoginDate(ts, userID, time.Now()); err != nil {
		u.log.Db.WithFields(lf).Errorln(u.logPrefix(), "не удалось обновить дату последнего входа пользователя:", err)
		return global.ErrInternalError
	}

	idMap := make(map[string]struct{}, len(deviceIDList))

	for _, id := range deviceIDList {
		idMap[id] = struct{}{}
	}

	_, exists := idMap[deviceID]
	if exists {
		return nil
	}

	if err = u.ri.Repository.Profile.CreateUserDeviceID(ts, userID, deviceID); err != nil {
		u.log.Db.WithFields(lf).Errorln(u.logPrefix(), "не удалось создать id нового устройства пользователя:", err)
		return global.ErrInternalError
	}

	return nil
}

func (u *Profile) OnPasswordLogin(ctx context.Context, param profile.PasswordLoginParam, deviceID string) (profile.PasswordLoginResponse, error) {
	lf := logrus.Fields{
		"tg_user_name": param.UserName,
	}

	var zero profile.PasswordLoginResponse

	ts := transaction.MustGetSession(ctx)

	userInfo, err := u.ri.Repository.Profile.FindProfileByTGUserNameOrID(ts, param.UserName, 0)
	switch err {
	case nil:
	case global.ErrNoData:
		return zero, err

	default:
		u.log.Db.WithFields(lf).Errorln(u.logPrefix(), "не удалось найти пользователя по телеграм юзернейму:", err)
		return zero, global.ErrInternalError
	}

	if !userInfo.IsPasswordSet() {
		return zero, global.ErrNoData
	}

	if err := u.ri.Profile.SetProfileLastLoginDate(ts, userInfo.ID, time.Now()); err != nil {
		u.log.Db.WithFields(lf).Errorln(u.logPrefix(), "не удалось обновить дату последнего входа пользователя:", err)
		return zero, global.ErrInternalError
	}

	valid := passencoder.CheckHashPassword(userInfo.Password.String, param.Password)
	if !valid {
		return zero, global.ErrNoData
	}

	lf["u_id"] = userInfo.ID

	deviceIDList, err := u.ri.Repository.Profile.FindUserDeviceIDList(ts, userInfo.ID)
	if err != nil && err != global.ErrNoData {
		u.log.Db.WithFields(lf).Errorln(u.logPrefix(), "не удалось найти id устройств пользователя:", err)
		return zero, global.ErrInternalError
	}

	var (
		needTwoStepAuth bool
		uuID            string
	)

	if len(deviceIDList) > 0 {
		needTwoStepAuth = true

		idMap := make(map[string]struct{}, len(deviceIDList))

		for _, id := range deviceIDList {
			idMap[id] = struct{}{}
		}

		_, exists := idMap[deviceID]

		if !exists {
			uuID = uuid.NewString()
			u.authChan.Create(uuID, u.config.SSETTL)

			_, err := u.b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID:    userInfo.TgID.Int64,
				Text:      fmt.Sprintf(telegram.TwoStepAuthMessage, userInfo.FirstName),
				ParseMode: "HTML",
				ReplyMarkup: &models.InlineKeyboardMarkup{
					InlineKeyboard: [][]models.InlineKeyboardButton{
						{
							{
								Text:         telegram.TwoStepAuthCallBackQueryDeclineView,
								CallbackData: telegram.TwoStepAuthCallBackQueryDecline(uuID),
							},
							{
								Text:         telegram.TwoStepAuthCallBackQueryAcceptView,
								CallbackData: telegram.TwoStepAuthCallBackQueryAccept(uuID),
							},
						},
					},
				},
			})

			if err != nil {
				u.log.Db.WithFields(lf).Errorln(u.logPrefix(), "не удалось отправить сообщения о подтверждении", err)
				return zero, global.ErrInternalError
			}
		}

		needTwoStepAuth = !exists
	}

	return profile.NewPasswordLoginResponse(needTwoStepAuth, userInfo.ID, userInfo.Access, uuID), nil
}

func (u *Profile) SetProfilePassword(ctx context.Context, password string, userID int) error {
	lf := logrus.Fields{
		"u_id": userID,
	}

	hashPassword, err := passencoder.CreateHashPassword(password)
	if err != nil {
		u.log.Db.WithFields(lf).Errorln(u.logPrefix(), "не удалось захешировать пароль:", err)
		return global.ErrInternalError
	}

	ts := transaction.MustGetSession(ctx)

	err = u.ri.Repository.Profile.SetProfilePassword(ts, userID, hashPassword)
	if err != nil {
		u.log.Db.WithFields(lf).Errorln(u.logPrefix(), "не удалось записать пароль:", err)
		return global.ErrInternalError
	}

	return nil
}

func (u *Profile) GetUserCommonInfo(ctx context.Context, userID int) (profile.UserCommonInfo, error) {
	var zero profile.UserCommonInfo

	lf := logrus.Fields{
		"u_id": userID,
	}

	ts := transaction.MustGetSession(ctx)

	user, err := u.ri.Repository.Profile.FindProfileByID(ts, userID)
	switch err {
	case nil:
	case global.ErrNoData:
		return zero, err

	default:
		u.log.Db.WithFields(lf).Errorln(u.logPrefix(), "не удалось найти пользователя по id:", err)
		return zero, global.ErrInternalError
	}

	return user.NewUserCommonInfo(), nil
}

func (u *Profile) FindProfileUsers(ctx context.Context, userID int) ([]profile.User, error) {
	lf := logrus.Fields{
		"u_id": userID,
	}

	ts := transaction.MustGetSession(ctx)

	data, err := u.ri.Profile.LoadUsersProfile(ts)
	switch err {
	case nil:
	case global.ErrNoData:
		return nil, err

	default:
		u.log.Db.WithFields(lf).Errorln(u.logPrefix(), "не удалось найти пользователей:", err)
		return nil, global.ErrInternalError
	}

	for i, user := range data {
		if user.ID == userID {
			data[i].IsYou = true
		}
	}

	return data, nil
}
