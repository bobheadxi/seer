package main

import (
	"context"
	"os"

	"go.bobheadxi.dev/seer/config"
	"go.bobheadxi.dev/seer/jobs"
	"go.bobheadxi.dev/seer/riot"
	"go.bobheadxi.dev/seer/server"
	"go.bobheadxi.dev/seer/store"
	"go.bobheadxi.dev/zapx"
	"go.uber.org/zap"
)

func main() {
	// load up configuration
	flags, err := config.LoadFlags(os.Args[1:])
	if err != nil {
		panic(err)
	}
	cfg := config.NewEnvConfig()
	meta := config.NewBuildMeta()

	// init logger
	log, err := zapx.New(flags.LogPath, flags.Dev, func(cfg *zap.Config) error {
		if flags.Dev == true {
			cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
		}
		return nil
	}, zapx.WithFields(map[string]interface{}{
		"build.commit": meta.Commit,
	}))
	if err != nil {
		panic(err)
	}
	defer log.Sync()

	// report basic config
	log.Info("configuration loaded",
		zap.Any("config", cfg),
		zap.Any("meta", meta),
		zap.Any("flags", flags))

	// spin up jobs manager UI
	if !flags.DisableJobsUI {
		go startJobsUI(cfg.RedisNamespace, cfg.DefaultRedisPool(), ":8081")
	}

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
	srv, err := server.New(log, rc, bs, je, meta)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Info("spinning up server",
		zap.String("port", flags.APIPort))
	if err := srv.Start(":"+flags.APIPort, newStopper()); err != nil {
		log.Fatal(err.Error())
	}
}
