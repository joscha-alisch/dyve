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
			cf: fakeCf{b: backend{
				Orgs: map[string]*Org{"abc": {"org"}},
			}},
		},
		{
			desc: "updates space",
			db: fakeDb{
				job: &ReconcileJob{Type: ReconcileSpace, Guid: "abc"},
			},
			cf: fakeCf{b: backend{
				Spaces: map[string]*Space{"abc": {"space"}},
			}},
		},
		{
			desc: "updates app",
			db: fakeDb{
				job: &ReconcileJob{Type: ReconcileApp, Guid: "abc"},
			},
			cf: fakeCf{b: backend{
				Apps: map[string]*App{"abc": {"app"}},
			}},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(tt *testing.T) {
			r := NewReconciler(&test.db, &test.cf)
			_, _ = r.Run()

			if !cmp.Equal(test.cf.b, test.db.b) {
				tt.Errorf("\ncf api and reconciled db differ in orgs: \n%s", cmp.Diff(test.cf.b, test.db.b))
			}
		})
	}

}

type fakeCf struct {
	b backend
}

func (f *fakeCf) GetApp(guid string) (App, error) {
	return *f.b.Apps[guid], nil
}

func (f *fakeCf) GetSpace(guid string) (Space, error) {
	return *f.b.Spaces[guid], nil
}

func (f *fakeCf) GetOrg(guid string) (Org, error) {
	return *f.b.Orgs[guid], nil
}

type fakeDb struct {
	job *ReconcileJob
	b   backend
}

func (f *fakeDb) UpdateApp(guid string, a App) error {
	if f.b.Apps == nil {
		f.b.Apps = make(map[string]*App)
	}
	f.b.Apps[guid] = &a
	return nil
}

func (f *fakeDb) UpdateSpace(guid string, s Space) error {
	if f.b.Spaces == nil {
		f.b.Spaces = make(map[string]*Space)
	}
	f.b.Spaces[guid] = &s
	return nil
}

func (f *fakeDb) UpdateOrg(guid string, o Org) error {
	if f.b.Orgs == nil {
		f.b.Orgs = make(map[string]*Org)
	}
	f.b.Orgs[guid] = &o
	return nil
}

func (f *fakeDb) FetchReconcileJob() *ReconcileJob {
	return f.job
}

type backend struct {
	Orgs   map[string]*Org
	Spaces map[string]*Space
	Apps   map[string]*App
}
