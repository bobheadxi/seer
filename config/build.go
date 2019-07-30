package config

import (
	"fmt"
	"os"
	"time"
)

var (
	// Version is the version or commit hash of this build - can be injected at build time
	Version string
)

// BuildMeta denotes build metadata
type BuildMeta struct {
	Version string
}

// NewBuildMeta instantiates a new build metadata struct from the environment.
// Currently leverages Heroku's Dyno Metadata: https://devcenter.heroku.com/articles/dyno-metadata
func NewBuildMeta() BuildMeta {
	v := firstOf(Version, os.Getenv("HEROKU_SLUG_COMMIT"), os.Getenv("SEER_VERSION"))
	if len(v) > 30 {
		v = v[:7] // trim commits
	} else if v == "" {
		v = fmt.Sprintf("unknown-%s", time.Now().Format(time.RFC3339))
	}
	return BuildMeta{
		Version: v,
	}
}
