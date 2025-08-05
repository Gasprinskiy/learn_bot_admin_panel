package usecase

import (
	"learn_bot_admin_panel/config"
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

// func (u *BotUsers) FindActiveUsers(ctx context.Context)
