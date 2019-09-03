package analytics

// PlayerAnalytics is a container for analytics on a single player
type PlayerAnalytics struct {
	AccountID    string              `json:"accountId"`
	Tier         string              `json:"seasonTier"`
	Aggregations *PlayerAggregations `json:"aggregations"`
}

// PlayerAggregations define aggregations for a player
type PlayerAggregations struct {
	Champions []ChampionAggregation `json:"champions"` // per champion in top 5
}

// ChampionAggregation compiles data about a played champion
type ChampionAggregation struct {
	ChampionID string `json:"championId"`
	Count      int    `json:"count"`

	Lanes []Count `json:"lanes"`
	Roles []Count `json:"roles"`

	Perks     []Count `json:"perks"`
	Summoners []Count `json:"summoners"`
	// Items     []Count `json:"items"`

	// Games GameAggregations `json:"games"`
	// Timelines TimelineAggregations `json:"timelines"`
}
