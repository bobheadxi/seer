package jobs

import (
	"runtime"
	"time"

	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
	"go.bobheadxi.dev/seer/riot"
	"go.bobheadxi.dev/seer/store"
	"go.uber.org/zap"
)

// Engine is the contract for managing job execution
type Engine interface {
	Queue(j Job) (string, error)

	Start()
	Close()
}

type engine struct {
	l     *zap.Logger
	pool  *work.WorkerPool
	queue *work.Enqueuer
}

// BaseJobContext denotes dependencies common to most jobs
type BaseJobContext struct {
	RiotAPI riot.API
	Store   store.Store
}

// NewJobsEngine instantiates a new job runn
func NewJobsEngine(l *zap.Logger, app string, redisPool *redis.Pool, b *BaseJobContext) Engine {
	pool := work.NewWorkerPool(nil, uint(runtime.NumCPU()), app, redisPool)

	matchesSync := &matchesSyncContext{l.Named("matches_sync"), b}
	pool.Job(jobMatchesSync, matchesSync.Run)

	return &engine{
		l:     l,
		pool:  pool,
		queue: work.NewEnqueuer(app, redisPool),
	}
}

// Queue adds a job to the job engine
func (e *engine) Queue(j Job) (string, error) {
	name, params := j.Name(), j.Params()
	log := e.l.With(zap.String("job.name", name), zap.Any("job.params", params))
	job, err := e.queue.Enqueue(j.Name(), params)
	if err != nil {
		log.Error("job failed to queue", zap.Error(err))
		return "", err
	}
	log.Info("job queued", zap.String("job.id", job.ID))
	return job.ID, nil
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
