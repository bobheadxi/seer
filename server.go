package main

import (
	"go.uber.org/zap"

	"go.bobheadxi.dev/seer/config"
	"go.bobheadxi.dev/seer/jobs"
	"go.bobheadxi.dev/seer/riot"
	"go.bobheadxi.dev/seer/server"
)

func startServer(
	log *zap.Logger,
	flags config.Flags,
	cfg config.Config,
	meta config.BuildMeta,
) {
	log.Info("instantiating dependencies")

	// riot api
	rc, err := riot.NewClient(log.Named("riot"), cfg.RiotAPIToken)
	if err != nil {
		log.Fatal(err.Error())
	}

	// storage backend
	s, err := newStorageBackend(log, flags.Store(), cfg)
	if err != nil {
		log.Fatal(err.Error())
	}

	// job queuer
	jq := jobs.NewJobQueue(log.Named("queue"), cfg.RedisNamespace, cfg.DefaultRedisPool())

	// server
	log.Info("creating server")
	srv, err := server.New(log.Named("api"), rc, s, jq, meta)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Info("spinning up server",
		zap.String("port", flags.APIPort))
	if err := srv.Start(":"+flags.APIPort, newStopper()); err != nil {
		log.Fatal(err.Error())
	}
}
