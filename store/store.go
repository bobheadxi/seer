package store

import (
	"context"
	"crypto/md5"
	"fmt"
	"sort"
	"strings"

	"go.bobheadxi.dev/seer/riot"
)

// Team represents a team
type Team struct {
	Region  riot.Region
	Members []*riot.Summoner
}

// GenerateTeamID generates a new team ID for this team based on its region and
// member player IDs
func (t *Team) GenerateTeamID() string {
	ids := make([]string, len(t.Members))
	for i, s := range t.Members {
		ids[i] = s.PlayerID
	}
	sort.Strings(ids)
	hash := md5.Sum([]byte(string(t.Region) + "|" + strings.Join(ids, "|")))
	return fmt.Sprintf("%x", hash)[:8]
}

// MatchData is a container for details about a match
type MatchData struct {
	Details *riot.MatchDetails
}

// Store defines the contract different storage backends for Seer should have
type Store interface {
	Create(ctx context.Context, teamID string, team *Team) error

	Get(ctx context.Context, teamID string) (*Team, []MatchData, error)
	Add(ctx context.Context, teamID string, matches []MatchData) error

	// LastUpdated(ctx context.Context, teamID string) (*time.Time, error)
}
