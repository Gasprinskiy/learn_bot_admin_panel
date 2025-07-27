package uimport

import (
	"learn_bot_admin_panel/config"
	"learn_bot_admin_panel/internal/chanel_bus"
	"learn_bot_admin_panel/internal/usecase"
	"learn_bot_admin_panel/rimport"
	"learn_bot_admin_panel/tools/logger"
)

type UsecaseImport struct {
	Usecase
}

func NewUsecaseImport(
	ri *rimport.RepositoryImports,
	log *logger.Logger,
	authChan *chanel_bus.AuthChan,
	conf *config.Config,
) *UsecaseImport {
	return &UsecaseImport{
		Usecase: Usecase{
			Profile: usecase.NewProfile(ri, log, authChan, conf),
		},
	}
}
