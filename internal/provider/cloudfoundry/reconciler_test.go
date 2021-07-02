package cloudfoundry

import (
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"testing"
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
		{
			desc: "no work to be done",
			db: fakeDb{
				job: nil,
			},
			cf:    fakeCf{},
			sleep: true,
		},
		{
			desc: "handle error org not found",
			db: fakeDb{
				job: &ReconcileJob{Type: ReconcileOrg, Guid: "not_exist"},
			},
			cf: fakeCf{},
			err: &errReconcileFailed{
				Err: errNotFound,
				Job: ReconcileJob{Type: ReconcileOrg, Guid: "not_exist"},
			},
		},
		{
			desc: "handle error space not found",
			db: fakeDb{
				job: &ReconcileJob{Type: ReconcileSpace, Guid: "not_exist"},
			},
			cf: fakeCf{},
			err: &errReconcileFailed{
				Err: errNotFound,
				Job: ReconcileJob{Type: ReconcileSpace, Guid: "not_exist"},
			},
		},
		{
			desc: "handle error app not found",
			db: fakeDb{
				job: &ReconcileJob{Type: ReconcileApp, Guid: "not_exist"},
			},
			cf: fakeCf{},
			err: &errReconcileFailed{
				Err: errNotFound,
				Job: ReconcileJob{Type: ReconcileApp, Guid: "not_exist"},
			},
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

func (f *fakeCf) GetApp(guid string) (App, error) {
	if f.b.Apps[guid] == nil {
		return App{}, errNotFound
	}
	return *f.b.Apps[guid], nil
}

func (f *fakeCf) GetSpace(guid string) (Space, error) {
	if f.b.Spaces[guid] == nil {
		return Space{}, errNotFound
	}
	return *f.b.Spaces[guid], nil
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

func (f *fakeDb) FetchReconcileJob() (ReconcileJob, bool) {
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
