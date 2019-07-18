package jobs

import (
	"runtime"
	"time"

	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
	"go.uber.org/zap"

	"go.bobheadxi.dev/seer/riot"
	"go.bobheadxi.dev/seer/store"
)

// Engine is the contract for managing job execution
type Engine interface {
	Start()
	Close()
}

type engine struct {
	l    *zap.Logger
	pool *work.WorkerPool
}

// BaseJobContext denotes dependencies common to most jobs
type BaseJobContext struct {
	RiotAPI riot.API
	Store   store.Store
}

// NewJobsEngine instantiates a new job runn
func NewJobsEngine(l *zap.Logger, app string, redisPool *redis.Pool, b *BaseJobContext) Engine {
	pool := work.NewWorkerPool(BaseJobContext{}, uint(runtime.NumCPU()), app, redisPool)

	matchesSync := &matchesSyncContext{l.Named("matches_sync"), b}
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
	}
}

func (e *engine) Start() {
	e.l.Info("starting job engine")
	e.pool.Start()
}

func (e *engine) Close() {
	e.l.Info("draining jobs and closing pool")
	t := time.AfterFunc(30*time.Second, e.pool.Stop)
	e.pool.Drain()
	t.Stop()
}
