package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
	"go.uber.org/multierr"
	"go.uber.org/zap"

	"go.bobheadxi.dev/seer/store"
)

// cachedStore is a cached overlay on top of a store.Store
// TODO: forgot I have redis, maybe this should use redis instead
type cachedStore struct {
	l *zap.Logger

	store store.Store
	pool  *redis.Pool
}

// New instantiates a new cached store
func New(l *zap.Logger, s store.Store, redisPool *redis.Pool) (store.Store, error) {
	// test pool connection
	conn, err := redisPool.Dial()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}
	conn.Close()

	return &cachedStore{l: l, store: s, pool: redisPool}, nil
}

func (c *cachedStore) Create(ctx context.Context, teamID string, team *store.Team) error {
	return c.store.Create(ctx, teamID, team)
}

func (c *cachedStore) Add(ctx context.Context, teamID string, matches store.Matches) error {
	return c.store.Add(ctx, teamID, matches)
}

// GetTeam is an expensive operation
func (c *cachedStore) GetTeam(ctx context.Context, teamID string) (*store.TeamWithAnalytics, error) {
	key := "/team/" + teamID
	var team *store.TeamWithAnalytics
	found, err := c.get(key, &team)
	if err != nil {
		c.l.Error("cache check failed", zap.Error(err))
	}
	if found {
		return team, nil
	}

	team, err = c.store.GetTeam(ctx, teamID)
	if err != nil {
		return nil, err
	}
	return team, c.put(key, team, 12*time.Hour)
}

// GetMatches needs to be live to facilitate sync, and is inexpensive
func (c *cachedStore) GetMatches(ctx context.Context, teamID string) ([]int64, error) {
	return c.store.GetMatches(ctx, teamID)
}

func (c *cachedStore) Close() error {
	var errors error
	if err := c.store.Close(); err != nil {
		multierr.Append(errors, fmt.Errorf("cache: underlying store failed to close: %w", err))
	}
	return errors
}

func (c *cachedStore) get(k string, out interface{}) (bool, error) {
	conn := c.pool.Get()
	defer conn.Close()

	resp, err := redis.String(conn.Do("GET", k))
	if err != nil {
		if errors.Is(err, redis.ErrNil) {
			return false, nil
		}
		return false, fmt.Errorf("failed to retrieve value '%s' from cache: %w", k, err)
	}
	if err := json.Unmarshal([]byte(resp), &out); err != nil {
		return false, fmt.Errorf("failed to read discovered value from cache: %w", err)
	}

	return true, nil
}

func (c *cachedStore) put(k string, v interface{}, ttl time.Duration) error {
	bytes, err := json.Marshal(&v)
	if err != nil {
		return fmt.Errorf("failed to marshal data for '%s': %w", k, err)
	}

	conn := c.pool.Get()
	defer conn.Close()

	conn.Send("SET", k, string(bytes))
	conn.Send("EXPIRE", k, ttl.Seconds())
	if err := conn.Flush(); err != nil {
		return fmt.Errorf("failed to write key '%s' to cache: %w", k, err)
	}

	c.l.Info("placed value into cache",
		zap.String("key", k),
		zap.Int("size", len(bytes)),
		zap.Duration("ttl", ttl))

	return nil
}
