package store

import (
	"context"

	"cloud.google.com/go/bigquery"
	"go.uber.org/zap"
)

type hybridBigQueryStore struct {
	gh Store
	bq *bigquery.Client
}

// BigQueryOpts defines options for a BigQuery connection
type BigQueryOpts struct {
	ProjectID string
}

// NewHybridBigQueryStore instantiates a new Store backed by GitHub issues and
// Google BigQuery https://cloud.google.com/bigquery/
func NewHybridBigQueryStore(ctx context.Context, l *zap.Logger, bqOpts BigQueryOpts, ghOpts GitHubStoreOpts) (Store, error) {
	ghStore, err := NewGitHubStore(ctx, l.Named("gh"), ghOpts)
	if err != nil {
		return nil, err
	}

	bqc, err := bigquery.NewClient(ctx, bqOpts.ProjectID)
	if err != nil {
		return nil, err
	}

	return &hybridBigQueryStore{
		gh: ghStore,
		bq: bqc,
	}, nil
}

func (s *hybridBigQueryStore) Create(ctx context.Context, teamID string, team *Team) error {
	return s.gh.Create(ctx, teamID, team)
}

func (s *hybridBigQueryStore) GetTeam(ctx context.Context, teamID string) (*Team, error) {
	return s.gh.GetTeam(ctx, teamID)
}

func (s *hybridBigQueryStore) GetMatches(ctx context.Context, teamID string) (Matches, error) {
	// TODO: get matches from BigQuery
	return nil, nil
}

func (s *hybridBigQueryStore) Add(ctx context.Context, teamID string, matches Matches) error {
	// TODO: add matches to BigQuery, update analytics in GitHub team issue
	return nil
}
