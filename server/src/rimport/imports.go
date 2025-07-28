package rimport

import (
	"learn_bot_admin_panel/config"
	http_rep "learn_bot_admin_panel/internal/repository/http"
	"learn_bot_admin_panel/internal/repository/postgres"
	"learn_bot_admin_panel/internal/repository/redis_cache"

	"github.com/redis/go-redis/v9"
)

type RepositoryImports struct {
	Repository
}

func NewRepositoryImports(config *config.Config, rdb *redis.Client) *RepositoryImports {
	return &RepositoryImports{
		Repository: Repository{
			AuthCache: redis_cache.NewAuthCache(rdb, config.SSETTL),
			Profile:   postgres.NewProfile(),
			TgBot:     http_rep.NewTgBot(config.TgApiURL, config.BotToken),
		},
	}
}
