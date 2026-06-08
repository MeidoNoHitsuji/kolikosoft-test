package repository

import (
	"context"

	"github.com/MeidoNoHitsuji/kolikosoft-test/internal/config"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
)

type Repository struct {
	db  *sqlx.DB
	log *zerolog.Logger
	cfg config.Config
}

func New(db *sqlx.DB, log *zerolog.Logger, cfg config.Config) *Repository {
	return &Repository{
		db:  db,
		log: log,
		cfg: cfg,
	}
}

func (r *Repository) BeginTx(ctx context.Context) (*sqlx.Tx, error) {
	return r.db.BeginTxx(ctx, nil)
}

func (r *Repository) Commit(tx *sqlx.Tx) error {
	return tx.Commit()
}

func (r *Repository) Rollback(tx *sqlx.Tx) error {
	return tx.Rollback()
}
