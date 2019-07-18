package main

import (
	"fmt"
	"os"

	"go.uber.org/zap"

	"go.bobheadxi.dev/seer/config"
	"go.bobheadxi.dev/zapx"
)

func main() {
	// load up configuration
	flags, err := config.LoadFlags(os.Args[1:])
	if err != nil {
		panic(err)
	}
	cfg := config.NewEnvConfig()
	meta := config.NewBuildMeta()

	// init logger
	log, err := zapx.New(flags.LogPath, flags.Dev,
		zapx.WithDebug(flags.Dev),
		zapx.WithFields(map[string]interface{}{
			"build.commit": meta.Commit,
		}))
	if err != nil {
		panic(err)
	}
	defer log.Sync()

	// report basic config
	log.Info("configuration loaded",
		zap.Any("config", cfg),
		zap.Any("meta", meta),
		zap.Any("flags", flags))

	// handle jobs-ui-only mode first
	if flags.Mode() == config.ModeJobsUIOnly {
		if flags.JobsUIPort == "" {
			flags.JobsUIPort = "8081"
		}
		startJobsUI(log.Named("jobs_ui"), cfg.RedisNamespace, cfg.DefaultRedisPool(), ":"+flags.JobsUIPort)
		return
	}

	// spin up jobs manager UI as well if configured to do so for all other modes
	if flags.JobsUIPort != "" {
		go startJobsUI(
			log.Named("jobs_ui"),
			cfg.RedisNamespace,
			cfg.DefaultRedisPool(),
			":"+flags.JobsUIPort)
	}
	<-newStopper()

	// handle other operation modes
	switch flags.Mode() {
	case config.ModeServer:
		startServer(log.Named("server"), flags, cfg, meta)

	case config.ModeWorker:
		startWorker(log.Named("worker"), flags, cfg, meta)

	case config.ModeAll:
		go startServer(log.Named("server"), flags, cfg, meta)
		startWorker(log.Named("worker"), flags, cfg, meta)

	default:
		fmt.Printf("unknown mode '%s'", flags.Mode())
		os.Exit(1)
	}
}
