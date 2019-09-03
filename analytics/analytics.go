package analytics

// Analytics is the root container for Seer analytics
type Analytics struct {
	// Team    map[string]*TeamAnalytics   `json:"team"`    // per season
	Players []*PlayerAnalytics `json:"players"` // per player. TODO: per season?
}

// NewAnalytics instantiates a new base analytics object
func NewAnalytics() *Analytics {
	return &Analytics{
		Players: make([]*PlayerAnalytics, 0),
	}
}

// Aggregate is a generic container for basic aggregations
type Aggregate struct {
	Average float64 `json:"average"`
	Max     float64 `json:"max"`
	Min     float64 `json:"min"`
	StdDev  float64 `json:"stddev"`
}

// Count is a generic container for count aggregations
type Count struct {
	Value string `json:"value"`
	Count int    `json:"count"`
}

// TimelineAggregations contains analytics from timeline data
type TimelineAggregations struct {
	// build order, lane/jungle path, etc.
}

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
