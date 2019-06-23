package jobs

import (
	"context"
	"strconv"
	"time"

	"github.com/gocraft/work"
	"go.bobheadxi.dev/seer/riot"
	"go.bobheadxi.dev/seer/store"
	"go.uber.org/zap"
)

const jobMatchesSync = "matches_sync"

type matchesSyncJob struct {
	teamID    string
	requestID string
}

// NewMatchesSyncJob instantiates a runnable job for syncing a team's latest matches
func NewMatchesSyncJob(teamID, requestID string) Job {
	return &matchesSyncJob{teamID, requestID}
}

func (m *matchesSyncJob) Name() string { return jobMatchesSync }

func (m *matchesSyncJob) Params() map[string]interface{} {
	return map[string]interface{}{
		"team.id":    m.teamID,
		"request.id": m.requestID,
	}
}

type matchesSyncContext struct {
	l *zap.Logger
	*BaseJobContext
}

func (m *matchesSyncContext) Run(job *work.Job) error {
	var (
		ctx       = context.Background()
		teamID    = job.ArgString("team.id")
		requestID = job.ArgString("request.id")
		start     = time.Now()
		log       = m.l.With(
			zap.String("job.id", job.ID),
			zap.String("team.id", teamID),
			zap.String("request.id", requestID))
	)

	log.Info("job started")
	defer log.Info("job completed", zap.Duration("job.duration", time.Since(start)))

	// get known team data
	log.Debug("looking for known team data")
	team, storedMatches, err := m.Store.Get(ctx, teamID)
	if err != nil {
		return err
	}

	// see what matches have already been tracked
	knownMatches := make(map[int64]bool)
	for _, match := range storedMatches {
		knownMatches[match.Details.GameID] = true
	}

	// check for shared matches
	log.Debug("check match history of each member",
		zap.String("riot.region", string(team.Region)))
	discoveredMatches := make(map[int64]int)
	api := m.RiotAPI.WithRegion(riot.Region(team.Region))
	for _, member := range team.Members {
		matches, err := api.Matches(ctx, member.AccountID)
		if err != nil {
			return err
		}
		for _, match := range matches {
			if knownMatches[match.GameID] {
				continue
			}
			discoveredMatches[match.GameID]++
		}
	}

	// look for match details
	log.Debug("querying for match details")
	var matchesToStore []store.MatchData
	for game, count := range discoveredMatches {
		if count < 4 {
			continue
		}
		details, err := api.MatchDetails(ctx, strconv.Itoa(int(game)))
		if err != nil {
			return err
		}
		matchesToStore = append(matchesToStore, store.MatchData{
			Details: details,
		})
	}

	if len(matchesToStore) == 0 {
		log.Info("no matches found for team")
		// TODO: maybe try paging back further?
		return nil
	}

	log.Debug("storing match data")
	if err := m.Store.Add(ctx, teamID, matchesToStore); err != nil {
		return err
	}

	return nil
}
