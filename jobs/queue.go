package jobs

import (
	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
	"go.uber.org/zap"
)

// Queue is an interface for accepting jobs in a queue
type Queue interface {
	Queue(j Job) (string, error)
}

type queue struct {
	l        *zap.Logger
	enquerer *work.Enqueuer
}

// NewJobQueue instantiates a new queue
func NewJobQueue(log *zap.Logger, app string, redisPool *redis.Pool) Queue {
	return &queue{
		l:        log,
		enquerer: work.NewEnqueuer(app, redisPool),
	}
}

// Queue adds a job to the job engine
func (q *queue) Queue(j Job) (string, error) {
	name, params := j.Name(), j.Params()
	log := q.l.With(zap.String("job.name", name), zap.Any("job.params", params))

	var job *work.Job
	var err error
	if j.Unique() {
		job, err = q.enquerer.EnqueueUnique(j.Name(), params)
	} else {
		job, err = q.enquerer.Enqueue(j.Name(), params)
	}
	if err != nil {
		log.Error("job failed to queue", zap.Error(err))
		return "", err
	}
	log.Info("job queued", zap.String("job.id", job.ID))
	return job.ID, nil
}
