package jobs

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/go-chi/chi/middleware"
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

func (m *matchesSyncJob) Unique() bool { return true }

func (m *matchesSyncJob) Params() map[string]interface{} {
	return map[string]interface{}{
		"team.id":    m.teamID,
		"request.id": m.requestID,
	}
}

type matchesSyncContext struct {
	l *zap.Logger
	q Queue
	*BaseJobContext
}

func (m *matchesSyncContext) Run(job *work.Job) error {
	var (
		teamID    = job.ArgString("team.id")
		requestID = job.ArgString("request.id")
		start     = time.Now()
		log       = m.l.With(
			zap.String("job.id", job.ID),
			zap.String("team.id", teamID),
			zap.String("request.id", requestID))
		ctx = context.WithValue(context.Background(), middleware.RequestIDKey, requestID)
	)

	log.Info("job started")
	defer log.Info("job completed", zap.Duration("job.duration", time.Since(start)))

	// get known team data
	log.Debug("looking for known team data")
	team, err := m.Store.GetTeam(ctx, teamID)
	if err != nil {
		log.Error("failed to find team", zap.Error(err))
		return fmt.Errorf("error while looking for team in store: %v", err)
	}

	storedMatches, err := m.Store.GetMatches(ctx, teamID)
	if err != nil {
		log.Error("failed to find matches", zap.Error(err))
		return fmt.Errorf("error while looking for matches in store: %v", err)
	}

	// see what matches have already been tracked
	knownMatches := make(map[int64]bool)
	for _, match := range storedMatches {
		knownMatches[match] = true
	}

	// check for shared matches
	log.Debug("check match history of each member",
		zap.String("riot.region", string(team.Region)))
	discoveredMatches := make(map[int64]int)
	api := m.RiotAPI.WithRegion(riot.Region(team.Region))
	for i, member := range team.Members {
		job.Checkin(fmt.Sprintf("member=%d name=%s", i, member.Name))
		matches, err := api.Matches(ctx, member.AccountID)
		if err != nil {
			log.Error("failed to find matches", zap.Error(err),
				zap.String("summoner", member.Name))
			return fmt.Errorf("error querying for matches for matches: %v", err)
		}
		for _, match := range matches {
			// ignore if we already have this match in store
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
		job.Checkin(fmt.Sprintf("game_count=%d game=%d", len(matchesToStore), game))
		details, err := api.MatchDetails(ctx, strconv.Itoa(int(game)))
		if err != nil {
			log.Error("failed to find game", zap.Error(err), zap.Int64("game", game))
			return fmt.Errorf("error querying for match details: %v", err)
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
	job.Checkin(fmt.Sprintf("new_game_count=%d", len(matchesToStore)))
	if err := m.Store.Add(ctx, teamID, matchesToStore); err != nil {
		log.Error("failed to store matches", zap.Error(err))
		return fmt.Errorf("error saving results to store: %v", err)
	}

	return nil
}
