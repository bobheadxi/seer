package main

import (
	"go.bobheadxi.dev/seer/config"
	"go.bobheadxi.dev/seer/jobs"
	"go.bobheadxi.dev/seer/riot"
	"go.uber.org/zap"
)

func startWorker(
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
	s, err := newStorageBackend(log, flags.CachePath, flags.Store(), cfg, meta)
	if err != nil {
		log.Fatal(err.Error())
	}
	if err != nil {
		log.Fatal(err.Error())
	}

	// jobs engine
	je := jobs.NewJobsEngine(log, cfg.RedisNamespace, cfg.DefaultRedisPool(), &jobs.BaseJobContext{
		RiotAPI: rc,
		Store:   s,
	})
	log.Info("spinning up jobs engine")
	go je.Start()
	<-newStopper()
	log.Info("stop signal received, closing")
	je.Close()
}
