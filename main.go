package main

import (
	"fmt"
	"os"
	"strings"

	"go.uber.org/zap"

	"go.bobheadxi.dev/seer/config"
	"go.bobheadxi.dev/zapx/zapx"
)

func main() {
	// load up configuration
	flags, err := config.LoadFlags(os.Args[1:])
	if err != nil {
		if strings.Contains(err.Error(), "help requested") {
			return
		}
		panic(err)
	}
	cfg, err := config.NewEnvConfig()
	if err != nil {
		panic(err)
	}
	meta := config.NewBuildMeta()

	// init logger
	log, err := zapx.New(flags.LogPath, flags.Dev,
		zapx.WithDebug(flags.Dev),
		zapx.WithFields(map[string]interface{}{
			"build.version": meta.Version,
		}))
	if err != nil {
		panic(err)
	}
	defer log.Sync()

	// init dynamic configuration
	if !cfg.NoDynamicConfiguration {
		log.Info("initializing dynamic configuration")
		cfg.InitDynamicConfig(log.Named("dynamic_config"))
	} else {
		log.Info("dynamic configuration disabled")
	}

	// report basic config
	log.Debug("configuration loaded",
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

	// handle other operation modes
	switch flags.Mode() {
	case config.ModeServer:
		startServer(log.Named("server"), flags, *cfg, meta)

	case config.ModeWorker:
		startWorker(log.Named("worker"), flags, *cfg, meta)

	case config.ModeAll:
		go startServer(log.Named("server"), flags, *cfg, meta)
		startWorker(log.Named("worker"), flags, *cfg, meta)

	default:
		fmt.Printf("unknown mode '%s'", flags.Mode())
		os.Exit(1)
	}
}
