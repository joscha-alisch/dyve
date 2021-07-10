package provider

import (
	"errors"
	"github.com/google/go-cmp/cmp"
	recon "github.com/joscha-alisch/dyve/internal/reconciliation"
	"github.com/joscha-alisch/dyve/pkg/provider/sdk"
	"testing"
	"time"
)

var fakeA = &fakeAppProvider{}

func TestGetAppProvider(t *testing.T) {
	tests := []struct {
		desc        string
		setup       func(m Manager)
		id          string
		expected    sdk.AppProvider
		expectedErr error
	}{
		{"returns provider", func(m Manager) {
			_ = m.AddAppProvider("test-id", fakeA)
		}, "test-id", fakeA, nil},
		{"returns not found", func(m Manager) {},
			"test-id", nil, ErrNotFound},
	}

	for _, test := range tests {
		t.Run(test.desc, func(tt *testing.T) {
			m := NewManager(nil)
			test.setup(m)

			res, err := m.GetAppProvider("test-id")

			if !errors.Is(err, test.expectedErr) {
				tt.Errorf("\nwanted error %v\n   got error %v", test.expectedErr, err)

			}
			if res != test.expected {
				tt.Errorf("\nwanted %v\n   got %v", test.expected, res)
			}
		})
	}

}

func TestAddAppProvider(t *testing.T) {
	tests := []struct {
		desc                string
		setup               func(m Manager)
		id                  string
		provider            sdk.AppProvider
		expectedErr         error
		expectedDbProviders []string
	}{
		{"adds provider to db", func(m Manager) {}, "test-id", fakeA, nil, []string{"test-id"}},
		{"returns already exists", func(m Manager) {
			_ = m.AddAppProvider("test-id", fakeA)
		}, "test-id", fakeA, ErrExists, []string{"test-id"}},
		{"returns provider nil", func(m Manager) {
		}, "test-id", nil, ErrNil, nil},
	}

	for _, test := range tests {
		t.Run(test.desc, func(tt *testing.T) {
			db := &fakeDb{}
			m := NewManager(db)
			test.setup(m)

			err := m.AddAppProvider(test.id, test.provider)

			if !errors.Is(err, test.expectedErr) {
				tt.Errorf("\nwanted %v\n   got %v", test.expectedErr, err)
			}

			if !cmp.Equal(test.expectedDbProviders, db.added) {
				tt.Errorf("\ndiff db providers: \n%s\n", cmp.Diff(test.expectedDbProviders, db.added))
			}
		})
	}

}

type fakeAppProvider struct {
}

func (f *fakeAppProvider) ListApps() ([]sdk.App, error) {
	panic("implement me")
}

func (f *fakeAppProvider) GetApp(id string) (sdk.App, error) {
	panic("implement me")
}

type fakeDb struct {
	added []string
}

func (f *fakeDb) AddAppProvider(providerId string) error {
	f.added = append(f.added, providerId)
	return nil
}

func (f *fakeDb) AcceptReconcileJob(olderThan time.Duration) (recon.Job, bool) {
	panic("implement me")
}

func (f *fakeDb) ListAppsPaginated(perPage int, page int) (sdk.AppPage, error) {
	panic("implement me")
}

func (f *fakeDb) DeleteAppProvider(providerId string) error {
	panic("implement me")
}

func (f *fakeDb) UpdateApps(providerId string, apps []sdk.App) error {
	panic("implement me")
}