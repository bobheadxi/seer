package config

import (
	"os"

	"golang.org/x/oauth2"
)

// Config exposes server configuration
type Config struct {
	RiotAPIToken string
	GitHubToken  string
}

// NewEnvConfig instatiates configuration from environment
func NewEnvConfig() Config {
	return Config{
		RiotAPIToken: os.Getenv("RIOT_API_TOKEN"),
		GitHubToken:  os.Getenv("GITHUB_TOKEN"),
	}
}

// RiotAPITokenSource inits a static token source from this configuration
func (c *Config) RiotAPITokenSource() oauth2.TokenSource {
	return oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: c.RiotAPIToken,
	})
}

// GitHubAPITokenSource inits a static token source from this configuration
func (c *Config) GitHubAPITokenSource() oauth2.TokenSource {
	return oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: c.RiotAPIToken,
	})
}
