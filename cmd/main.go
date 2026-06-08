package main

import (
	"context"
	"flag"

	"github.com/MeidoNoHitsuji/kolikosoft-test/internal/controller"
	"github.com/MeidoNoHitsuji/kolikosoft-test/internal/db"
	"github.com/MeidoNoHitsuji/kolikosoft-test/internal/layer"
	"github.com/MeidoNoHitsuji/kolikosoft-test/internal/log"
	"github.com/MeidoNoHitsuji/kolikosoft-test/internal/repository"
	"github.com/MeidoNoHitsuji/kolikosoft-test/internal/service"
	"github.com/MeidoNoHitsuji/kolikosoft-test/internal/shutdowner"
	"github.com/MeidoNoHitsuji/kolikosoft-test/internal/web"
)

func main() {
	ctx := context.Background()

	cfg := layer.NewConfig(parseFlags())
	log.InitLogger(&cfg)

	sd := shutdowner.NewShutdowner(ctx, &log.Logger)

	database := db.NewDatabase(ctx, cfg, &log.Logger)
	db.MustMigrate(ctx, database, &log.Logger)
	db.Shutdown(sd, database)

	rep := repository.New(database, &log.Logger, cfg)

	srvHolder := service.NewServiceHolder(rep)
	ctrHolder := controller.NewControllerHolder(srvHolder)

	serve := web.NewServe(ctrHolder, &log.Logger, cfg, sd)
	serve.Run()
	serve.Shutdown(ctx)

	sd.WaitShutdown()
}

func parseFlags() string {
	var configPath string
	flag.StringVar(&configPath, "config", "./config/application.json", "path to config file")
	flag.Parse()

	return configPath
}
