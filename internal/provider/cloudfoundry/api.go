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

// API is an abstraction around the CloudFoundry functionality.
type API interface {
	ListOrgs() ([]Org, error)
	ListSpaces(orgGuid string) ([]Space, error)
	ListApps(spaceGuid string) ([]App, error)
	GetApp(guid string) (App, error)
	GetRoutes(appId string) (Routes, error)
	GetInstances(appId string) (Instances, error)
}

// CfCli is a wrapper interface for the official cloudfoundry client extracting the needed functions.
type CfCli interface {
	ListOrgs() ([]cf.Org, error)
	GetOrgByGuid(guid string) (cf.Org, error)
	GetSpaceByGuid(guid string) (cf.Space, error)
	ListSpacesByOrgGuid(orgGuid string) ([]cf.Space, error)
	ListAppsBySpaceGuid(spaceGuid string) ([]cf.App, error)
	GetAppRoutes(appGuid string) ([]cf.Route, error)
	GetAppInstances(guid string) (map[string]cf.AppInstance, error)
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

func (a *api) GetInstances(appId string) (Instances, error) {
	instances, err := a.cli.GetAppInstances(appId)
	if err != nil {
		return nil, err
	}

	var res Instances
	for _, instance := range instances {
		res = append(res, Instance{
			State: instance.State,
			Since: instance.Since.Time,
		})
	}
	return res, nil
}

func (a *api) GetRoutes(appId string) (Routes, error) {
	routes, err := a.cli.GetAppRoutes(appId)
	if err != nil {
		return nil, err
	}

	domains := map[string]string{}

	var res Routes
	for _, route := range routes {
		if domains[route.DomainGuid] == "" {
			domain, err := route.Domain()
			if err != nil {
				return nil, err
			}
			domains[route.DomainGuid] = domain.Name
		}

		res = append(res, Route{
			Host: route.Host + "." + domains[route.DomainGuid],
			Path: route.Path,
			Port: route.Port,
		})
	}
	return res, nil
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
