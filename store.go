package main

import (
	"context"
	"fmt"

	"go.bobheadxi.dev/seer/config"
	"go.bobheadxi.dev/seer/store"
	"go.uber.org/zap"
)

func newStorageBackend(log *zap.Logger, storage config.Store, cfg config.Config) (s store.Store, err error) {
	ctx := context.Background()
	log.Info(fmt.Sprintf("configuring '%s' storage backend", storage))
	switch storage {
	case config.StoreGitHub:
		s, err = store.NewGitHubStore(ctx, log, store.GitHubStoreOpts{
			Auth: cfg.GitHubAPITokenSource(), Repo: cfg.GitHubStoreRepo,
		})
	case config.StoreHybridBigQuery:
		s, err = store.NewHybridBigQueryStore(ctx, log, store.BigQueryOpts{
			ProjectID: cfg.GCPProjectID,
			ConnOpts:  cfg.GCPConnOpts(),
			DataOpts:  cfg.BigQuery,
		}, store.GitHubStoreOpts{
			Auth: cfg.GitHubAPITokenSource(), Repo: cfg.GitHubStoreRepo,
		})
	default:
		panic(fmt.Sprintf("unsupported storage backend '%s'", storage))
	}
	return
}
