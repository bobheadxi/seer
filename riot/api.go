package riot

import (
	"context"
)

// API exposes the top-level client
type API interface {
	WithRegion(r Region) RegionalAPI
}

// RegionalAPI exposes the API for a region
type RegionalAPI interface {
	Summoner(ctx context.Context, name string) (*Summoner, error)
	SummonerByAccount(ctx context.Context, account string) (*Summoner, error)

	Matches(ctx context.Context, account string) ([]Match, error)
	MatchDetails(ctx context.Context, matchID string) (*MatchDetails, error)
	MatchTimeline(ctx context.Context, matchID string) (*MatchTimeline, error)
}

type (
	// Summoner represents a user
	Summoner struct {
		AccountID     string `json:"accountId"`
		Name          string `json:"name"`
		SummonerLevel int    `json:"summonerLevel"`
		ProfileIconID int    `json:"profileIconId"`
		PlayerID      string `json:"id"`
	}
)
