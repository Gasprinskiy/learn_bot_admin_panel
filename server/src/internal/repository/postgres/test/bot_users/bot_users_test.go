package bot_users_test

import (
	"fmt"
	"learn_bot_admin_panel/config"
	"learn_bot_admin_panel/internal/entity/bot_users"
	"learn_bot_admin_panel/internal/repository/postgres"
	"learn_bot_admin_panel/internal/transaction"
	"learn_bot_admin_panel/tools/dump"
	"log"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

func TestFindBotRegisteredUsers(t *testing.T) {
	r := require.New(t)

	config := config.NewConfig()

	pgdb, err := sqlx.Connect("pgx", config.PostgresURL)
	if err != nil {
		log.Fatalln("не удалось подключиться к базе postgres: ", err)
	}
	defer pgdb.Close()

	r.NoError(pgdb.Ping())

	repo := postgres.NewBotUsers()

	ts := transaction.NewSQLSession(pgdb)
	err = ts.Start()
	r.NoError(err)
	defer ts.Rollback()

	t.Run("freaking_test", func(t *testing.T) {
		data, err := repo.FindBotRegisteredUsers(ts, bot_users.FindBotRegisteredUsersParam{
			Limit: 10,
		})
		r.NoError(err)

		fmt.Println("DATA: ", dump.Struct(data))
	})
}
