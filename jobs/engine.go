package jobs

import (
	"fmt"
	"runtime"
	"time"

	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
	"go.uber.org/multierr"
	"go.uber.org/zap"
)

// Engine is the contract for managing job execution
type Engine interface {
	Start()
	Close()
}

type engine struct {
	l    *zap.Logger
	pool *work.WorkerPool
	jobs *work.Client
}

// NewJobsEngine instantiates a new job runn
func NewJobsEngine(l *zap.Logger, app string, redisPool *redis.Pool, b *BaseJobContext) Engine {
	pool := work.NewWorkerPool(BaseJobContext{}, uint(runtime.NumCPU()), app, redisPool)
	queue := NewJobQueue(l.Named("queue"), app, redisPool)

	matchesSync := &matchesSyncContext{l.Named("matches_sync"), queue, b}
	pool.JobWithOptions(jobMatchesSync, work.JobOptions{
		Priority: 10,
		MaxFails: 3,
		// TODO: create backoff calculator that checks for rate limits
		// Backoff:
	}, matchesSync.Run)

	teamAnalytics := &teamAnalyticsContext{l.Named("team_analytics"), b}
	pool.JobWithOptions(jobTeamAnalytics, work.JobOptions{
		Priority: 5,
		MaxFails: 3,
		// TODO: create backoff calculator that checks for rate limits
		// Backoff:
	}, teamAnalytics.Run)

	return &engine{
		l:    l,
		pool: pool,
		jobs: work.NewClient(app, redisPool),
	}
}

func (e *engine) Start() {
	e.l.Info("starting job engine")
	e.pool.Start()

	// spin up status reporter
	var check int
	interval := time.NewTicker(5 * time.Minute) // TODO: configure?
	for range interval.C {
		e.reportState(check)
		check++
	}
}

func (e *engine) Close() {
	e.l.Info("draining jobs and closing pool")
	t := time.AfterFunc(30*time.Second, e.pool.Stop)
	e.pool.Drain()
	t.Stop()
}

func (e *engine) reportState(id int) {
	var errors error
	heartbeats, err := e.jobs.WorkerPoolHeartbeats()
	if err != nil {
		multierr.Append(errors, fmt.Errorf("could not get hearbeats: %v", err))
	}

	// TODO: report more details if desired?
	_, pagesDead, err := e.jobs.DeadJobs(0)
	if err != nil {
		multierr.Append(errors, fmt.Errorf("could not get dead jobs: %v", err))
	}
	_, pagesScheduled, err := e.jobs.ScheduledJobs(0)
	if err != nil {
		multierr.Append(errors, fmt.Errorf("could not get scheduled jobs: %v", err))
	}
	_, pagesRetry, err := e.jobs.RetryJobs(0)
	if err != nil {
		multierr.Append(errors, fmt.Errorf("could not get retry jobs: %v", err))
	}

	// report state
	e.l.Info("engine state report",
		zap.Int("jobstatereport.id", id),
		zap.Any("heartbeats", heartbeats),
		zap.Int64("dead_jobs", pagesDead*20),
		zap.Int64("scheduled_jobs", pagesScheduled*20),
		zap.Int64("retry_jobs", pagesRetry))
	if errors != nil {
		e.l.Error("error(s) encountered on engine state report",
			zap.Int("jobstatereport.id", id),
			zap.Error(errors))
	}
}
