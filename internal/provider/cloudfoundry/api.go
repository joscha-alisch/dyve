package cloudfoundry

import cf "github.com/cloudfoundry-community/go-cfclient"

type Login struct {
	Api  string
	User string
	Pass string
}

/**
A cloudfoundry.API is a wrapper around the official client.
*/
type API interface {
	GetOrg(guid string) (Org, error)
	GetSpace(guid string) (Space, []App, error)
	GetApp(guid string) (App, error)
}

type CfCli interface {
	ListOrgs() ([]cf.Org, error)
	GetOrgByGuid(guid string) (cf.Org, error)
	GetSpaceByGuid(guid string) (cf.Space, error)
	ListSpacesByOrgGuid(orgGuid string) ([]cf.Space, error)
	ListAppsBySpaceGuid(spaceGuid string) ([]cf.App, error)
}

func NewApi(cli CfCli) API {
	return &api{
		cli: cli,
	}
}

type api struct {
	cli CfCli
}

func (a *api) GetApp(guid string) (App, error) {
	return App{}, nil
}

func (a *api) GetSpace(guid string) (Space, []App, error) {
	return Space{}, nil, nil
}

func (a *api) GetOrg(guid string) (Org, error) {
	o, err := a.cli.GetOrgByGuid(guid)
	if err != nil {
		return Org{}, err
	}

	spaces, err := a.cli.ListSpacesByOrgGuid(o.Guid)
	if err != nil {
		return Org{}, err
	}

	var spaceGuids []string
	for _, space := range spaces {
		spaceGuids = append(spaceGuids, space.Guid)
	}

	return Org{
		Guid:   o.Guid,
		Name:   o.Name,
		Spaces: spaceGuids,
	}, nil
}
