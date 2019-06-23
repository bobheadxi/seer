package riot

import (
	"context"
	"path"
	"sync"

	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

const (
	queue5v5Flex  = "440"
	queue5v5Draft = "400"
	queueClash    = "700"
)

type client struct {
	l    *zap.Logger
	http *httpClient

	regions map[Region]*regionalClient
	rmux    sync.RWMutex
}

// NewClient instantiates an API implementation
func NewClient(l *zap.Logger, auth oauth2.TokenSource) API {
	httpC := oauth2.NewClient(context.Background(), auth)
	return &client{
		l:       l,
		http:    &httpClient{httpC},
		regions: make(map[Region]*regionalClient),
	}
}

func (c *client) WithRegion(r Region) RegionalAPI {
	c.rmux.RLock()
	regional, exists := c.regions[r]
	c.rmux.Unlock()
	if !exists {
		regional = &regionalClient{
			host: apiHost(r),
			http: c.http,
		}
		c.regions[r] = regional
	}
	return regional
}

type regionalClient struct {
	host string
	http *httpClient
}

func (c *regionalClient) Summoner(ctx context.Context, name string) (*Summoner, error) {
	var sum Summoner
	return &sum, c.http.Get(ctx,
		path.Join(c.host, pathSummonerByName, name),
		nil,
		&sum)
}

func (c *regionalClient) SummonerByAccount(ctx context.Context, account string) (*Summoner, error) {
	var sum Summoner
	return &sum, c.http.Get(ctx,
		path.Join(c.host, pathSummonerByAccount, account),
		nil,
		&sum)
}

func (c *regionalClient) Matches(ctx context.Context, account string) ([]Match, error) {
	var matches []Match
	return matches, c.http.Get(ctx,
		path.Join(c.host, pathMatchesByAccount, account),
		newParams().Q("queue", queue5v5Draft).Q("queue", queue5v5Flex).Q("queue", queueClash),
		&matches)
}

func (c *regionalClient) MatchDetails(ctx context.Context, matchID string) (*MatchDetails, error) {
	var details MatchDetails
	return &details, c.http.Get(ctx,
		path.Join(c.host, pathMatchByMatchID, matchID),
		nil,
		&details)
}
