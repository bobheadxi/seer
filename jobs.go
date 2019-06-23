package main

import (
	"github.com/gocraft/work/webui"
	"github.com/gomodule/redigo/redis"
)

func startJobsUI(namespace string, pool *redis.Pool, addr string) {
	webui.NewServer(namespace, pool, addr).Start()
}
