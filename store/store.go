package store

import (
	"context"
	"crypto/md5"
	"fmt"
	"sort"
	"strings"

	"go.bobheadxi.dev/seer/analytics"
	"go.bobheadxi.dev/seer/riot"
)

// Store defines the contract different storage backends for Seer should have
type Store interface {
	Create(ctx context.Context, teamID string, team *Team) error
	Add(ctx context.Context, teamID string, matches Matches) error

	GetTeam(ctx context.Context, teamID string) (*Team, error)
	GetMatches(ctx context.Context, teamID string) ([]int64, error)

	// LastUpdated(ctx context.Context, teamID string) (*time.Time, error)

	Close() error
}

// Team represents a team
type Team struct {
	Region    riot.Region          `json:"region"`
	Members   []*riot.Summoner     `json:"members"`
	Analytics *analytics.Analytics `json:"analytics"`
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
	Details  *riot.MatchDetails  `json:"details"`
	Timeline *riot.MatchTimeline `json:"timeline"`
}

// Matches is a collection of MatchData. It can be sorted so that the most recent
// games come first.
type Matches []MatchData

// Less reports whether the element with
// index i should sort before the element with index j.
func (m Matches) Len() int { return len(m) }

// Less reports whether the element with
// index i should sort before the element with index j.
func (m Matches) Less(i, j int) bool {
	return m[i].Details.GameCreation > m[i].Details.GameCreation
}

// Swap swaps the elements with indexes i and j.
func (m Matches) Swap(i, j int) {
	tmp := m[i]
	m[i] = m[j]
	m[j] = tmp
}
