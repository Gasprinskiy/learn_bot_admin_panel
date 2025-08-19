package rimport

import "learn_bot_admin_panel/internal/repository"

type Repository struct {
	AuthCache     repository.AuthCache
	Profile       repository.Profile
	TgBot         repository.TgBot
	BotUsers      repository.BotUsers
	NotifyMessage repository.NotifyMessage
}
