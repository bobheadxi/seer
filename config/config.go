package config

import (
	"os"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/gomodule/redigo/redis"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"google.golang.org/api/option"
	configcat "gopkg.in/configcat/go-sdk.v1"
)

// Config exposes server configuration
type Config struct {
	// Redis
	RedisNamespace string `env:"REDIS_NAMESPACE,required"`
	RedisURL       string `env:"REDIS_URL"` // redis conn string
	RedisAddr      string `env:"REDIS_ADDR"`

	// GitHub
	GitHubToken     string `env:"GITHUB_TOKEN,required"`
	GitHubStoreRepo GitHubStoreRepo

	// Google Cloud
	GCPProjectID   string `env:"GCP_PROJECT_ID"`
	GCPCredentials string `env:"GCP_CREDENTIALS"`
	BigQuery       BigQuery

	// dynamic configuration
	ConfigCatKey string `env:"CONFIGCAT_KEY"`
	dynamic      *configcat.Client
}

// NewEnvConfig instatiates configuration from environment
func NewEnvConfig(l *zap.Logger) (*Config, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}
	catCfg := configcat.DefaultClientConfig()
	catCfg.Logger = &catLogger{l}
	cfg.dynamic = configcat.NewCustomClient(cfg.ConfigCatKey, catCfg)
	return &cfg, nil
}

const riotAPIKey = "RIOT_API_TOKEN"

// RiotAPIToken generates a token for the Riot API
func (c *Config) RiotAPIToken() string {
	return c.dynamic.GetValue(riotAPIKey, os.Getenv(riotAPIKey)).(string)
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
