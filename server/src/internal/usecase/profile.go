package usecase

import (
	"context"
	"encoding/json"
	"learn_bot_admin_panel/config"
	"learn_bot_admin_panel/internal/chanel_bus"
	"learn_bot_admin_panel/internal/entity/global"
	"learn_bot_admin_panel/internal/entity/profile"
	"learn_bot_admin_panel/internal/transaction"
	"learn_bot_admin_panel/rimport"
	"learn_bot_admin_panel/tools/logger"
	"learn_bot_admin_panel/tools/passencoder"
	"learn_bot_admin_panel/tools/str"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type Profile struct {
	ri       *rimport.RepositoryImports
	log      *logger.Logger
	authChan *chanel_bus.AuthChan
	config   *config.Config
}

func NewProfile(
	ri *rimport.RepositoryImports,
	log *logger.Logger,
	authChan *chanel_bus.AuthChan,
	config *config.Config,
) *Profile {
	return &Profile{ri, log, authChan, config}
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

func (u *Profile) TgAuthVerify(ctx context.Context, userName, text string) (message string, err error) {
	var chanel chanel_bus.SessionChanel

	splitted := str.SplitStringByEmptySpace(text)
	if len(splitted) < 2 {
		return message, global.ErrInvalidParam
	}

	ts := transaction.MustGetSession(ctx)

	user, err := u.ri.Repository.Profile.FindProfileByTGUserName(ts, userName)
	switch err {
	case nil:
	case global.ErrNoData:
		chanel.Error = err
		err = global.ErrNoData

	default:
		u.log.Db.Errorln(u.logPrefix(), "не удалось найти пользователя по юзернейму в телеграм: ", err)
		chanel.Error = global.ErrInternalError
		err = global.ErrInternalError
	}

	chanel.User = user
	done := u.authChan.Write(splitted[1], chanel)
	if !done {
		return message, global.ErrExpired
	}

	if err != nil {
		return message, err
	}

	return profile.AuthSuccessfulyMessage, nil
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
		userData := authChanel.User

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

	idMap := make(map[string]struct{}, len(deviceIDList))

	for _, id := range deviceIDList {
		idMap[id] = struct{}{}
	}

	_, exists := idMap[deviceID]
	if !exists {
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

	userInfo, err := u.ri.Repository.Profile.FindProfileByTGUserName(ts, param.UserName)
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

	var needTwoStepAuth bool

	if len(deviceIDList) > 0 {
		needTwoStepAuth = true

		idMap := make(map[string]struct{}, len(deviceIDList))

		for _, id := range deviceIDList {
			idMap[id] = struct{}{}
		}

		_, exists := idMap[deviceID]

		needTwoStepAuth = !exists
	}

	return profile.NewPasswordLoginResponse(needTwoStepAuth, userInfo.ID, userInfo.Access), nil
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
