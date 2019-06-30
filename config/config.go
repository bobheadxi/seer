package config

import (
	"os"
	"time"

	"github.com/gomodule/redigo/redis"
	"golang.org/x/oauth2"
)

// Config exposes server configuration
type Config struct {
	RedisNamespace string
	RedisURL       string // redis conn string
	RedisAddr      string

	RiotAPIToken string // TODO: need a mechanism to update this

	GitHubToken     string
	GitHubStoreRepo GitHubStoreRepo
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

func firstOf(vars ...string) string {
	for _, s := range vars {
		if s != "" {
			return s
		}
	}
	return ""
}
