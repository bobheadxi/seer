package analytics

type Analytics struct {
	Team    map[string]*TeamAnalytics              `json:"team"`    // per season
	Players map[string]map[string]*PlayerAnalytics `json:"players"` // per season, per player
}

type PlayerAnalytics struct {
	Tier         string              `json:"season_tier"`
	Aggregations *PlayerAggregations `json:"aggregations"`
}

type PlayerAggregations struct {
	Lanes []string `json:"lanes"` // top 2
	Roles []string `json:"roles"` // top 2

	Champions map[string]ChampionAggregation `json:"champions"` // per champion in top 5
}

type Aggregate struct {
	Average float64 `json:"average"`
	Max     float64 `json:"max"`
	Min     float64 `json:"min"`
	StdDev  float64 `json:"stddev"`
}

type Perks struct {
	Primary   string
	Secondary string
}

type Summoners struct {
	Spell1 string
	Spell2 string
}

type ChampionAggregation struct {
	Lane string
	Role string

	Perks     Perks
	Summoners Summoners
	Items     []string

	Games     GameAggregations
	Timelines TimelineAggregations
}

type TimelineAggregations struct {
	// TODO
}

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

type TeamAnalytics struct {
	// TODO
}
