package main

import (
	"context"

	"go.bobheadxi.dev/seer/config"
	"go.bobheadxi.dev/seer/jobs"
	"go.bobheadxi.dev/seer/riot"
	"go.bobheadxi.dev/seer/store"
	"go.uber.org/zap"
)

func startWorker(
	log *zap.Logger,
	flags config.Flags,
	cfg config.Config,
	meta config.BuildMeta,
) {
	log.Info("instantiating dependencies")
	rc, err := riot.NewClient(log.Named("riot"), cfg.RiotAPIToken)
	if err != nil {
		log.Fatal(err.Error())
	}
	bs, err := store.NewGitHubStore(context.Background(), log, store.GitHubStoreOpts{
		Auth: cfg.GitHubAPITokenSource(), Repo: cfg.GitHubStoreRepo,
	})
	if err != nil {
		log.Fatal(err.Error())
	}
	je := jobs.NewJobsEngine(log, cfg.RedisNamespace, cfg.DefaultRedisPool(), &jobs.BaseJobContext{
		RiotAPI: rc,
		Store:   bs,
	})

	log.Info("spinning up jobs engine")
	go je.Start()
	<-newStopper()
	log.Info("stop signal received, closing")
	je.Close()
}
