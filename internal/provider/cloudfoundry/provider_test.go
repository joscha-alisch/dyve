package cloudfoundry

import (
	"errors"
	"github.com/google/go-cmp/cmp"
	"github.com/joscha-alisch/dyve/pkg/provider/sdk"
	"testing"
)

func TestListApps(t *testing.T) {
	tests := []struct {
		desc     string
		db       Database
		expected []sdk.App
	}{
		{"lists apps", &fakeDb{b: backend{Apps: map[string]*App{
			"app-guid-a": {AppInfo: AppInfo{Guid: "app-guid-a", Name: "app-name-a"}},
			"app-guid-b": {AppInfo: AppInfo{Guid: "app-guid-b", Name: "app-name-b"}},
		}}}, []sdk.App{
			{Id: "app-guid-a", Name: "app-name-a"},
			{Id: "app-guid-b", Name: "app-name-b"},
		}},
	}

	for _, test := range tests {
		t.Run(test.desc, func(tt *testing.T) {
			p := NewAppProvider(test.db)
			apps, _ := p.ListApps()

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
		{"gets app", &fakeDb{b: backend{Apps: map[string]*App{
			"app-guid-a": {AppInfo: AppInfo{Guid: "app-guid-a", Name: "app-name-a"}},
			"app-guid-b": {AppInfo: AppInfo{Guid: "app-guid-b", Name: "app-name-b"}},
		}}}, "app-guid-a", sdk.App{Id: "app-guid-a", Name: "app-name-a"}, nil},
	}

	for _, test := range tests {
		t.Run(test.desc, func(tt *testing.T) {
			p := NewAppProvider(test.db)
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