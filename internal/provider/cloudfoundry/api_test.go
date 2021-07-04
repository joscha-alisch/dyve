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
		expected *Org
	}{
		{"gets org", "org-a", cfBackend{
			orgs:   map[string]*cf.Org{"org-a": {Guid: "org-a", Name: "org-name"}},
			spaces: map[string]*cf.Space{"space-a": {Guid: "space-a", OrganizationGuid: "org-a"}},
		}, &Org{
			Guid:   "org-a",
			Name:   "org-name",
			Spaces: []string{"space-a"},
		}},
	}

	for _, test := range tests {
		t.Run(test.desc, func(tt *testing.T) {
			api := NewApi(&fakeCfClient{test.state})

			o, _ := api.GetOrg(test.guid)
			if !cmp.Equal(*test.expected, o) {
				tt.Errorf("\norg was different: \n%s\n", cmp.Diff(*test.expected, o))
			}
		})
	}
}

func TestGetSpace(t *testing.T) {
	tests := []struct {
		desc          string
		guid          string
		state         cfBackend
		expectedSpace *Space
		expectedApps  []App
	}{
		{"gets space", "space-a", cfBackend{
			spaces: map[string]*cf.Space{"space-a": {Guid: "space-a", Name: "space-name", OrganizationGuid: "org-a"}},
			apps:   map[string]*cf.App{"app-a": {Guid: "app-a", Name: "app-name", SpaceGuid: "space-a"}},
		}, &Space{
			Guid: "space-a",
			Org:  "org-a",
			Name: "space-name",
			Apps: []string{"app-a"},
		}, []App{
			{Guid: "app-a", Name: "app-name", Space: "space-a", Org: "org-a"},
		}},
	}

	for _, test := range tests {
		t.Run(test.desc, func(tt *testing.T) {
			api := NewApi(&fakeCfClient{test.state})

			s, apps, _ := api.GetSpace(test.guid)
			if !cmp.Equal(*test.expectedSpace, s) {
				tt.Errorf("\nspace was different: \n%s\n", cmp.Diff(*test.expectedSpace, s))
			}
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
		expected *CFInfo
	}{
		{"gets info", cfBackend{
			orgs: map[string]*cf.Org{"org-a": {Guid: "org-a", Name: "org-name"}},
		}, &CFInfo{
			Orgs: []string{"org-a"},
		}},
	}

	for _, test := range tests {
		t.Run(test.desc, func(tt *testing.T) {
			api := NewApi(&fakeCfClient{test.state})

			o, _ := api.GetCFInfo()
			if !cmp.Equal(*test.expected, o) {
				tt.Errorf("\ninfo was different: \n%s\n", cmp.Diff(*test.expected, o))
			}
		})
	}
}

type fakeCfClient struct {
	b cfBackend
}

type cfBackend struct {
	orgs   map[string]*cf.Org
	spaces map[string]*cf.Space
	apps   map[string]*cf.App
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
