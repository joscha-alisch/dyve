package cloudfoundry

import (
	cf "github.com/cloudfoundry-community/go-cfclient"
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestGetOrg(t *testing.T) {
	tests := []struct {
		desc     string
		guid     string
		state    cfBackend
		expected []Space
	}{
		{"gets org spaces", "org-a", cfBackend{
			orgs: map[string]*cf.Org{"org-a": {Guid: "org-a", Name: "org-name"}},
			spaces: map[string]*cf.Space{"space-a": {
				Guid: "space-a", Name: "space-name", OrganizationGuid: "org-a",
			}},
		}, []Space{
			{SpaceInfo: SpaceInfo{Guid: "space-a", Name: "space-name", Org: OrgInfo{Guid: "org-a"}}},
		}},
	}

	for _, test := range tests {
		t.Run(test.desc, func(tt *testing.T) {
			api := NewApi(&fakeCfClient{test.state})

			spaces, _ := api.ListSpaces(test.guid)
			if !cmp.Equal(test.expected, spaces) {
				tt.Errorf("\nspaces were different: \n%s\n", cmp.Diff(test.expected, spaces))
			}
		})
	}
}

func TestGetSpace(t *testing.T) {
	tests := []struct {
		desc         string
		guid         string
		state        cfBackend
		expectedApps []App
	}{
		{"gets space", "space-a", cfBackend{
			spaces: map[string]*cf.Space{"space-a": {Guid: "space-a", Name: "space-name", OrganizationGuid: "org-a"}},
			apps:   map[string]*cf.App{"app-a": {Guid: "app-a", Name: "app-name", SpaceGuid: "space-a"}},
		}, []App{
			{AppInfo{Guid: "app-a", Name: "app-name", Space: SpaceInfo{Guid: "space-a"}}},
		}},
	}

	for _, test := range tests {
		t.Run(test.desc, func(tt *testing.T) {
			api := NewApi(&fakeCfClient{test.state})

			apps, _ := api.ListApps(test.guid)
			if !cmp.Equal(test.expectedApps, apps) {
				tt.Errorf("\napps were different: \n%s\n", cmp.Diff(test.expectedApps, apps))
			}
		})
	}
}

func TestGetCfInfo(t *testing.T) {
	tests := []struct {
		desc     string
		state    cfBackend
		expected []Org
	}{
		{"gets info", cfBackend{
			orgs: map[string]*cf.Org{"org-a": {Guid: "org-a", Name: "org-name"}},
		}, []Org{
			{OrgInfo: OrgInfo{Guid: "org-a", Name: "org-name", Cf: CFInfo{Guid: CFGuid}}},
		}},
	}

	for _, test := range tests {
		t.Run(test.desc, func(tt *testing.T) {
			api := NewApi(&fakeCfClient{test.state})

			orgs, _ := api.ListOrgs()
			if !cmp.Equal(test.expected, orgs) {
				tt.Errorf("\norgs were different: \n%s\n", cmp.Diff(test.expected, orgs))
			}
		})
	}
}

func TestGetInstances(t *testing.T) {
	tests := []struct {
		desc        string
		id          string
		state       cfBackend
		expected    Instances
		expectedErr error
	}{
		{desc: "gets instances", id: "app-a", state: cfBackend{
			instances: map[string]map[string]cf.AppInstance{
				"app-a": {"0": {State: "STOPPED"}},
			},
		}, expected: Instances{
			{
				State: "STOPPED",
			},
		}},
	}

	for _, test := range tests {
		t.Run(test.desc, func(tt *testing.T) {
			api := NewApi(&fakeCfClient{test.state})

			res, _ := api.GetInstances(test.id)
			if !cmp.Equal(test.expected, res) {
				tt.Errorf("\nresult mismatch: \n%s\n", cmp.Diff(test.expected, res))
			}
		})
	}
}

func TestGetRoutes(t *testing.T) {
	tests := []struct {
		desc        string
		id          string
		state       cfBackend
		expected    Routes
		expectedErr error
	}{
		{desc: "not found", id: "app-a", state: cfBackend{
			routes: map[string][]cf.Route{
				"app-a": {},
			},
		}, expected: nil},
		{desc: "gets routes", id: "app-a", state: cfBackend{
			routes: map[string][]cf.Route{
				"app-a": {},
			},
		}, expected: nil},
	}

	for _, test := range tests {
		t.Run(test.desc, func(tt *testing.T) {
			api := NewApi(&fakeCfClient{test.state})

			res, _ := api.GetRoutes(test.id)
			if !cmp.Equal(test.expected, res) {
				tt.Errorf("\nresult mismatch: \n%s\n", cmp.Diff(test.expected, res))
			}
		})
	}
}

type fakeCfClient struct {
	b cfBackend
}

func (f *fakeCfClient) GetAppRoutes(appGuid string) ([]cf.Route, error) {
	if f.b.routes[appGuid] == nil {
		return nil, errNotFound
	}
	return f.b.routes[appGuid], nil
}

func (f *fakeCfClient) GetAppInstances(guid string) (map[string]cf.AppInstance, error) {
	if f.b.instances[guid] == nil {
		return nil, errNotFound
	}
	return f.b.instances[guid], nil
}

type cfBackend struct {
	orgs      map[string]*cf.Org
	spaces    map[string]*cf.Space
	apps      map[string]*cf.App
	routes    map[string][]cf.Route
	instances map[string]map[string]cf.AppInstance
}

func (f *fakeCfClient) ListOrgs() ([]cf.Org, error) {
	var orgs []cf.Org
	for _, org := range f.b.orgs {
		orgs = append(orgs, *org)
	}
	return orgs, nil
}

func (f *fakeCfClient) GetOrgByGuid(guid string) (cf.Org, error) {
	return *f.b.orgs[guid], nil
}

func (f *fakeCfClient) GetSpaceByGuid(guid string) (cf.Space, error) {
	return *f.b.spaces[guid], nil
}

func (f *fakeCfClient) ListSpacesByOrgGuid(orgGuid string) ([]cf.Space, error) {
	var res []cf.Space
	for _, space := range f.b.spaces {
		if space.OrganizationGuid == orgGuid {
			res = append(res, *space)
		}
	}
	return res, nil
}

func (f *fakeCfClient) ListAppsBySpaceGuid(spaceGuid string) ([]cf.App, error) {
	var res []cf.App
	for _, app := range f.b.apps {
		if app.SpaceGuid == spaceGuid {
			res = append(res, *app)
		}
	}
	return res, nil
}
