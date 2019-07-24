package store

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"cloud.google.com/go/bigquery"
	"github.com/go-chi/chi/middleware"
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
		//return nil, fmt.Errorf("hybrid-store: failed to init github store: %v", err)
	}

	l.Info("initializing BigQuery client",
		zap.String("project", bqOpts.ProjectID),
		zap.Any("data_opts", bqOpts.DataOpts))
	bqc, err := bigquery.NewClient(ctx, bqOpts.ProjectID, bqOpts.ConnOpts...)
	if err != nil {
		//return nil, fmt.Errorf("hybrid-store: failed to initialize bigquery client: %v", err)
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
	log := s.l.With(zap.String("request.id", middleware.GetReqID(ctx)), zap.String("team.id", teamID))
	opTimer := time.Now()
	defer log.Info("add complete", zap.Duration("duration", time.Since(opTimer)))

	// convert matches into a BigQuery compatible json document (new-line demited)
	timer := time.Now()
	var buf bytes.Buffer
	for _, m := range matches {
		bytes, err := json.Marshal(&m)
		if err != nil {
			log.Error("failed to unmarshal a match",
				zap.Error(err),
				zap.Any("match", m))
		}
		buf.Write(bytes)
		buf.WriteRune('\n')
	}
	log.Info("marshal complete", zap.Duration("duration", time.Since(timer)))

	// create BiqQuery data source for upload
	timer = time.Now()
	data := bigquery.NewReaderSource(&buf)
	data.SourceFormat = bigquery.JSON
	data.AutoDetect = true

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
	log.Info("upload complete", zap.Duration("duration", time.Since(timer)))

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
