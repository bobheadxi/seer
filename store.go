package main

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"go.bobheadxi.dev/seer/config"
	"go.bobheadxi.dev/seer/store"
	"go.bobheadxi.dev/seer/store/cache"
)

func newStorageBackend(
	log *zap.Logger,
	cachePath string,
	storage config.Store,
	cfg config.Config,
	md config.BuildMeta,
) (s store.Store, err error) {
	ctx := context.Background()
	log.Info(fmt.Sprintf("configuring '%s' storage backend", storage))
	switch storage {
	case config.StoreBigQuery:
		s, err = store.NewBigQueryStore(ctx, log, store.BigQueryOpts{
			ServiceVersion: md.Version,
			ProjectID:      cfg.GCPProjectID,
			ConnOpts:       cfg.GCPConnOpts(),
			DataOpts:       cfg.BigQuery,
		})
	default:
		log.Fatal(fmt.Sprintf("unsupported storage backend '%s'", storage))
	}

	if cachePath != "" {
		s, err = cache.New(log.Named("cache"), s, cfg.DefaultRedisPool())
		if err != nil {
			log.Fatal("unable to start up cache", zap.Error(err))
		}
	}

	return
}
