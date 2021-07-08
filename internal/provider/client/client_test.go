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
			{"id-a", "name-a"},
			{"id-b", "name-b"},
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

func TestListAppsPaged(t *testing.T) {
	tests := []struct {
		desc        string
		perPage     int
		page        int
		appPage    sdk.AppPage
		expectedErr error
	}{
		{desc: "returns apps", perPage: 2, page: 4, appPage: sdk.AppPage{
			TotalResults: 10,
			TotalPages:   10,
			PerPage:      2,
			Page:         4,
			Apps:         []sdk.App{
				{"id-a", "name-a"},
				{"id-b", "name-b"},
			},
		}},
	}

	for _, test := range tests {
		t.Run(test.desc, func(tt *testing.T) {
			f := &fakeAppProvider{
				appPage: test.appPage,
			}
			handler := sdk.NewAppProviderHandler(f)
			s := httptest.NewServer(handler)
			defer s.Close()

			c := NewAppProviderClient(s.URL, nil)

			page, err := c.ListAppsPaged(test.perPage, test.page)
			if !errors.Is(err, test.expectedErr) {
				tt.Errorf("\nwanted err: %v\ngot: %v", test.expectedErr, err)
			}

			if !cmp.Equal(test.appPage, page) {
				tt.Errorf("\ndiff between apps\n%s\n", cmp.Diff(test.appPage, page))
			}
		})
	}

}

type fakeAppProvider struct {
	apps []sdk.App
	app  sdk.App
	appPage sdk.AppPage
}

func (f fakeAppProvider) ListAppsPaged(perPage int, page int) (sdk.AppPage, error) {
	if f.appPage.PerPage == perPage && f.appPage.Page == page {
		return f.appPage, nil
	}
	return sdk.AppPage{}, errors.New("something went wrong")
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

func (f fakeAppProvider) Search(term string, limit int) ([]sdk.AppSearchResult, error) {
	panic("implement me")
}
