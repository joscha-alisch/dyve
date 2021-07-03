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
	panic("implement me")
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
	panic("implement me")
}
