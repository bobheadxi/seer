package config

import "os"

var (
	// Commit is the commit hash of this build - can be injected at build time
	Commit string
)

// BuildMeta denotes build metadata
type BuildMeta struct {
	Commit string
}

// NewBuildMeta instantiates a new build metadata struct from the environment.
// Currently leverages Heroku's Dyno Metadata: https://devcenter.heroku.com/articles/dyno-metadata
func NewBuildMeta() BuildMeta {
	c := firstOf(Commit, os.Getenv("HEROKU_SLUG_COMMIT"))
	if len(c) > 8 {
		c = c[:7]
	}
	return BuildMeta{
		Commit: c,
	}
}
