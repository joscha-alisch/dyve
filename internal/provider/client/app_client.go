package client

import (
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
	return &appProviderClient{
		baseClient: newBaseClient(uri+"/apps", c),
	}
}

type appProviderClient struct {
	baseClient
}

func (a *appProviderClient) ListApps() ([]sdk.App, error) {
	r := listAppsResponse{}
	err := a.get(&r, nil)
	if err != nil {
		return nil, err
	}
	return r.Result, nil
}

func (a *appProviderClient) GetApp(id string) (sdk.App, error) {
	r := getAppResponse{}
	err := a.get(&r, nil, id)
	if err != nil {
		return sdk.App{}, err
	}
	return r.Result, nil
}
