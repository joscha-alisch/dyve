package cloudfoundry

import (
	"errors"
	"github.com/google/go-cmp/cmp"
	"github.com/joscha-alisch/dyve/pkg/provider/sdk"
	"testing"
)

var someErr = errors.New("error")

func TestListApps(t *testing.T) {
	tests := []struct {
		desc        string
		db          Database
		expected    []sdk.App
		expectedErr error
	}{
		{desc: "lists apps", db: &fakeDb{b: backend{Apps: map[string]*App{
			"app-guid-a": {AppInfo: AppInfo{Guid: "app-guid-a", Name: "app-name-a"}},
			"app-guid-b": {AppInfo: AppInfo{Guid: "app-guid-b", Name: "app-name-b"}},
		}}}, expected: []sdk.App{
			{Id: "app-guid-a", Name: "app-name-a"},
			{Id: "app-guid-b", Name: "app-name-b"},
		}},
		{desc: "returns error", db: &fakeDb{err: someErr}, expected: nil, expectedErr: someErr},
	}

	for _, test := range tests {
		t.Run(test.desc, func(tt *testing.T) {
			p := NewProvider(test.db, nil)
			apps, err := p.ListApps()

			if err != test.expectedErr {
				tt.Errorf("\ndiff between errors: \n%s\n", cmp.Diff(test.expectedErr, err))
			}

			if !cmp.Equal(test.expected, apps) {
				tt.Errorf("\ndiff between returned apps: \n%s\n", cmp.Diff(test.expected, apps))
			}
		})
	}
}

func TestGetApp(t *testing.T) {
	tests := []struct {
		desc        string
		db          Database
		id          string
		expected    sdk.App
		expectedErr error
	}{
		{desc: "gets app", db: &fakeDb{b: backend{Apps: map[string]*App{
			"app-guid-a": {AppInfo: AppInfo{Guid: "app-guid-a", Name: "app-name-a"}},
			"app-guid-b": {AppInfo: AppInfo{Guid: "app-guid-b", Name: "app-name-b"}},
		}}}, id: "app-guid-a", expected: sdk.App{Id: "app-guid-a", Name: "app-name-a"}},
		{desc: "returns error", db: &fakeDb{err: someErr}, expected: sdk.App{}, expectedErr: someErr},
	}

	for _, test := range tests {
		t.Run(test.desc, func(tt *testing.T) {
			p := NewProvider(test.db, nil)
			app, err := p.GetApp(test.id)
			if !errors.Is(test.expectedErr, err) {
				tt.Errorf("\nwanted error: \n%s, got %s\n", test.expectedErr.Error(), err.Error())
			}
			if !cmp.Equal(test.expected, app) {
				tt.Errorf("\ndiff between returned app: \n%s\n", cmp.Diff(test.expected, app))
			}
		})
	}
}

func TestGetAppInstances(t *testing.T) {
	tests := []struct {
		desc        string
		db          Database
		cf          API
		id          string
		expected    sdk.AppInstances
		expectedErr error
	}{
		{desc: "gets app instances cached", db: &fakeDb{b: backend{Cache: map[string]interface{}{
			"app-guid-a/instances": sdk.AppInstances{{
				State: "STOPPED",
				Since: someTime,
			}},
		}}}, id: "app-guid-a", expected: sdk.AppInstances{{
			State: "STOPPED",
			Since: someTime,
		}}},
		{desc: "returns error", db: &fakeDb{err: someErr}, expected: nil, expectedErr: someErr},
		{desc: "gets app instances from CF", db: &fakeDb{}, cf: &fakeCf{b: backend{
			AppInstances: map[string]Instances{
				"app-guid-a": {{
					State: "STOPPED",
					Since: someTime,
				}, {
					State: "STARTED",
					Since: someTime,
				}, {
					State: "SOME OTHER UNKNOWN STATE",
					Since: someTime,
				}},
			}},
		}, id: "app-guid-a", expected: sdk.AppInstances{{
			State: "stopped",
			Since: someTime,
		}, {
			State: "running",
			Since: someTime,
		}, {
			State: "unknown",
			Since: someTime,
		}}},
	}

	for _, test := range tests {
		t.Run(test.desc, func(tt *testing.T) {
			p := NewProvider(test.db, test.cf)
			app, err := p.GetAppInstances(test.id)
			if !errors.Is(test.expectedErr, err) {
				tt.Errorf("\nwanted error: \n%s, got %s\n", test.expectedErr.Error(), err.Error())
			}
			if !cmp.Equal(test.expected, app) {
				tt.Errorf("\ndiff between returned app: \n%s\n", cmp.Diff(test.expected, app))
			}
		})
	}
}

func TestGetAppRoutes(t *testing.T) {
	tests := []struct {
		desc        string
		db          Database
		cf          API
		id          string
		expected    sdk.AppRouting
		expectedErr error
	}{
		{desc: "gets app routing cached", db: &fakeDb{b: backend{Cache: map[string]interface{}{
			"app-guid-a/routing": sdk.AppRouting{Routes: []sdk.AppRoute{
				{
					Host:    "host",
					Path:    "path",
					AppPort: 231,
				},
			}},
		}}}, id: "app-guid-a", expected: sdk.AppRouting{Routes: []sdk.AppRoute{
			{
				Host:    "host",
				Path:    "path",
				AppPort: 231,
			},
		}}},
		{desc: "returns error", db: &fakeDb{err: someErr}, expected: sdk.AppRouting{}, expectedErr: someErr},
		{desc: "gets app routing from CF", db: &fakeDb{}, cf: &fakeCf{b: backend{
			AppRoutes: map[string]Routes{
				"app-guid-a": {{
					Host: "host",
					Path: "path",
					Port: 4223,
				}}}},
		}, id: "app-guid-a", expected: sdk.AppRouting{Routes: sdk.AppRoutes{
			{
				Host:    "host",
				Path:    "path",
				AppPort: 4223,
			},
		}}},
	}

	for _, test := range tests {
		t.Run(test.desc, func(tt *testing.T) {
			p := NewProvider(test.db, test.cf)
			app, err := p.GetAppRouting(test.id)
			if !errors.Is(test.expectedErr, err) {
				tt.Errorf("\nwanted error: \n%s, got %s\n", test.expectedErr.Error(), err.Error())
			}
			if !cmp.Equal(test.expected, app) {
				tt.Errorf("\ndiff between returned app: \n%s\n", cmp.Diff(test.expected, app))
			}
		})
	}
}
