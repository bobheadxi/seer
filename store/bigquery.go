package store

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"cloud.google.com/go/bigquery"
	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"

	"go.bobheadxi.dev/seer/config"
	"go.bobheadxi.dev/seer/riot"
	"go.bobheadxi.dev/seer/store/bigqueries"
)

type bigQueryStore struct {
	l  *zap.Logger
	bq *bigquery.Client

	project string
	cfg     config.BigQuery
}

// BigQueryOpts defines options for a BigQuery connection
type BigQueryOpts struct {
	ProjectID string
	ConnOpts  []option.ClientOption

	DataOpts config.BigQuery
}

// NewBigQueryStore instantiates a new Store backed by GitHub issues and
// Google BigQuery https://cloud.google.com/bigquery/
func NewBigQueryStore(ctx context.Context, l *zap.Logger, bqOpts BigQueryOpts) (Store, error) {
	l.Info("initializing BigQuery client",
		zap.String("project", bqOpts.ProjectID),
		zap.Any("data_opts", bqOpts.DataOpts))
	bqc, err := bigquery.NewClient(ctx, bqOpts.ProjectID, bqOpts.ConnOpts...)
	if err != nil {
		return nil, fmt.Errorf("hybrid-store: failed to initialize bigquery client: %v", err)
	}

	return &bigQueryStore{
		l:  l,
		bq: bqc,

		project: bqOpts.ProjectID,
		cfg:     bqOpts.DataOpts,
	}, nil
}

func (s *bigQueryStore) Create(ctx context.Context, teamID string, team *Team) error {
	log := s.l.With(zap.String("request.id", middleware.GetReqID(ctx)), zap.String("team.id", teamID))

	// format data
	members := make([]string, len(team.Members))
	for i, m := range team.Members {
		members[i] = fmt.Sprintf("'%s'", m.AccountID)
	}
	teamBytes, err := json.Marshal(team.Members)
	if err != nil {
		return err
	}

	// set up table view
	rawQuery, err := bigqueries.ReadFile(bigqueries.TeamMatchesView)
	if err != nil {
		return err
	}
	view := &bigquery.TableMetadata{
		Name:        "", // TODO?
		Description: string(teamBytes),
		Labels: map[string]string{
			"team":   teamID,
			"region": team.Region.ToLower(),
		},
		ViewQuery: fmt.Sprintf(string(rawQuery),
			strings.Join(members, ","),
			s.project,
			s.cfg.DatasetID,
			s.cfg.MatchesTableID),
	}
	log.Debug("view configuration instantiated", zap.Any("view_configuration", view))

	// create view
	log.Info("creating team view in BigQuery",
		zap.String("table_id", teamView(teamID)))
	if err := s.bqDataset().Table(teamView(teamID)).Create(ctx, view); err != nil {
		return fmt.Errorf("unable to create team: %v", err)
	}

	return nil
}

func (s *bigQueryStore) GetTeam(ctx context.Context, teamID string) (*Team, error) {
	t := s.bqTeamView(teamID)
	md, err := t.Metadata(ctx)
	if err != nil {
		return nil, err
	}
	if md.Type != bigquery.ViewTable {
		return nil, fmt.Errorf("expected table to be a view, go '%s'", md.Type)
	}

	var members []*riot.Summoner
	if err := json.Unmarshal([]byte(md.Description), &members); err != nil {
		return nil, err
	}

	return &Team{
		Region:    riot.ParseRegion(md.Labels["region"]),
		Members:   members,
		Analytics: nil, // TODO
	}, nil
}

func (s *bigQueryStore) GetMatches(ctx context.Context, teamID string) ([]int64, error) {
	rawQuery, err := bigqueries.ReadFile(bigqueries.TeamGamesQuery)
	if err != nil {
		return nil, err
	}
	query := s.bq.Query(fmt.Sprintf(string(rawQuery),
		s.project,
		s.cfg.DatasetID,
		teamView(teamID)))

	it, err := query.Read(ctx)
	if err != nil {
		return nil, err
	}
	var matches []int64
	for {
		var values []bigquery.Value
		err := it.Next(&values)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		id, ok := values[0].(int64)
		if !ok {
			return nil, fmt.Errorf("unknown type %T", ok)
		}
		matches = append(matches, id)
	}
	return nil, nil
}

func (s *bigQueryStore) Add(ctx context.Context, teamID string, matches Matches) error {
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
	loader.Labels = map[string]string{"team": teamID}
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

	return nil
}

func (s *bigQueryStore) Close() error { return s.Close() }

func (s *bigQueryStore) bqDataset() *bigquery.Dataset { return s.bq.Dataset(s.cfg.DatasetID) }

func (s *bigQueryStore) bqMatchesTable() *bigquery.Table {
	return s.bqDataset().Table(s.cfg.MatchesTableID)
}

func (s *bigQueryStore) bqTeamView(teamID string) *bigquery.Table {
	return s.bqDataset().Table(teamView(teamID))
}

func teamView(teamID string) string { return fmt.Sprintf("team_%s", teamID) }
