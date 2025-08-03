package repository

import (
	"context"
	"learn_bot_admin_panel/internal/entity/profile"
	"learn_bot_admin_panel/internal/entity/telegram"
	"learn_bot_admin_panel/internal/transaction"
)

type Profile interface {
	CreateProfile(ts transaction.Session, param profile.CreateProfileParam) (int, error)
	FindProfileByTGUserName(ts transaction.Session, userName string) (profile.User, error)
	FindProfileByID(ts transaction.Session, userID int) (profile.User, error)
	FindUserDeviceIDList(ts transaction.Session, userID int) ([]string, error)
	CreateUserDeviceID(ts transaction.Session, userID int, deviceID string) error
	SetProfilePassword(ts transaction.Session, userID int, password string) error
	SetProfileTGID(ts transaction.Session, userID int, TGID int64) error
}

type TgBot interface {
	GetBotInfo() (telegram.BotInfoResponse, error)
}

type AuthCache interface {
	SetTempUserData(ctx context.Context, tempKey string, user profile.User) error
	GetTempUserData(ctx context.Context, tempKey string) (profile.User, error)
}
