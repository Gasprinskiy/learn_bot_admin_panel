package postgres

import (
	"learn_bot_admin_panel/internal/transaction"

	"github.com/jmoiron/sqlx"
)

func SqlxTx(ts transaction.Session) *sqlx.Tx {
	tx := ts.Tx()
	return tx.(*sqlx.Tx)
}
