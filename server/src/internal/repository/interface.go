package repository

import (
	"context"
	"learn_bot_admin_panel/internal/entity/bot_users"
	"learn_bot_admin_panel/internal/entity/profile"
	"learn_bot_admin_panel/internal/entity/telegram"
	"learn_bot_admin_panel/internal/transaction"
)

type Profile interface {
	CreateProfile(ts transaction.Session, param profile.CreateProfileParam) (int, error)
	FindProfileByTGUserNameOrID(ts transaction.Session, userName string, TGID int64) (profile.User, error)
	FindProfileByID(ts transaction.Session, userID int) (profile.User, error)
	FindUserDeviceIDList(ts transaction.Session, userID int) ([]string, error)
	CreateUserDeviceID(ts transaction.Session, userID int, deviceID string) error
	SetProfilePassword(ts transaction.Session, userID int, password string) error
	SetProfileTGID(ts transaction.Session, userID int, TGID int64) error
}

type BotUsers interface {
	FindBotRegisteredUsers(ts transaction.Session, param bot_users.FindBotRegisteredUsersInnerParam) ([]bot_users.BotUserProfile, error)
	FindUserByID(ts transaction.Session, id int) (bot_users.BotUserProfile, error)
	LoadAllBotSubscriptionTypes(ts transaction.Session) ([]bot_users.BotSubscriptionType, error)
	CreateSubscriptionPurchase(ts transaction.Session, param bot_users.Purchase) (int64, error)
}

type TgBot interface {
	GetBotInfo() (telegram.BotInfoResponse, error)
}

type AuthCache interface {
	SetTempUserData(ctx context.Context, tempKey string, user profile.User) error
	GetTempUserData(ctx context.Context, tempKey string) (profile.User, error)
}

type NotifyMessage interface {
	SendInviteLink(ctx context.Context, TGID int64) (bool, error)
}
