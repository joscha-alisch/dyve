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
	GetCFInfo() (CFInfo, error)
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

func NewDefaultApi(l Login) (API, error) {
	cli, err := cf.NewClient(&cf.Config{
		ApiAddress: l.Api,
		Username:   l.User,
		Password:   l.Pass,
	})
	if err != nil {
		return nil, err
	}

	return &api{cli: cli}, nil
}

func NewApi(cli CfCli) API {
	return &api{
		cli: cli,
	}
}

type api struct {
	cli CfCli
}

func (a *api) GetCFInfo() (CFInfo, error) {
	orgs, err := a.cli.ListOrgs()
	if err != nil {
		return CFInfo{}, err
	}

	var res []string
	for _, org := range orgs {
		res = append(res, org.Guid)
	}

	return CFInfo{
		Orgs: res,
	}, nil
}

func (a *api) GetApp(guid string) (App, error) {
	return App{}, nil
}

func (a *api) GetSpace(guid string) (Space, []App, error) {
	s, err := a.cli.GetSpaceByGuid(guid)
	if err != nil {
		return Space{}, nil, err
	}

	cfApps, err := a.cli.ListAppsBySpaceGuid(s.Guid)
	if err != nil {
		return Space{}, nil, err
	}

	var appGuids []string
	var apps []App
	for _, app := range cfApps {
		appGuids = append(appGuids, app.Guid)
		apps = append(apps, App{
			Guid:  app.Guid,
			Name:  app.Name,
			Org:   s.OrganizationGuid,
			Space: s.Guid,
		})
	}
	return Space{
		Guid: s.Guid,
		Org:  s.OrganizationGuid,
		Name: s.Name,
		Apps: appGuids,
	}, apps, nil
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
