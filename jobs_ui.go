package main

import (
	"fmt"

	"github.com/gocraft/work/webui"
	"github.com/gomodule/redigo/redis"
	"go.uber.org/zap"
)

func startJobsUI(log *zap.Logger, namespace string, pool *redis.Pool, addr string) {
	log.Info(fmt.Sprintf("starting jobs UI on '%s' for namespace '%s'", addr, namespace))
	webui.NewServer(namespace, pool, addr).Start()
	<-newStopper()
}
