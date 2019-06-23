package riot

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
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
	http *http.Client
}

func (c *httpClient) Get(ctx context.Context, url string, query params, out interface{}) error {
	req, err := http.NewRequest(http.MethodGet, url+query.Build(), nil)
	if err != nil {
		return err
	}

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
	return json.Unmarshal(b, out)
}

func codeToErr(c int) error {
	// TODO
	return nil
}
