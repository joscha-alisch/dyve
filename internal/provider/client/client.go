package client

import (
	"encoding/json"
	"fmt"
	"github.com/joscha-alisch/dyve/pkg/provider/sdk"
	"net/http"
)

type listAppsResponse struct {
	Status int
	Err    string
	Result []sdk.App
}

type getAppResponse struct {
	Status int
	Err    string
	Result sdk.App
}

func NewAppProviderClient(uri string, c *http.Client) sdk.AppProvider {
	if c == nil {
		c = http.DefaultClient
	}

	return &appProviderClient{
		basePath: uri,
		c:        c,
	}
}

type appProviderClient struct {
	c        *http.Client
	basePath string
}

func (a *appProviderClient) ListApps() ([]sdk.App, error) {
	req, err := http.NewRequest("GET", a.basePath+"/apps", nil)
	if err != nil {
		return nil, err
	}

	res, err := a.c.Do(req)
	if err != nil {
		return nil, err
	}

	r := listAppsResponse{}

	err = json.NewDecoder(res.Body).Decode(&r)
	if err != nil {
		return nil, err
	}

	return r.Result, nil
}

func (a *appProviderClient) GetApp(id string) (sdk.App, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/apps/%s", a.basePath, id), nil)
	if err != nil {
		return sdk.App{}, err
	}

	res, err := a.c.Do(req)
	if err != nil {
		return sdk.App{}, err
	}

	if res.StatusCode == http.StatusNotFound {
		return sdk.App{}, sdk.ErrNotFound
	}

	r := getAppResponse{}

	err = json.NewDecoder(res.Body).Decode(&r)
	if err != nil {
		return sdk.App{}, err
	}

	return r.Result, nil
}

func (a *appProviderClient) Search(term string, limit int) ([]sdk.AppSearchResult, error) {
	panic("implement me")
}
