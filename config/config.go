package config

import (
	"os"
	"time"

	"github.com/gomodule/redigo/redis"
	"golang.org/x/oauth2"
	"google.golang.org/api/option"
)

// Config exposes server configuration
type Config struct {
	// Redis
	RedisNamespace string
	RedisURL       string // redis conn string
	RedisAddr      string

	// Riot API
	RiotAPIToken string // TODO: need a mechanism to update this

	// GitHub
	GitHubToken     string
	GitHubStoreRepo GitHubStoreRepo

	// Google Cloud
	GCPProjectID   string
	GCPCredentials string
	BigQuery       BigQuery
}

// NewEnvConfig instatiates configuration from environment
func NewEnvConfig() Config {
	return Config{
		RedisNamespace: os.Getenv("REDIS_NAMESPACE"),
		RedisURL:       os.Getenv("REDIS_URL"),
		RedisAddr:      os.Getenv("REDIS_ADDR"),

		RiotAPIToken: os.Getenv("RIOT_API_TOKEN"),
		GitHubToken:  os.Getenv("GITHUB_TOKEN"),

		GitHubStoreRepo: GitHubStoreRepo{
			Owner: os.Getenv("GITHUB_STORE_OWNER"),
			Repo:  os.Getenv("GITHUB_STORE_REPO"),
		},

		GCPProjectID:   os.Getenv("GCP_PROJECT_ID"),
		GCPCredentials: os.Getenv("GCP_CREDENTIALS"),
		BigQuery: BigQuery{
			DatasetID:      os.Getenv("BIGQUERY_DATASET_ID"),
			MatchesTableID: os.Getenv("BIGQUERY_TABLE_ID_MATCHES"),
		},
	}
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
