package analytics

// PlayerAnalytics is a container for analytics on a single player
type PlayerAnalytics struct {
	Tier         string              `json:"season_tier"`
	Aggregations *PlayerAggregations `json:"aggregations"`
}

type PlayerAggregations struct {
	Lanes []string `json:"lanes"` // top 2
	Roles []string `json:"roles"` // top 2

	Champions map[string]ChampionAggregation `json:"champions"` // per champion in top 5
}

type Perks struct {
	Primary   string `json:"primary"`
	Secondary string `json:"secondary"`
}

type Summoners struct {
	Spell1 string `json:"spell1"`
	Spell2 string `json:"spell2"`
}
