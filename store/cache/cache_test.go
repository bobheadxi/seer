package cache

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/dgraph-io/badger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"

	"go.bobheadxi.dev/seer/store"
	"go.bobheadxi.dev/seer/store/mocks"
	"go.bobheadxi.dev/zapx/zapx"
)

func TestCacheGetTeam(t *testing.T) {
	l := zaptest.NewLogger(t)
	s := mocks.FakeStore{}
	tmpDir := os.TempDir()
	defer os.RemoveAll(tmpDir)

	db, err := badger.Open(badger.DefaultOptions(tmpDir).WithLogger(zapx.NewFormatLogger(l)))
	require.NoError(t, err)

	c := New(l, &s, db)
	defer c.Close()

	ctx := context.Background()
	exampleTeam := "1234"
	testTeamData := &store.TeamWithAnalytics{
		Team: &store.Team{Region: "some_region"},
	}
	s.GetTeamReturnsOnCall(0, testTeamData, nil)
	got, err := c.GetTeam(ctx, exampleTeam)
	assert.NoError(t, err)
	assert.EqualValues(t, testTeamData, got)

	s.GetTeamReturnsOnCall(1, nil, errors.New("should have hit cache"))
	got, err = c.GetTeam(ctx, exampleTeam)
	assert.NoError(t, err)
	assert.EqualValues(t, testTeamData, got)
}
