package cloudfoundry

import (
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	recon "github.com/joscha-alisch/dyve/internal/reconciliation"
	"testing"
	"time"
)

func TestReconciler(t *testing.T) {
	tests := []struct {
		desc  string
		db    fakeDb
		cf    fakeCf
		sleep bool
		err   error
	}{
		{
			desc: "updates org",
			db: fakeDb{
				job: &recon.Job{Type: ReconcileSpaces, Guid: "org-a-guid"},
				b: backend{
					Orgs: map[string]*Org{"org-a-guid": {
						OrgInfo: OrgInfo{Name: "org", Guid: "org-a-guid"},
					}},
				},
			},
			cf: fakeCf{b: backend{
				Orgs: map[string]*Org{"org-a-guid": {OrgInfo: OrgInfo{Name: "org", Guid: "org-a-guid"}}},
				Spaces: map[string]*Space{
					"space-a-guid": {SpaceInfo: SpaceInfo{Guid: "space-a-guid", Org: OrgInfo{Guid: "org-a-guid"}}},
					"space-b-guid": {SpaceInfo: SpaceInfo{Guid: "space-b-guid", Org: OrgInfo{Guid: "org-a-guid"}}},
				},
			}},
		},
		{
			desc: "updates orgs",
			db: fakeDb{
				job: &recon.Job{Type: ReconcileOrganizations, Guid: "main"},
				b:   backend{CfApis: map[string]*CF{"main": {CFInfo: CFInfo{Guid: "main"}}}},
			},
			cf: fakeCf{b: backend{
				CfApis: map[string]*CF{"main": {CFInfo: CFInfo{Guid: "main"}}},
				Orgs: map[string]*Org{
					"org-a": {OrgInfo: OrgInfo{Guid: "org-a"}},
					"org-b": {OrgInfo: OrgInfo{Guid: "org-b"}},
				},
			}},
		},
		{
			desc: "updates space",
			db: fakeDb{
				job: &recon.Job{Type: ReconcileApps, Guid: "space-a"},
				b: backend{
					Spaces: map[string]*Space{"space-a": {SpaceInfo: SpaceInfo{Guid: "space-a"}}},
				},
			},
			cf: fakeCf{b: backend{
				Spaces: map[string]*Space{"space-a": {SpaceInfo: SpaceInfo{Guid: "space-a"}}},
				Apps: map[string]*App{
					"app-a": {AppInfo: AppInfo{Guid: "app-a", Space: SpaceInfo{Guid: "space-a"}}},
					"app-b": {AppInfo: AppInfo{Guid: "app-b", Space: SpaceInfo{Guid: "space-a"}}},
				},
			}},
		},
		{
			desc: "no work to be done",
			db: fakeDb{
				job: nil,
			},
			cf:    fakeCf{},
			sleep: true,
		},
		{
			desc: "removes org when not found",
			db: fakeDb{
				job: &recon.Job{Type: ReconcileSpaces, Guid: "not_exist"},
				b: backend{
					Orgs: map[string]*Org{
						"not_exist": {OrgInfo: OrgInfo{Guid: "not_exist"}},
					},
				},
			},
			cf: fakeCf{},
		},
		{
			desc: "removes space when not found",
			db: fakeDb{
				job: &recon.Job{Type: ReconcileApps, Guid: "not_exist"},
				b: backend{
					Spaces: map[string]*Space{
						"not_exist": {SpaceInfo: SpaceInfo{Guid: "not_exist"}},
					},
				},
			},
			cf: fakeCf{},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(tt *testing.T) {
			r := NewReconciler(&test.db, &test.cf, 1*time.Minute)
			worked, err := r.Run()
			if worked == test.sleep {
				tt.Errorf("\nexpected return: %v, was: %v", !test.sleep, worked)
			}

			if !cmp.Equal(test.err, err, cmpopts.EquateErrors()) {
				tt.Errorf("\nerr not as expected: \n%s", cmp.Diff(test.err, err, cmpopts.EquateErrors()))
			}

			if !cmp.Equal(test.cf.b, test.db.b) {
				tt.Errorf("\ncf api and reconciled db differ in orgs: \n%s", cmp.Diff(test.cf.b, test.db.b))
			}

		})
	}

}

