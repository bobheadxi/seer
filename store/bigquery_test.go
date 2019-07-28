package store

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"

	"go.bobheadxi.dev/seer/config"
	"go.bobheadxi.dev/seer/riot"
	"go.bobheadxi.dev/seer/store/fixtures"
)

func TestBigQuery_Integration(t *testing.T) {
	godotenv.Load("../.env")
	if os.Getenv("BIGQUERY_INTEGRATION") != "true" {
		t.Skip("BIGQUERY_INTEGRATION not set to true")
	}

	var (
		ctx      = context.Background()
		l        = zaptest.NewLogger(t)
		testTeam = "test_team"
	)

	// initialize store
	cfg, err := config.NewEnvConfig()
	require.NoError(t, err)
	cfg.BigQuery.MatchesTableID = "matches_integration_test" // hard override for test
	bqs, err := NewBigQueryStore(ctx, l, BigQueryOpts{
		ProjectID: cfg.GCPProjectID,
		ConnOpts:  cfg.GCPConnOpts(),
		DataOpts:  cfg.BigQuery,
	})
	require.NoError(t, err)

	// delete tables
	matchesTable := bqs.(*bigQueryStore).bqMatchesTable()
	matchesTable.Delete(ctx)
	teamView := bqs.(*bigQueryStore).bqTeamView(testTeam)
	teamView.Delete(ctx)

	// read test data
	var match1, match2, match3, match4, match5 riot.MatchDetails
	require.NoError(t, json.Unmarshal([]byte(fixtures.TestMatch1), &match1))
	require.NoError(t, json.Unmarshal([]byte(fixtures.TestMatch2), &match2))
	require.NoError(t, json.Unmarshal([]byte(fixtures.TestMatch3), &match3))
	require.NoError(t, json.Unmarshal([]byte(fixtures.TestMatch4), &match4))
	require.NoError(t, json.Unmarshal([]byte(fixtures.TestMatch5), &match5))
	var team Team
	require.NoError(t, json.Unmarshal([]byte(fixtures.TestTeam), &team))

	// upload test data
	t.Run("Add()", func(t *testing.T) {
		require.NoError(t, bqs.Add(ctx, "integration-test",
			Matches{{&match1}, {&match2}, {&match3}, {&match4}, {&match5}}))
	})

	// create test team
	t.Run("Create()", func(t *testing.T) {
		require.NoError(t, bqs.Create(ctx, testTeam, &team))
	})

	// get matches
	t.Run("GetMatches()", func(t *testing.T) {
		matches, err := bqs.GetMatches(ctx, testTeam)
		require.NoError(t, err)
		assert.ElementsMatch(t, []int64{3072165694, 3059336276, 3057579582}, matches)
	})
}
