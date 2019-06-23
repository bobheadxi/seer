package riot

import (
	"fmt"
)

const (
	pathSummonerByName    = "/lol/summoner/v4/summoners/by-name/"
	pathSummonerByAccount = "/lol/summoner/v4/summoners/by-account/"
	pathMatchesByAccount  = "/lol/match/v4/matchlists/by-account/"
	pathMatchByMatchID    = "/lol/match/v4/matches/"
)

func apiHost(region Region) string {
	return fmt.Sprintf("https://%s.api.riotgames.com", region)
}
