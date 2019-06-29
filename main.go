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
	dev      = flag.Bool("dev", os.Getenv("DEV") == "true", "toggle dev mode")
	logPath  = flag.String("logpath", "", "path for log storage")
	noJobsUI = flag.Bool("no-jobs-ui", false, "disable jobs UI")
	apiPort  = flag.String("port", "8080", "port to serve Seer API on")
)

func main() {
	// load up configuration
	cfg := config.NewEnvConfig()
	meta := config.NewBuildMeta()

	// init logger
	log, err := zapx.New(*logPath, *dev, func(cfg *zap.Config) error {
		if *dev == true {
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
	log.Debug("configuration loaded", zap.Any("config", cfg))

	// spin up jobs manager UI
	if !*noJobsUI {
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

	log.Info("spinning up server")
	if err := srv.Start(":"+*apiPort, newStopper()); err != nil {
		log.Fatal(err.Error())
	}
}
