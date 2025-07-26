package uimport

import (
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
) *UsecaseImport {
	return &UsecaseImport{
		Usecase: Usecase{
			Profile: usecase.NewProfile(ri, log),
		},
	}
}
