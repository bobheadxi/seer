package analytics

// Analytics is the root container for Seer analytics
type Analytics struct {
	Team    map[string]*TeamAnalytics              `json:"team"`    // per season
	Players map[string]map[string]*PlayerAnalytics `json:"players"` // per season, per player
}

// Aggregate is a generic container for basic aggregations
type Aggregate struct {
	Average float64 `json:"average"`
	Max     float64 `json:"max"`
	Min     float64 `json:"min"`
	StdDev  float64 `json:"stddev"`
}

// ChampionAggregation compiles data about a played champion
type ChampionAggregation struct {
	Lane string `json:"lane"`
	Role string `json:"role"`

	Perks     Perks     `json:"perks"`
	Summoners Summoners `json:"summoners"`
	Items     []string  `json:"items"`

	Games GameAggregations `json:"games"`
	// Timelines TimelineAggregations `json:"timelines"`
}

/*
type TimelineAggregations struct { }
*/

// GameAggregations compiles a basic overview of games
type GameAggregations struct {
	Kills   Aggregate
	Deaths  Aggregate
	Assists Aggregate

	Vision Aggregate
	Gold   Aggregate

	DamageDealt Aggregate
	DamageTaken Aggregate

	Minions        Aggregate
	JungleFriendly Aggregate
	JungleEnemy    Aggregate
}
