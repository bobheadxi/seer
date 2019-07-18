package jobs

import (
	"time"

	"github.com/gocraft/work"
	"go.uber.org/zap"
)

const jobTeamAnalytics = "team_analytics"

type teamAnalyticsJob struct {
	teamID    string
	requestID string
}

// NewTeamAnalyticsJob instantiates a runnable job for syncing a team's latest matches
func NewTeamAnalyticsJob(teamID, requestID string) Job {
	return &matchesSyncJob{teamID, requestID}
}

func (t *teamAnalyticsJob) Name() string { return jobTeamAnalytics }

func (t *teamAnalyticsJob) Unique() bool { return true }

func (t *teamAnalyticsJob) Params() map[string]interface{} {
	return map[string]interface{}{
		"team.id":    t.teamID,
		"request.id": t.requestID,
	}
}

type teamAnalyticsContext struct {
	l *zap.Logger
	*BaseJobContext
}

func (t *teamAnalyticsContext) Run(job *work.Job) error {
	var (
		teamID    = job.ArgString("team.id")
		requestID = job.ArgString("request.id")
		start     = time.Now()
		log       = t.l.With(
			zap.String("job.id", job.ID),
			zap.String("team.id", teamID),
			zap.String("request.id", requestID))
		// ctx = context.WithValue(context.Background(), middleware.RequestIDKey, requestID)
	)

	log.Info("job started")
	defer log.Info("job completed", zap.Duration("job.duration", time.Since(start)))

	// TODO

	return nil
}
