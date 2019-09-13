package main

import (
	"context"
	"fmt"
	"os"

	"github.com/dgraph-io/badger"
	"go.uber.org/zap"

	"go.bobheadxi.dev/seer/config"
	"go.bobheadxi.dev/seer/store"
	"go.bobheadxi.dev/seer/store/cache"
	"go.bobheadxi.dev/zapx/zapx"
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
		log.Error(fmt.Sprintf("unsupported storage backend '%s'", storage))
		os.Exit(1)
	}

	if cachePath != "" {
		badgerLogger := log.Named("badger.db").
			WithOptions(zapx.WrapWithLevel(zap.ErrorLevel))
		opts := badger.DefaultOptions(cachePath).
			WithLogger(zapx.NewFormatLogger(badgerLogger))
		db, err := badger.Open(opts)
		if err != nil {
			log.Error(fmt.Sprintf("unable to open cache at '%s': %v", cachePath, err))
			os.Exit(1)
		}
		s = cache.New(log.Named("cache"), s, db)
	}

	return
}
