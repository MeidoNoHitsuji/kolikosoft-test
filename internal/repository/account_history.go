package repository

import (
	"context"

	backErr "github.com/MeidoNoHitsuji/kolikosoft-test/internal/errors"
	"github.com/MeidoNoHitsuji/kolikosoft-test/internal/model"
	"github.com/jmoiron/sqlx"
)

func (r *Repository) TxAddHistory(ctx context.Context, tx *sqlx.Tx, data *model.AddHistory) error {
	const query = `
		INSERT INTO account_histories (
			account_id,
			old_balance,
			new_balance,
			comment
		)
		VALUES (
			$1,
			$2,
			$3,
			$4
		)
	`

	if _, err := tx.ExecContext(
		ctx,
		query,
		data.AccountID,
		data.OldBalance,
		data.NewBalance,
		data.Comment,
	); err != nil {
		r.log.Err(err).Msg("Не удалось добавить историю баланса")
		return backErr.ErrInternalServer
	}

	return nil
}
