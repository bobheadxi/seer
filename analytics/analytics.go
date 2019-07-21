package analytics

type Analytics struct {
	Team    *TeamAnalytics              `json:"team"`
	Players map[string]*PlayerAnalytics `json:"players"`
}

type PlayerAnalytics struct {
	Overview     *PlayerOverview     `json:"overview"`
	Aggregations *PlayerAggregations `json:"aggregations"`
}

type PlayerOverview struct {
	SeasonTier string `json:"season_tier"`
	PeakTier   string `json:"peak_tier"`
}

type PlayerAggregations struct {
}

type TeamAnalytics struct {
	Overview     *TeamOverview     `json:"overview"`
	Aggregations *TeamAggregations `josn:"aggregations"`
}

type TeamOverview struct {
}

type TeamAggregations struct {
}
