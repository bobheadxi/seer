package config

import (
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/gomodule/redigo/redis"
	"golang.org/x/oauth2"
	"google.golang.org/api/option"
)

// Config exposes server configuration
type Config struct {
	// Redis
	RedisNamespace string `env:"REDIS_NAMESPACE,required"`
	RedisURL       string `env:"REDIS_URL"` // redis conn string
	RedisAddr      string `env:"REDIS_ADDR"`

	// Riot API
	RiotAPIToken string `env:"RIOT_API_TOKEN,required"` // TODO: need a mechanism to update this

	// GitHub
	GitHubToken     string `env:"GITHUB_TOKEN,required"`
	GitHubStoreRepo GitHubStoreRepo

	// Google Cloud
	GCPProjectID   string `env:"GCP_PROJECT_ID"`
	GCPCredentials string `env:"GCP_CREDENTIALS"`
	BigQuery       BigQuery
}

// NewEnvConfig instatiates configuration from environment
func NewEnvConfig() (Config, error) {
	var cfg Config
	return cfg, env.Parse(&cfg)
}

// GitHubAPITokenSource inits a static token source from this configuration
func (c *Config) GitHubAPITokenSource() oauth2.TokenSource {
	return oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: c.GitHubToken,
	})
}

// DefaultRedisPool inits a default redis configuration
func (c *Config) DefaultRedisPool(opts ...redis.DialOption) *redis.Pool {
	dialFunc := func() (redis.Conn, error) { return redis.DialURL(c.RedisURL, opts...) }
	if c.RedisAddr != "" {
		dialFunc = func() (redis.Conn, error) { return redis.Dial("tcp", c.RedisAddr, opts...) }
	}
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial:        dialFunc,
	}
}

// GCPConnOpts returns options needed to connect to GCP services
func (c *Config) GCPConnOpts() []option.ClientOption {
	var opts []option.ClientOption
	if c.GCPCredentials != "" {
		opts = []option.ClientOption{
			option.WithCredentialsJSON([]byte(c.GCPCredentials)),
		}
	}
	return opts
}

func firstOf(vars ...string) string {
	for _, s := range vars {
		if s != "" {
			return s
		}
	}
	return ""
}