type fakeCf struct {
	b backend
}

func (f *fakeCf) ListOrgs() ([]Org, error) {
	var res []Org
	for _, org := range f.b.Orgs {
		res = append(res, *org)
	}
	return res, nil
}

func (f *fakeCf) ListSpaces(orgGuid string) ([]Space, error) {
	var res []Space

	if f.b.Orgs[orgGuid] == nil {
		return nil, errNotFound
	}

	for _, space := range f.b.Spaces {
		if space.Org.Guid == orgGuid {
			res = append(res, *space)
		}
	}
	return res, nil
}

func (f *fakeCf) ListApps(spaceGuid string) ([]App, error) {
	if f.b.Spaces[spaceGuid] == nil {
		return nil, errNotFound
	}

	var res []App
	for _, app := range f.b.Apps {
		if app.Space.Guid == spaceGuid {
			res = append(res, *app)
		}
	}
	return res, nil
}

func (f *fakeCf) GetApp(guid string) (App, error) {
	if f.b.Apps[guid] == nil {
		return App{}, errNotFound
	}
	return *f.b.Apps[guid], nil
}

type fakeDb struct {
	job *recon.Job
	b   backend
}

func (f *fakeDb) GetApp(id string) (App, error) {
	if f.b.Apps[id] == nil {
		return App{}, errNotFound
	}

	return *f.b.Apps[id], nil
}

func (f *fakeDb) ListApps() ([]App, error) {
	var res []App
	for _, app := range f.b.Apps {
		res = append(res, *app)
	}
	return res, nil
}

func (f *fakeDb) UpsertOrgs(cfGuid string, orgs []Org) error {
	if f.b.Orgs == nil {
		f.b.Orgs = make(map[string]*Org)
	}
	for _, org := range orgs {
		org := org
		f.b.Orgs[org.Guid] = &org
	}
	return nil
}

func (f *fakeDb) UpsertOrgSpaces(orgGuid string, spaces []Space) error {
	if f.b.Spaces == nil {
		f.b.Spaces = make(map[string]*Space)
	}

	for _, space := range spaces {
		space := space
		f.b.Spaces[space.Guid] = &space
	}
	return nil
}

func (f *fakeDb) UpsertSpaceApps(spaceGuid string, apps []App) error {
	if f.b.Apps == nil {
		f.b.Apps = make(map[string]*App)
	}
	for _, app := range apps {
		app := app
		f.b.Apps[app.Guid] = &app
	}
	return nil
}

func (f *fakeDb) DeleteApp(guid string) {
	delete(f.b.Apps, guid)
	if len(f.b.Apps) == 0 {
		f.b.Apps = nil
	}
}

func (f *fakeDb) DeleteSpace(guid string) {
	delete(f.b.Spaces, guid)
	if len(f.b.Spaces) == 0 {
		f.b.Spaces = nil
	}
}

func (f *fakeDb) DeleteOrg(guid string) {
	delete(f.b.Orgs, guid)
	if len(f.b.Orgs) == 0 {
		f.b.Orgs = nil
	}
}

func (f *fakeDb) UpsertSpace(s Space) error {
	if f.b.Spaces == nil {
		f.b.Spaces = make(map[string]*Space)
	}
	f.b.Spaces[s.Guid] = &s
	return nil
}

func (f *fakeDb) AcceptReconcileJob(olderThan time.Duration) (recon.Job, bool) {
	if f.job == nil {
		return recon.Job{}, false
	}
	return *f.job, true
}

type backend struct {
	CfApis map[string]*CF
	Orgs   map[string]*Org
	Spaces map[string]*Space
	Apps   map[string]*App
}
