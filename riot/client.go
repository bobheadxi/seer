package riot

import (
	"context"
	"net/http"
	"path"
	"sync"

	"go.uber.org/zap"
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
func NewClient(l *zap.Logger, auth func() string) (API, error) {
	return &client{
		l:       l,
		http:    &httpClient{l.Named("http"), http.DefaultClient, auth},
		regions: make(map[Region]*regionalClient),
	}, nil
}

func (c *client) WithRegion(r Region) RegionalAPI {
	c.rmux.RLock()
	regional, exists := c.regions[r]
	c.rmux.RUnlock()
	if !exists {
		regional = &regionalClient{
			host: apiHost(r),
			http: c.http,
		}
		c.rmux.Lock()
		c.regions[r] = regional
		c.rmux.Unlock()
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

type matchesResponse struct {
	Matches []Match `json:"matches"`
}

func (c *regionalClient) Matches(ctx context.Context, account string) ([]Match, error) {
	var resp matchesResponse
	if err := c.http.Get(ctx,
		path.Join(c.host, pathMatchesByAccount, account),
		newParams().Q("queue", queue5v5Draft).Q("queue", queue5v5Flex).Q("queue", queueClash),
		&resp,
	); err != nil {
		return nil, err
	}
	return resp.Matches, nil
}

func (c *regionalClient) MatchDetails(ctx context.Context, matchID string) (*MatchDetails, error) {
	var details MatchDetails
	return &details, c.http.Get(ctx,
		path.Join(c.host, pathMatchByMatchID, matchID),
		nil,
		&details)
}

func (c *regionalClient) MatchTimeline(ctx context.Context, matchID string) (*MatchTimeline, error) {
	var details MatchTimeline
	return &details, c.http.Get(ctx,
		path.Join(c.host, pathTimelineByMatchID, matchID),
		nil,
		&details)
}
