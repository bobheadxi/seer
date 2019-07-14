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

	// handle any special operational mode
	switch flags.Mode {
	case "jobs-ui":
		if flags.JobsUIPort == "" {
			flags.JobsUIPort = "8081"
		}
		startJobsUI(
			log.Named("jobs_ui"),
			cfg.RedisNamespace,
			cfg.DefaultRedisPool(),
			":"+flags.JobsUIPort)
	case "server":
		startServer(log, flags, cfg, meta)
	default:
		fmt.Printf("unknown mode '%s'", flags.Mode)
		os.Exit(1)
	}
}
