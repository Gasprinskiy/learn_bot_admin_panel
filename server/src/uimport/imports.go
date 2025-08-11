package uimport

import (
	"learn_bot_admin_panel/config"
	"learn_bot_admin_panel/internal/chanel_bus"
	"learn_bot_admin_panel/internal/entity/profile"
	"learn_bot_admin_panel/internal/usecase"
	"learn_bot_admin_panel/rimport"
	"learn_bot_admin_panel/tools/logger"

	"github.com/go-telegram/bot"
)

type UsecaseImport struct {
	Usecase
}

func NewUsecaseImport(
	ri *rimport.RepositoryImports,
	log *logger.Logger,
	authChan *chanel_bus.BusChanel[profile.User],
	conf *config.Config,
	b *bot.Bot,
) *UsecaseImport {
	return &UsecaseImport{
		Usecase: Usecase{
			Jwt:      usecase.NewJwt(ri, log, conf),
			Profile:  usecase.NewProfile(ri, log, authChan, conf, b),
			BotUsers: usecase.NewBotUsers(ri, log, conf, b),
		},
	}
}
