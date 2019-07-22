package store

import (
	"bytes"
	"context"
	"fmt"

	"cloud.google.com/go/bigquery"
	"github.com/gocarina/gocsv"
	"go.uber.org/zap"
	"google.golang.org/api/option"

	"go.bobheadxi.dev/seer/config"
)

type hybridBigQueryStore struct {
	l *zap.Logger

	gh Store
	bq *bigquery.Client

	cfg config.BigQuery
}

// BigQueryOpts defines options for a BigQuery connection
type BigQueryOpts struct {
	ProjectID string
	ConnOpts  []option.ClientOption

	DataOpts config.BigQuery
}

// NewHybridBigQueryStore instantiates a new Store backed by GitHub issues and
// Google BigQuery https://cloud.google.com/bigquery/
func NewHybridBigQueryStore(ctx context.Context, l *zap.Logger, bqOpts BigQueryOpts, ghOpts GitHubStoreOpts) (Store, error) {
	ghStore, err := NewGitHubStore(ctx, l.Named("gh"), ghOpts)
	if err != nil {
		return nil, fmt.Errorf("hybrid-store: failed to init github store: %v", err)
	}

	l.Info("initializing BigQuery client",
		zap.String("project", bqOpts.ProjectID),
		zap.Any("data_opts", bqOpts.DataOpts))
	bqc, err := bigquery.NewClient(ctx, bqOpts.ProjectID, bqOpts.ConnOpts...)
	if err != nil {
		return nil, fmt.Errorf("hybrid-store: failed to initialize bigquery client: %v", err)
	}

	return &hybridBigQueryStore{
		l: l,

		gh: ghStore,
		bq: bqc,

		cfg: bqOpts.DataOpts,
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
	// convert matches into a BiqQuery data source for upload
	var buf bytes.Buffer
	if err := gocsv.Marshal(matches, &buf); err != nil {
		return fmt.Errorf("failed to marshal matches as csv: %v", err)
	}
	data := bigquery.NewReaderSource(&buf)
	data.SourceFormat = bigquery.CSV
	data.SkipLeadingRows = 1
	data.AutoDetect = true
	// TODO: it seems like a schema must be defined qq for nested objects.
	// maybe consider multiple tables and join?

	// run upload
	// TODO: there are somewhat strict limitations on this: https://cloud.google.com/bigquery/quotas#load_jobs
	// consider pooling multiple team's new matches together
	loader := s.bqMatchesTable().LoaderFrom(data)
	job, err := loader.Run(ctx)
	if err != nil {
		return fmt.Errorf("could not run data load: %v", err)
	}
	status, err := job.Wait(ctx)
	if err != nil {
		return fmt.Errorf("data load job failed to complete: %v", err)
	}
	if err := status.Err(); err != nil {
		return fmt.Errorf("data load job completed with error: %v", err)
	}

	// TODO: run and update analytics in GitHub team issue. or use a view:
	// https://cloud.google.com/bigquery/docs/views

	return nil
}

func (s *hybridBigQueryStore) Close() error { return s.Close() }

func (s *hybridBigQueryStore) bqDataset() *bigquery.Dataset {
	return s.bq.Dataset(s.cfg.DatasetID)
}

func (s *hybridBigQueryStore) bqMatchesTable() *bigquery.Table {
	return s.bqDataset().Table(s.cfg.MatchesTableID)
}
