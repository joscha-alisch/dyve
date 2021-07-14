package client

import (
	"encoding/json"
	"fmt"
	"github.com/joscha-alisch/dyve/pkg/provider/sdk"
	"net/http"
	"strings"
)

func newBaseClient(basePath string, c *http.Client) baseClient {
	if c == nil {
		c = http.DefaultClient
	}

	if !strings.HasSuffix(basePath, "/") {
		basePath = basePath + "/"
	}

	return baseClient{
		c:        c,
		basePath: basePath,
	}
}

type baseClient struct {
	c        *http.Client
	basePath string
}

func (a *baseClient) get(resp interface{}, query map[string]string, path ...string) error {
	fullPath := a.basePath
	for _, s := range path {
		fullPath = fullPath + "/" + s
	}

	var queries []string
	for k, v := range query {
		queries = append(queries, fmt.Sprintf("%s=%s", k, v))
	}
	if len(queries) != 0 {
		fullPath += "?"
		fullPath += strings.Join(queries, "&")
	}

	req, err := http.NewRequest("GET", fullPath, nil)
	if err != nil {
		return err
	}

	res, err := a.c.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode == http.StatusNotFound {
		return sdk.ErrNotFound
	}

	err = json.NewDecoder(res.Body).Decode(&resp)
	if err != nil {
		return err
	}

	return nil
}
