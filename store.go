package main

import (
	"context"
	"fmt"
	"os"

	"go.bobheadxi.dev/seer/config"
	"go.bobheadxi.dev/seer/store"
	"go.uber.org/zap"
)

func newStorageBackend(log *zap.Logger, storage config.Store, cfg config.Config, md config.BuildMeta) (s store.Store, err error) {
	ctx := context.Background()
	log.Info(fmt.Sprintf("configuring '%s' storage backend", storage))
	switch storage {
	case config.StoreGitHub:
		s, err = store.NewGitHubStore(ctx, log, store.GitHubStoreOpts{
			Auth: cfg.GitHubAPITokenSource(), Repo: cfg.GitHubStoreRepo,
		})
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
	return
}
