package rimport

import (
	"learn_bot_admin_panel/config"

	"github.com/redis/go-redis/v9"
)

type RepositoryImports struct {
	Repository
}

func NewRepositoryImports(config *config.Config, rdb *redis.Client) *RepositoryImports {
	return &RepositoryImports{
		Repository: Repository{
			// UserCache: redis_cache.NewUserCache(rdb, config.RedisTtl),
			// Profile: postgres.NewProfile(),
		},
	}
}
