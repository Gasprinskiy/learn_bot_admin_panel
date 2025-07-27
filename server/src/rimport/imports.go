package rimport

import (
	"learn_bot_admin_panel/config"
	http_rep "learn_bot_admin_panel/internal/repository/http"
	"learn_bot_admin_panel/internal/repository/postgres"

	"github.com/redis/go-redis/v9"
)

type RepositoryImports struct {
	Repository
}

func NewRepositoryImports(config *config.Config, rdb *redis.Client) *RepositoryImports {
	return &RepositoryImports{
		Repository: Repository{
			// UserCache: redis_cache.NewUserCache(rdb, config.RedisTtl),
			Profile: postgres.NewProfile(),
			TgBot:   http_rep.NewTgBot(config.TgApiURL, config.BotToken),
		},
	}
}
