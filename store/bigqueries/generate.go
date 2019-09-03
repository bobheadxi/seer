package bigqueries

//go:generate go run github.com/UnnoTed/fileb0x b0x.yml

const (
	// TeamMatchesView is the source for the team matches view query
	TeamMatchesView = "sql/team_matches_view.sql"
	// TeamGamesQuery is the source for the team games query
	TeamGamesQuery = "sql/team_games.sql"

	// AnalyticsQuery is the source for the team analytics query
	AnalyticsQuery = "sql/player_analytics.sql"
)
