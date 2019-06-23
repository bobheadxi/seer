package store

import (
	"context"

	"go.bobheadxi.dev/seer/riot"
)

// Team represents a team
type Team struct {
	Region  riot.Region
	Members []*riot.Summoner
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
