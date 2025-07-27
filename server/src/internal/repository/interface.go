package repository

import (
	"learn_bot_admin_panel/internal/entity/profile"
	"learn_bot_admin_panel/internal/entity/telegram"
	"learn_bot_admin_panel/internal/transaction"
)

type Profile interface {
	CreateProfile(ts transaction.Session, param profile.CreateProfileParam) (int, error)
	FindProfileByTGUserName(ts transaction.Session, userName string) (profile.User, error)
	// SetProfilePassword(ts transaction.Session, ID int, password string) error
}

type TgBot interface {
	GetBotInfo() (telegram.BotInfoResponse, error)
}
