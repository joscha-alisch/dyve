package client

import (
	"errors"
	"github.com/google/go-cmp/cmp"
	"github.com/joscha-alisch/dyve/pkg/provider/sdk"
	"net/http/httptest"
	"testing"
)

func TestGetApp(t *testing.T) {
	tests := []struct {
		desc        string
		id          string
		app         sdk.App
		expectedErr error
	}{
		{desc: "returns apps", id: "id-a", app: sdk.App{Id: "id-a", Name: "name-a"}},
		{desc: "returns not found err", id: "not-exist", app: sdk.App{}, expectedErr: sdk.ErrNotFound},
	}

	for _, test := range tests {
		t.Run(test.desc, func(tt *testing.T) {
			f := &fakeAppProvider{
				app: test.app,
			}
			handler := sdk.NewAppProviderHandler(f)
			s := httptest.NewServer(handler)
			defer s.Close()

			c := NewAppProviderClient(s.URL, nil)

			apps, err := c.GetApp(test.id)
			if !errors.Is(err, test.expectedErr) {
				tt.Errorf("\nwanted err: %v\ngot: %v", test.expectedErr, err)
			}

			if !cmp.Equal(test.app, apps) {
				tt.Errorf("\ndiff between apps\n%s\n", cmp.Diff(test.app, apps))
			}
		})
	}

}

func TestListApps(t *testing.T) {
	tests := []struct {
		desc        string
		apps        []sdk.App
		expectedErr error
	}{
		{desc: "returns apps", apps: []sdk.App{
			{Id: "id-a", Name: "name-a"},
			{Id: "id-b", Name: "name-b"},
		}},
	}

	for _, test := range tests {
		t.Run(test.desc, func(tt *testing.T) {
			f := &fakeAppProvider{
				apps: test.apps,
			}
			handler := sdk.NewAppProviderHandler(f)
			s := httptest.NewServer(handler)
			defer s.Close()

			c := NewAppProviderClient(s.URL, nil)

			apps, err := c.ListApps()
			if !errors.Is(err, test.expectedErr) {
				tt.Errorf("\nwanted err: %v\ngot: %v", test.expectedErr, err)
			}

			if !cmp.Equal(test.apps, apps) {
				tt.Errorf("\ndiff between apps\n%s\n", cmp.Diff(test.apps, apps))
			}
		})
	}
}

type fakeAppProvider struct {
	apps    []sdk.App
	app     sdk.App
	appPage sdk.AppPage
}

func (f fakeAppProvider) ListApps() ([]sdk.App, error) {
	return f.apps, nil
}

func (f fakeAppProvider) GetApp(id string) (sdk.App, error) {
	if f.app.Id == id {
		return f.app, nil
	}
	return sdk.App{}, sdk.ErrNotFound
}
