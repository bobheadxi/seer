package config

import (
	"flag"
	"os"
)

// Mode denotes operational mode for a run
type Mode string

const (
	// ModeAll starts the server, worker, and if enabled, the jobs UI
	ModeAll Mode = "all"
	// ModeServer starts the server and if enabled, the jobs UI
	ModeServer Mode = "server"
	// ModeWorker starts the worker and if enabled, the jobs UI
	ModeWorker Mode = "worker"
	// ModeJobsUIOnly starts just the jobs UI
	ModeJobsUIOnly Mode = "jobs-ui-only"
)

// Flags contains command-line flag configuration
type Flags struct {
	Dev        bool
	LogPath    string
	JobsUIPort string
	APIPort    string

	mode string
}

// Mode returns the configured operational mode
func (f Flags) Mode() Mode { return Mode(f.mode) }

// LoadFlags loads flags from the given set of arguments
func LoadFlags(args []string) (Flags, error) {
	var v Flags
	flags := flag.NewFlagSet("seer", flag.ContinueOnError)
	flags.BoolVar(&v.Dev, "dev", os.Getenv("DEV") == "true", "toggle dev mode")
	flags.StringVar(&v.LogPath, "logpath", "", "path for log storage")
	flags.StringVar(&v.JobsUIPort, "jobs-ui", "", "enable jobs UI on given port")
	flags.StringVar(&v.APIPort, "port", "8080", "port to serve Seer API on")

	flags.StringVar(&v.mode, "mode", string(ModeAll), "operation mode to run in")
	return v, flags.Parse(args)
}
