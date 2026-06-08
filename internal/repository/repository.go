package repository

import (
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
