package store

import (
	"time"

	"go.uber.org/zap"
)

func logDuration(log *zap.Logger, msg, key string, start time.Time) {
	log.Info(msg, zap.Duration(key, time.Since(start)))
}
