package config

import (
	"flag"
	"os"
)

// Flags contains command-line flag configuration
type Flags struct {
	Dev           bool
	LogPath       string
	DisableJobsUI bool
	APIPort       string
}

// LoadFlags loads flags from the given set of arguments
func LoadFlags(args []string) (Flags, error) {
	var v Flags
	flags := flag.NewFlagSet("seer", flag.ContinueOnError)
	flags.BoolVar(&v.Dev, "dev", os.Getenv("DEV") == "true", "toggle dev mode")
	flags.StringVar(&v.LogPath, "logpath", "", "path for log storage")
	flags.BoolVar(&v.DisableJobsUI, "no-jobs-ui", false, "disable jobs UI")
	flags.StringVar(&v.APIPort, "port", "8080", "port to serve Seer API on")
	return v, flags.Parse(args)
}
