package main

import (
	"context"
	"flag"
	"os"

	"go.bobheadxi.dev/seer/config"
	"go.bobheadxi.dev/seer/jobs"
	"go.bobheadxi.dev/seer/riot"
	"go.bobheadxi.dev/seer/server"
	"go.bobheadxi.dev/seer/store"
	"go.bobheadxi.dev/zapx"
	"go.uber.org/zap"
)

var (
	dev     = flag.Bool("dev", os.Getenv("DEV") == "true", "toggle dev mode")
	logPath = flag.String("logpath", "", "path for log storage")
)

func main() {
	cfg := config.NewEnvConfig()

	log, err := zapx.New(*logPath, *dev, func(cfg *zap.Config) error {
		if *dev == true {
			cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	defer log.Sync()
	log.Debug("configuration loaded", zap.Any("config", cfg))

	go startJobsUI(cfg.RedisNamespace, cfg.DefaultRedisPool(), ":8081")

	log.Info("instantiating dependencies")
	rc, err := riot.NewClient(log.Named("riot"), cfg.RiotAPIToken)
	if err != nil {
		log.Fatal(err.Error())
	}
	bs, err := store.NewGitHubStore(context.Background(), log, cfg.GitHubAPITokenSource(), cfg.GitHubStoreRepo)
	if err != nil {
		log.Fatal(err.Error())
	}
	je := jobs.NewJobsEngine(log, cfg.RedisNamespace, cfg.DefaultRedisPool(), &jobs.BaseJobContext{
		RiotAPI: rc,
		Store:   bs,
	})

	log.Info("creating server")
	srv, err := server.New(log, rc, bs, je)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Info("spinning up server")
	if err := srv.Start(":8080", newStopper()); err != nil {
		log.Fatal(err.Error())
	}
}
