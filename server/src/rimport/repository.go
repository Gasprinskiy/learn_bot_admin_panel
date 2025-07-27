package rimport

import "learn_bot_admin_panel/internal/repository"

type Repository struct {
	Profile repository.Profile
	TgBot   repository.TgBot
}
