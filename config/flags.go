package config

import (
	"flag"
	"os"
)

// Flags contains command-line flag configuration
type Flags struct {
	Dev        bool
	LogPath    string
	JobsUIPort string
	APIPort    string
	Mode       string
}

// LoadFlags loads flags from the given set of arguments
func LoadFlags(args []string) (Flags, error) {
	var v Flags
	flags := flag.NewFlagSet("seer", flag.ContinueOnError)
	flags.BoolVar(&v.Dev, "dev", os.Getenv("DEV") == "true", "toggle dev mode")
	flags.StringVar(&v.LogPath, "logpath", "", "path for log storage")
	flags.StringVar(&v.JobsUIPort, "jobs-ui", "", "enable jobs UI on given port")
	flags.StringVar(&v.APIPort, "port", "8080", "port to serve Seer API on")
	flags.StringVar(&v.Mode, "mode", "server", "operation mode to run in")
	return v, flags.Parse(args)
}
