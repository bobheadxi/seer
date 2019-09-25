package cache

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"

	"go.bobheadxi.dev/seer/config"
	"go.bobheadxi.dev/seer/store"
	"go.bobheadxi.dev/seer/store/mocks"
)

func TestCacheGetTeam(t *testing.T) {
	godotenv.Load("../../.env")
	if os.Getenv("CACHE_INTEGRATION") != "true" {
		t.Skip("CACHE_INTEGRATION not set to true")
	}

	l := zaptest.NewLogger(t)
	s := mocks.FakeStore{}
	cfg, err := config.NewEnvConfig()
	require.NoError(t, err)

	c, err := New(l, &s, cfg.DefaultRedisPool())
	require.NoError(t, err)
	defer c.Close()

	// first get should place team into cache
	ctx := context.Background()
	exampleTeam := "1234"
	testTeamData := &store.TeamWithAnalytics{
		Team: &store.Team{Region: "some_region"},
	}
	s.GetTeamReturnsOnCall(0, testTeamData, nil)
	got, err := c.GetTeam(ctx, exampleTeam)
	assert.NoError(t, err)
	assert.EqualValues(t, testTeamData, got)

	// second get should hit cache
	s.GetTeamReturnsOnCall(1, nil, errors.New("should have hit cache"))
	got, err = c.GetTeam(ctx, exampleTeam)
	assert.NoError(t, err)
	assert.EqualValues(t, testTeamData, got)
}
