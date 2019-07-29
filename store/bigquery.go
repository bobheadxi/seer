package store

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
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
	opTimer := time.Now()
	defer logDuration(log, "team create complete", "Create.duration", opTimer)

	// format data
	members := make([]string, len(team.Members))
	for i, m := range team.Members {
		members[i] = fmt.Sprintf("'%s'", m.AccountID)
	}
	membersBytes, err := json.Marshal(team.Members)
	if err != nil {
		return err
	}

	// set up table view
	rawQuery, err := bigqueries.ReadFile(bigqueries.TeamMatchesView)
	if err != nil {
		return err
	}
	view := &bigquery.TableMetadata{
		Name:        "", // TODO? this could be pretty-name
		Description: string(membersBytes),
		Labels: map[string]string{
			"team":   teamID,
			"region": team.Region.ToLower(),
		},
		ViewQuery: fmt.Sprintf(string(rawQuery),
			strings.Join(members, ","),
			s.project,
			s.cfg.DatasetID,
			s.cfg.MatchesTableID),
		ExpirationTime: time.Now().Add(365 * 24 * time.Hour),
	}
	log.Debug("view configuration instantiated", zap.Any("view_configuration", view))

	// create view
	log.Info("creating team view in BigQuery",
		zap.String("table_id", teamView(teamID)),
		zap.Strings("members", members))
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
	log := s.l.With(zap.String("request.id", middleware.GetReqID(ctx)), zap.String("team.id", teamID))
	opTimer := time.Now()
	defer logDuration(log, "matches get complete", "GetMatches.duration", opTimer)

	rawQuery, err := bigqueries.ReadFile(bigqueries.TeamGamesQuery)
	if err != nil {
		return nil, err
	}
	query := s.bq.Query(fmt.Sprintf(string(rawQuery),
		s.project,
		s.cfg.DatasetID,
		teamView(teamID)))

	// execute query
	job, stats, err := s.execQuery(ctx, query)
	if err != nil {
		return nil, err
	}
	log.Info("query for matches complete",
		zap.Bool("cache_hit", stats.CacheHit),
		zap.Int64("bytes_billed", stats.TotalBytesBilled))

	// parse results
	it, err := job.Read(ctx)
	if err != nil {
		return nil, err
	}
	var matches []int64
	for {
		var values []bigquery.Value
		if err := it.Next(&values); err == iterator.Done {
			break
		} else if err != nil {
			return nil, fmt.Errorf("failed to fetch next row: %+v", err)
		}
		if len(values) == 0 {
			return nil, errors.New("unexpected empty value in query response")
		}
		id, ok := values[0].(int64)
		if !ok {
			return nil, fmt.Errorf("unknown type %T", ok)
		}
		matches = append(matches, id)
	}
	return matches, nil
}

func (s *bigQueryStore) Add(ctx context.Context, teamID string, matches Matches) error {
	log := s.l.With(zap.String("request.id", middleware.GetReqID(ctx)), zap.String("team.id", teamID))
	opTimer := time.Now()
	defer logDuration(log, "matches add complete", "Add.duration", opTimer)

	// convert matches into a BigQuery compatible json document (new-line demited)
	timer := time.Now()
	var err error
	var matchesBytes, timelinesBytes bytes.Buffer
	var matchesEnc, timelinesEnc = json.NewEncoder(&matchesBytes), json.NewEncoder(&timelinesBytes)
	for _, m := range matches {
		// encode match
		if err = matchesEnc.Encode(&m.Details); err != nil {
			log.Error("failed to unmarshal a match",
				zap.Error(err),
				zap.Any("match", m))
			continue
		}

		// encode timeline
		m.Timeline.GameID = m.Details.GameID
		if err = timelinesEnc.Encode(&m.Timeline); err != nil {
			log.Error("failed to unmarshal a timeline",
				zap.Error(err),
				zap.Any("match", m))
			continue
		}
	}
	log.Info("marshal complete", zap.Duration("duration", time.Since(timer)))

	// upload data
	for _, set := range []struct {
		key  string
		data io.Reader
	}{
		{s.cfg.MatchesTableID, &matchesBytes},
		{s.cfg.TimelinesTableID, &timelinesBytes},
	} {
		// create BiqQuery data source for upload
		timer = time.Now()
		data := bigquery.NewReaderSource(set.data)
		data.SourceFormat = bigquery.JSON
		data.AutoDetect = true

		// run upload
		// TODO: there are somewhat strict limitations on this: https://cloud.google.com/bigquery/quotas#load_jobs
		// consider pooling multiple team's new matches together
		loader := s.bqDataset().Table(set.key).LoaderFrom(data)
		loader.Labels = map[string]string{"team": teamID}
		stats, err := s.execLoader(ctx, loader)
		if err != nil {
			return fmt.Errorf("%s: %v", set.key, err)
		}
		log.Info("upload complete",
			zap.Any("stats", stats),
			zap.Duration(fmt.Sprintf("%s.duration", set.key), time.Since(timer)))
	}

	return nil
}

func (s *bigQueryStore) Close() error { return s.Close() }

func (s *bigQueryStore) bqDataset() *bigquery.Dataset { return s.bq.Dataset(s.cfg.DatasetID) }

func (s *bigQueryStore) bqTeamView(teamID string) *bigquery.Table {
	return s.bqDataset().Table(teamView(teamID))
}

func (s *bigQueryStore) execQuery(
	ctx context.Context,
	query *bigquery.Query,
) (*bigquery.Job, *bigquery.QueryStatistics, error) {
	job, err := query.Run(ctx)
	if err != nil {
		return nil, nil, err
	}
	status, err := job.Wait(ctx)
	if err != nil {
		return nil, nil, err
	}
	if err := status.Err(); err != nil {
		return nil, nil, err
	}
	if status.Statistics == nil {
		return nil, nil, errors.New("no statistics attached to query")
	}
	queryStats, ok := status.Statistics.Details.(*bigquery.QueryStatistics)
	if !ok {
		return nil, nil, fmt.Errorf("did not receive query stats, got %T", status.Statistics.Details)
	}

	return job, queryStats, nil
}

func (s *bigQueryStore) execLoader(
	ctx context.Context,
	loader *bigquery.Loader,
) (*bigquery.LoadStatistics, error) {
	job, err := loader.Run(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not run data load: %v", err)
	}
	status, err := job.Wait(ctx)
	if err != nil {
		return nil, fmt.Errorf("data load job failed to complete: %v", err)
	}
	if err := status.Err(); err != nil {
		return nil, fmt.Errorf("data load job completed with error: %v, errors: %+v", err, status.Errors)
	}
	if status.Statistics == nil {
		return nil, errors.New("no statistics attached to loader")
	}
	loadStats, ok := status.Statistics.Details.(*bigquery.LoadStatistics)
	if !ok {
		return nil, fmt.Errorf("did not receive query stats, got %T", status.Statistics.Details)
	}
	return loadStats, nil
}

func teamView(teamID string) string { return fmt.Sprintf("team_%s", teamID) }
