package cloudfoundry

import (
	cf "github.com/cloudfoundry-community/go-cfclient"
)

const CFGuid = "main"

type CFLogin struct {
	Api  string
	User string
	Pass string
}

/**
A cloudfoundry.API is a wrapper around the official client.
*/
type API interface {
	ListOrgs() ([]Org, error)
	ListSpaces(orgGuid string) ([]Space, error)
	ListApps(spaceGuid string) ([]App, error)
	GetApp(guid string) (App, error)
}

type CfCli interface {
	ListOrgs() ([]cf.Org, error)
	GetOrgByGuid(guid string) (cf.Org, error)
	GetSpaceByGuid(guid string) (cf.Space, error)
	ListSpacesByOrgGuid(orgGuid string) ([]cf.Space, error)
	ListAppsBySpaceGuid(spaceGuid string) ([]cf.App, error)
}

func NewDefaultApi(l CFLogin) (API, error) {
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

func (a *api) ListOrgs() ([]Org, error) {
	orgs, err := a.cli.ListOrgs()
	if err != nil {
		return nil, err
	}

	var res []Org
	for _, org := range orgs {
		res = append(res, Org{
			OrgInfo: OrgInfo{
				Guid: org.Guid,
				Name: org.Name,
				Cf: CFInfo{
					Guid: CFGuid,
				},
			},
		})
	}

	return res, nil
}

func (a *api) GetApp(guid string) (App, error) {
	return App{}, nil
}

func (a *api) ListApps(spaceGuid string) ([]App, error) {
	cfApps, err := a.cli.ListAppsBySpaceGuid(spaceGuid)
	if err != nil {
		return nil, err
	}

	var apps []App
	for _, app := range cfApps {
		apps = append(apps, App{
			AppInfo{
				Guid: app.Guid,
				Name: app.Name,
				Space: SpaceInfo{
					Guid: spaceGuid,
				},
			},
		})
	}

	return apps, nil
}

func (a *api) ListSpaces(orgGuid string) ([]Space, error) {
	spaces, err := a.cli.ListSpacesByOrgGuid(orgGuid)
	if err != nil {
		return nil, err
	}

	var res []Space
	for _, space := range spaces {
		res = append(res, Space{
			SpaceInfo: SpaceInfo{
				Guid: space.Guid,
				Name: space.Name,
				Org: OrgInfo{
					Guid: orgGuid,
				},
			},
		})
	}

	return res, nil
}
