package main

import (
	"context"
	"flag"

	"github.com/MeidoNoHitsuji/kolikosoft-test/internal/db"
	"github.com/MeidoNoHitsuji/kolikosoft-test/internal/layer"
	"github.com/MeidoNoHitsuji/kolikosoft-test/internal/log"
	"github.com/MeidoNoHitsuji/kolikosoft-test/internal/shutdowner"
)

func main() {
	ctx := context.Background()

	cfg := layer.NewConfig(parseFlags())
	log.InitLogger(&cfg)

	sd := shutdowner.NewShutdowner(ctx, &log.Logger)

	database := db.NewDatabase(ctx, cfg, &log.Logger)
	db.MustMigrate(ctx, database, &log.Logger)
	db.Shutdown(sd, database)

	//rep := repository.New(database, &log.Logger, cfg)

	sd.WaitShutdown()
}

func parseFlags() string {
	var configPath string
	flag.StringVar(&configPath, "config", "./config/application.json", "path to config file")
	flag.Parse()

	return configPath
}
