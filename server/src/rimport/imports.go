package rimport

import (
	"learn_bot_admin_panel/config"
	"learn_bot_admin_panel/internal/repository/grpc_client"
	http_rep "learn_bot_admin_panel/internal/repository/http"
	"learn_bot_admin_panel/internal/repository/postgres"
	"learn_bot_admin_panel/internal/repository/redis_cache"

	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
)

type RepositoryImports struct {
	Repository
}

func NewRepositoryImports(
	config *config.Config,
	rdb *redis.Client,
	grpcConn *grpc.ClientConn,
) *RepositoryImports {
	return &RepositoryImports{
		Repository: Repository{
			AuthCache:     redis_cache.NewAuthCache(rdb, config.SSETTL),
			TgBot:         http_rep.NewTgBot(config.TgApiURL, config.BotToken),
			Profile:       postgres.NewProfile(),
			BotUsers:      postgres.NewBotUsers(),
			NotifyMessage: grpc_client.NewnotifyMessageGRPCRepository(grpcConn),
			Kicker:        grpc_client.NewKicker(grpcConn),
		},
	}
}
