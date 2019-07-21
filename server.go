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
	jq := jobs.NewJobQueue(log.Named("queue"), cfg.RedisNamespace, cfg.DefaultRedisPool())

	log.Info("creating server")
	srv, err := server.New(log.Named("api"), rc, bs, jq, meta)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Info("spinning up server",
		zap.String("port", flags.APIPort))
	if err := srv.Start(":"+flags.APIPort, newStopper()); err != nil {
		log.Fatal(err.Error())
	}
}
