package cloudfoundry

import (
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
