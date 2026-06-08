package db

import (
	"context"
	"embed"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

func MustMigrate(_ context.Context, db *sqlx.DB, log *zerolog.Logger) {
	if err := Migrate(db); err != nil {
		log.Fatal().Err(err).Msg("ошибка применения миграций goose")
	}

	if err := goose.Version(db.DB, "migrations"); err != nil {
		log.Fatal().Err(err).Msg("ошибка проверки версии бд")
	}

	log.Info().Msg("миграции goose применены")
}

func Migrate(db *sqlx.DB) error {
	goose.SetBaseFS(migrationsFS)
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	return goose.Up(db.DB, "migrations")
}
