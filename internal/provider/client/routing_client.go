package client

import (
	"github.com/joscha-alisch/dyve/pkg/provider/sdk"
	"net/http"
)

type getRoutingResponse struct {
	Status int
	Err    string
	Result sdk.AppRouting
}

func NewRoutingProviderClient(uri string, c *http.Client) sdk.RoutingProvider {
	return &routingProviderClient{
		baseClient: newBaseClient(uri+"/routing", c),
	}
}

type routingProviderClient struct {
	baseClient
}

func (c *routingProviderClient) GetAppRouting(id string) (sdk.AppRouting, error) {
	r := getRoutingResponse{}
	err := c.get(&r, nil, id)
	if err != nil {
		return sdk.AppRouting{}, err
	}
	return r.Result, nil
}
