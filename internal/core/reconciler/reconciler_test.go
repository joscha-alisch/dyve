package reconciler

import (
	"errors"
	"github.com/google/go-cmp/cmp"
	"github.com/joscha-alisch/dyve/internal/core/database"
	"github.com/joscha-alisch/dyve/internal/core/provider"
	recon "github.com/joscha-alisch/dyve/internal/reconciliation"
	"github.com/joscha-alisch/dyve/pkg/provider/sdk"
	"testing"
	"time"
)

func TestName(t *testing.T) {
	tests := []struct {
		desc           string
		job            recon.Job
		before         map[string][]sdk.App
		providerId     string
		providerApps   []sdk.App
		providerErr    error
		after          map[string][]sdk.App
		expectedErr    error
		expectedWorked bool
	}{
		{"adds apps", recon.Job{
			Type: database.ReconcileAppProvider,
			Guid: "app-provider",
		}, nil, "app-provider", []sdk.App{
			{Id: "app-a", Name: "app-a"},
			{Id: "app-b", Name: "app-b"},
		}, nil, map[string][]sdk.App{
			"app-provider": {
				{Id: "app-a", Name: "app-a"},
				{Id: "app-b", Name: "app-b"},
			},
		}, nil, true},
		{"removes apps if provider not found", recon.Job{
			Type: database.ReconcileAppProvider,
			Guid: "not-exist",
		}, map[string][]sdk.App{
			"not-exist": {
				{Id: "app-a", Name: "app-a"},
				{Id: "app-b", Name: "app-b"},
			},
		}, "", nil, nil, map[string][]sdk.App{},
			nil, true},
	}

	for _, test := range tests {
		t.Run(test.desc, func(tt *testing.T) {
			db := &fakeDb{job: test.job, apps: test.before}
			r := NewReconciler(db, &fakeManager{
				test.providerId,
				&fakeProvider{
					apps: test.providerApps,
					err:  test.providerErr,
				},
			})
			worked, err := r.Run()
			if !errors.Is(err, test.expectedErr) {
				tt.Errorf("\nwanted err %v\n   got err %v", test.expectedErr, err)
			}

			if worked != test.expectedWorked {
				tt.Errorf("\nwanted worked: %v\n   got worked: %v", test.expectedWorked, worked)
			}

			if !cmp.Equal(test.after, db.apps) {
				tt.Errorf("\nstate diff: \n%s\n", cmp.Diff(test.after, db.apps))
			}
		})
	}

}

type fakeManager struct {
	providerId string
	provider   sdk.AppProvider
}

func (f *fakeManager) AddAppProvider(id string, p sdk.AppProvider) error {
	panic("implement me")
}

func (f *fakeManager) GetAppProvider(id string) (sdk.AppProvider, error) {
	if id == f.providerId {
		return f.provider, nil
	}
	return nil, provider.ErrNotFound
}

type fakeProvider struct {
	apps []sdk.App
	err  error
}

func (f fakeProvider) ListApps() ([]sdk.App, error) {
	return f.apps, f.err
}

func (f fakeProvider) GetApp(id string) (sdk.App, error) {
	panic("implement me")
}

type fakeDb struct {
	job  recon.Job
	apps map[string][]sdk.App
}

func (f *fakeDb) DeleteAppProvider(providerId string) error {
	delete(f.apps, providerId)
	return nil
}

func (f *fakeDb) UpdateApps(providerId string, apps []sdk.App) error {
	if f.apps == nil {
		f.apps = make(map[string][]sdk.App)
	}

	f.apps[providerId] = apps
	return nil
}

func (f *fakeDb) AcceptReconcileJob(olderThan time.Duration) (recon.Job, bool) {
	return f.job, true
}

func (f *fakeDb) ListAppsPaginated(perPage int, page int) (sdk.AppPage, error) {
	panic("implement me")
}
