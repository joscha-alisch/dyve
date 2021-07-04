package cloudfoundry

import (
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
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
				job: &ReconcileJob{Type: ReconcileOrg, Guid: "abc"},
			},
			cf: fakeCf{b: backend{
				Orgs: map[string]*Org{"abc": {Name: "org", Guid: "abc"}},
			}},
		},
		{
			desc: "updates orgs",
			db: fakeDb{
				job: &ReconcileJob{Type: ReconcileCF},
			},
			cf: fakeCf{b: backend{
				Orgs: map[string]*Org{"abc": {Guid: "abc"}},
			}},
		},
		{
			desc: "updates space",
			db: fakeDb{
				job: &ReconcileJob{Type: ReconcileSpace, Guid: "abc"},
			},
			cf: fakeCf{b: backend{
				Spaces: map[string]*Space{"abc": {Name: "space", Guid: "abc"}},
			}},
		},
		{
			desc: "updates apps for space",
			db: fakeDb{
				job: &ReconcileJob{Type: ReconcileSpace, Guid: "abc"},
			},
			cf: fakeCf{b: backend{
				Spaces: map[string]*Space{"abc": {Name: "space", Guid: "abc", Apps: []string{
					"app-a",
				}}},
				Apps: map[string]*App{
					"app-a": {
						Guid:  "app-a",
						Name:  "a",
						Org:   "a",
						Space: "abc",
					},
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
				job: &ReconcileJob{Type: ReconcileOrg, Guid: "not_exist"},
				b: backend{
					Orgs: map[string]*Org{
						"not_exist": {Guid: "not_exist"},
					},
				},
			},
			cf: fakeCf{},
		},
		{
			desc: "removes space when not found",
			db: fakeDb{
				job: &ReconcileJob{Type: ReconcileSpace, Guid: "not_exist"},
				b: backend{
					Spaces: map[string]*Space{
						"not_exist": {Guid: "not_exist"},
					},
				},
			},
			cf: fakeCf{},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(tt *testing.T) {
			r := NewReconciler(&test.db, &test.cf)
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

func (f *fakeCf) GetCFInfo() (CFInfo, error) {
	var orgs []string
	for _, o := range f.b.Orgs {
		orgs = append(orgs, o.Guid)
	}
	return CFInfo{
		Orgs: orgs,
	}, nil
}

func (f *fakeCf) GetApp(guid string) (App, error) {
	if f.b.Apps[guid] == nil {
		return App{}, errNotFound
	}
	return *f.b.Apps[guid], nil
}

func (f *fakeCf) GetSpace(guid string) (Space, []App, error) {
	if f.b.Spaces[guid] == nil {
		return Space{}, nil, errNotFound
	}
	s := *f.b.Spaces[guid]
	var apps []App
	for _, app := range s.Apps {
		a, _ := f.GetApp(app)
		apps = append(apps, a)
	}
	return s, apps, nil
}

func (f *fakeCf) GetOrg(guid string) (Org, error) {
	if f.b.Orgs[guid] == nil {
		return Org{}, errNotFound
	}
	return *f.b.Orgs[guid], nil
}

type fakeDb struct {
	job *ReconcileJob
	b   backend
}

func (f *fakeDb) UpsertCfInfo(i CFInfo) error {
	for _, org := range i.Orgs {
		_ = f.UpsertOrg(Org{Guid: org})
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

func (f *fakeDb) UpsertApps(apps []App) error {
	for _, app := range apps {
		if f.b.Apps == nil {
			f.b.Apps = make(map[string]*App)
		}
		f.b.Apps[app.Guid] = &app
	}

	return nil
}

func (f *fakeDb) UpsertSpace(s Space) error {
	if f.b.Spaces == nil {
		f.b.Spaces = make(map[string]*Space)
	}
	f.b.Spaces[s.Guid] = &s
	return nil
}

func (f *fakeDb) UpsertOrg(o Org) error {
	if f.b.Orgs == nil {
		f.b.Orgs = make(map[string]*Org)
	}
	f.b.Orgs[o.Guid] = &o
	return nil
}

func (f *fakeDb) AcceptReconcileJob(olderThan time.Duration) (ReconcileJob, bool) {
	if f.job == nil {
		return ReconcileJob{}, false
	}
	return *f.job, true
}

type backend struct {
	Orgs   map[string]*Org
	Spaces map[string]*Space
	Apps   map[string]*App
}
