package client

import (
	"github.com/joscha-alisch/dyve/pkg/provider/sdk"
	"net/http"
)

type getInstancesResponse struct {
	Status int
	Err    string
	Result sdk.AppInstances
}

func NewInstancesProviderClient(uri string, c *http.Client) sdk.InstancesProvider {
	return &instancesProviderClient{
		baseClient: newBaseClient(uri+"/instances", c),
	}
}

type instancesProviderClient struct {
	baseClient
}

func (c *instancesProviderClient) GetAppInstances(id string) (sdk.AppInstances, error) {
	r := getInstancesResponse{}
	err := c.get(&r, nil, id)
	if err != nil {
		return nil, err
	}
	return r.Result, nil
}
