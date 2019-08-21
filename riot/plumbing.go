package riot

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
)

type params []string

func newParams() params { return []string{} }

func (p params) Q(key, val string) params {
	p = append(p, key+"="+val)
	return p
}

func (p params) Build() string {
	if len(p) == 0 {
		return ""
	}
	return "?" + strings.Join(p, "&")
}

type httpClient struct {
	l    *zap.Logger
	http *http.Client

	auth func() string
}

func (c *httpClient) Get(ctx context.Context, url string, query params, out interface{}) error {
	log := c.l.With(zap.String("request.id", middleware.GetReqID(ctx)), zap.String("target.url", url))

	req, err := http.NewRequest(http.MethodGet, "https://"+url+query.Build(), nil)
	if err != nil {
		return err
	}
	log.Debug("request created", zap.Any("target.parsed_url", req.URL))
	req.Header.Set("X-Riot-Token", c.auth())

	resp, err := c.http.Do(req.WithContext(ctx))
	if err != nil {
		return err
	}
	if err := codeToErr(resp.StatusCode); err != nil {
		return err
	}

	b, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return err
	}
	log.Debug("request completed", zap.Int("body.length", len(b)), zap.Int("code", resp.StatusCode))

	return json.Unmarshal(b, out)
}

func codeToErr(c int) error {
	// TODO
	return nil
}
