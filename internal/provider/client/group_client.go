package client

import (
	"github.com/joscha-alisch/dyve/pkg/provider/sdk"
	"net/http"
)

func NewGroupProviderClient(uri string, c *http.Client) sdk.GroupProvider {
	return &groupProviderClient{
		baseClient: newBaseClient(uri+"/groups", c),
	}
}

type listGroupsResponse struct {
	Status int
	Err    string
	Result []sdk.Group
}

type getGroupResponse struct {
	Status int
	Err    string
	Result sdk.Group
}

type groupProviderClient struct {
	baseClient
}

func (p *groupProviderClient) ListGroups() ([]sdk.Group, error) {
	r := listGroupsResponse{}
	err := p.get(&r, nil)
	if err != nil {
		return nil, err
	}
	return r.Result, nil
}

func (p *groupProviderClient) GetGroup(id string) (sdk.Group, error) {
	r := getGroupResponse{}
	err := p.get(&r, nil, id)
	if err != nil {
		return sdk.Group{}, err
	}
	return r.Result, nil
}
