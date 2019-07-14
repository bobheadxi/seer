package main

import (
	"context"

	"go.uber.org/zap"

	"go.bobheadxi.dev/seer/config"
	"go.bobheadxi.dev/seer/jobs"
	"go.bobheadxi.dev/seer/riot"
	"go.bobheadxi.dev/seer/server"
	"go.bobheadxi.dev/seer/store"
)

func startServer(
	log *zap.Logger,
	flags config.Flags,
	cfg config.Config,
	meta config.BuildMeta,
) {
	// spin up jobs manager UI as well if configured to do so
	if flags.JobsUIPort != "" {
		go startJobsUI(
			log.Named("jobs_ui"),
			cfg.RedisNamespace,
			cfg.DefaultRedisPool(),
			":"+flags.JobsUIPort)
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
