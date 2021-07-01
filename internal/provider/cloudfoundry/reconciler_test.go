package cloudfoundry

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestReconciler(t *testing.T) {
	tests := []struct {
		desc  string
		db    fakeDb
		cf    fakeCf
		sleep bool
	}{
		{
			desc: "updates org",
			db: fakeDb{
				job: &ReconcileJob{Type: ReconcileOrg, Guid: "abc"},
			},
			cf: fakeCf{
				orgs: map[string]*Org{"abc": {"org"}},
			},
		},
		{
			desc: "updates space",
			db: fakeDb{
				job: &ReconcileJob{Type: ReconcileSpace, Guid: "abc"},
			},
			cf: fakeCf{
				spaces: map[string]*Space{"abc": {"space"}},
			},
		},
		{
			desc: "updates app",
			db: fakeDb{
				job: &ReconcileJob{Type: ReconcileApp, Guid: "abc"},
			},
			cf: fakeCf{
				apps: map[string]*App{"abc": {"app"}},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(tt *testing.T) {
			r := NewReconciler(&test.db, &test.cf)
			_, _ = r.Run()

			expected := map[string]interface{}{
				"orgs":   test.cf.orgs,
				"spaces": test.cf.spaces,
				"apps":   test.cf.apps,
			}
			res := map[string]interface{}{
				"orgs":   test.db.orgs,
				"spaces": test.db.spaces,
				"apps":   test.db.apps,
			}

			if !cmp.Equal(expected, res) {
				tt.Errorf("\ncf api and reconciled db differ in orgs: \n%s", cmp.Diff(expected, res))
			}
		})
	}

}

type fakeCf struct {
	orgs   map[string]*Org
	spaces map[string]*Space
	apps   map[string]*App
}

func (f *fakeCf) GetApp(guid string) (App, error) {
	return *f.apps[guid], nil
}

func (f *fakeCf) GetSpace(guid string) (Space, error) {
	return *f.spaces[guid], nil
}

func (f *fakeCf) GetOrg(guid string) (Org, error) {
	return *f.orgs[guid], nil
}

type fakeDb struct {
	job    *ReconcileJob
	orgs   map[string]*Org
	spaces map[string]*Space
	apps   map[string]*App
}

func (f *fakeDb) UpdateApp(guid string, a App) error {
	if f.apps == nil {
		f.apps = make(map[string]*App)
	}
	f.apps[guid] = &a
	return nil
}

func (f *fakeDb) UpdateSpace(guid string, s Space) error {
	if f.spaces == nil {
		f.spaces = make(map[string]*Space)
	}
	f.spaces[guid] = &s
	return nil
}

func (f *fakeDb) UpdateOrg(guid string, o Org) error {
	if f.orgs == nil {
		f.orgs = make(map[string]*Org)
	}
	f.orgs[guid] = &o
	return nil
}

func (f *fakeDb) FetchReconcileJob() *ReconcileJob {
	return f.job
}
