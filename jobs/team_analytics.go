package jobs

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gocraft/work"
	"go.bobheadxi.dev/seer/store"
	"go.uber.org/zap"
)

const jobTeamAnalytics = "team_analytics"

type teamAnalyticsJob struct {
	teamID    string
	requestID string
	syncID    string

	team       []byte
	newMatches []byte
}

// NewTeamAnalyticsJob instantiates a runnable job for syncing a team's latest matches
func NewTeamAnalyticsJob(
	teamID, requestID, syncID string,
	team *store.Team,
	newMatches []store.MatchData,
) (Job, error) {
	var err error
	var newMatchesData, teamData []byte
	if len(newMatches) > 0 {
		newMatchesData, err = json.Marshal(newMatches)
		if err != nil {
			return nil, fmt.Errorf("could not marshal matches: %v", err)
		}
	}

	if team != nil {
		teamData, err = json.Marshal(team)
		if err != nil {
			return nil, fmt.Errorf("could not marshal team: %v", err)
		}
	}

	return &teamAnalyticsJob{
		teamID, requestID, syncID,
		teamData,
		newMatchesData,
	}, nil
}

func (t *teamAnalyticsJob) Name() string { return jobTeamAnalytics }

func (t *teamAnalyticsJob) Unique() bool { return true }

func (t *teamAnalyticsJob) Params() map[string]interface{} {
	return map[string]interface{}{
		"sync.id":          t.syncID,
		"team.id":          t.teamID,
		"request.id":       t.requestID,
		"team_data":        string(t.team),
		"new_matches_data": string(t.newMatches),
	}
}

type teamAnalyticsContext struct {
	l *zap.Logger
	*BaseJobContext
}

func (t *teamAnalyticsContext) Run(job *work.Job) error {
	var (
		teamID         = job.ArgString("team.id")
		requestID      = job.ArgString("request.id")
		teamData       = job.ArgString("team_data")
		newMatchesData = job.ArgString("new_matches_data")

		start = time.Now()
		log   = t.l.With(
			zap.String("job.id", job.ID),
			zap.String("team.id", teamID),
			zap.String("request.id", requestID))
		//ctx = context.WithValue(context.Background(), middleware.RequestIDKey, requestID)
	)

	log.Info("job started")
	defer log.Info("job completed", zap.Duration("job.duration", time.Since(start)))

	var team store.Team
	if err := json.Unmarshal([]byte(teamData), &team); err != nil {
		log.Error("unable to unmarshal team from parameter",
			zap.Error(err),
			zap.String("raw_data", teamData))
		return fmt.Errorf("unable to read team: %v", err)
	}

	var newMatches []store.MatchData
	if err := json.Unmarshal([]byte(newMatchesData), &newMatches); err != nil {
		log.Error("unable to unmarshal new matches from parameter",
			zap.Error(err),
			zap.String("raw_data", newMatchesData))
		return fmt.Errorf("unabled to read new matches: %v", err)
	}

	// TODO

	return nil
}
