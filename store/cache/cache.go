package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/dgraph-io/badger"
	"go.uber.org/multierr"
	"go.uber.org/zap"

	"go.bobheadxi.dev/seer/store"
)

// cachedStore is a cached overlay on top of a store.Store
// TODO: forgot I have redis, maybe this should use redis instead
type cachedStore struct {
	l *zap.Logger

	store store.Store
	db    *badger.DB
}

// New instantiates a new cached store
func New(l *zap.Logger, s store.Store, db *badger.DB) store.Store {
	return &cachedStore{l: l, store: s, db: db}
}

func (c *cachedStore) Create(ctx context.Context, teamID string, team *store.Team) error {
	return c.store.Create(ctx, teamID, team)
}

func (c *cachedStore) Add(ctx context.Context, teamID string, matches store.Matches) error {
	return c.store.Add(ctx, teamID, matches)
}

// GetTeam is an expensive operation
func (c *cachedStore) GetTeam(ctx context.Context, teamID string) (*store.TeamWithAnalytics, error) {
	key := []byte("/team/" + teamID)
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
	if err := c.db.Close(); err != nil {
		multierr.Append(errors, fmt.Errorf("cache: underlying database failed to close: %w", err))
	}
	return errors
}

func (c *cachedStore) get(k []byte, out interface{}) (bool, error) {
	if err := c.db.View(func(tx *badger.Txn) error {
		item, err := tx.Get(k)
		if err != nil {
			return err
		}
		if item.IsDeletedOrExpired() {
			return nil
		}

		return item.Value(func(b []byte) error {
			if err := json.Unmarshal(b, &out); err != nil {
				return fmt.Errorf("failed to read from cache: %w", err)
			}
			return nil
		})
	}); err != nil {
		if errors.Is(err, badger.ErrKeyNotFound) {
			return false, nil
		}
		return false, fmt.Errorf("error occured when retrieving team from cache: %w", err)
	}

	return true, nil
}

func (c *cachedStore) put(k []byte, v interface{}, ttl time.Duration) error {
	bytes, err := json.Marshal(&v)
	if err != nil {
		return err
	}
	tx := c.db.NewTransaction(true)
	if err := tx.SetEntry(badger.NewEntry(k, bytes).WithTTL(ttl)); err != nil {
		c.l.Error("failed to write result to cache", zap.Error(err))
		return err
	}
	return tx.Commit()
}
