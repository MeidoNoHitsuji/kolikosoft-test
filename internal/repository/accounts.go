package repository

import (
	"context"
	"database/sql"
	"errors"

	backErr "github.com/MeidoNoHitsuji/kolikosoft-test/internal/errors"
	"github.com/MeidoNoHitsuji/kolikosoft-test/internal/model"
	"github.com/jmoiron/sqlx"
)

func (r *Repository) GetAccountByID(ctx context.Context, id int64) (*model.Account, error) {
	const query = `SELECT
			id,
			balance
		FROM accounts
		WHERE id = $1
	`

	var account model.Account

	if err := r.db.GetContext(ctx, &account, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, backErr.ErrAccountNotFound
		}

		return nil, err
	}

	return &account, nil
}

func (r *Repository) GetAccountHistoriesByAccountID(ctx context.Context, accountID int64) ([]model.AccountHistory, error) {
	const query = `
		SELECT
			id,
			account_id,
			old_balance,
			new_balance,
			comment,
			created_at
		FROM account_histories
		WHERE account_id = $1
		ORDER BY created_at DESC
		LIMIT 100
	`

	var histories []model.AccountHistory
	if err := r.db.SelectContext(ctx, &histories, query, accountID); err != nil {
		return nil, err
	}

	return histories, nil
}

func (r *Repository) TxUpdateBalance(ctx context.Context, tx *sqlx.Tx, data *model.AccountAddMoney) (*model.AddHistory, error) {
	const query = `
		UPDATE accounts
		SET balance = balance + $2
		WHERE id = $1::bigint
		RETURNING
			id AS account_id,
			balance - $2 AS old_balance,
			balance AS new_balance
	`

	var result model.AddHistory

	if err := tx.GetContext(ctx, &result, query, data.ID, data.Value); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, backErr.ErrAccountNotFound
		}
		r.log.Err(err).Msg("Не удалось обновить баланс")
		return nil, backErr.ErrInternalServer
	}

	result.Comment = "Пополнение баланса"
	return &result, nil
}
