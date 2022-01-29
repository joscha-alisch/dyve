package client

import (
	"errors"
	"github.com/google/go-cmp/cmp"
	"github.com/joscha-alisch/dyve/pkg/provider/sdk"
	"net/http/httptest"
	"testing"
)

func TestGetAppInstances(t *testing.T) {
	tests := []struct {
		desc        string
		id          string
		state       fakeInstancesProvider
		expectedErr error
	}{
		{desc: "returns instances", id: "id-a", state: fakeInstancesProvider{
			instances: sdk.AppInstances{{
				State: "stopped",
			}},
		}},
		{desc: "returns not found err", id: "not-exist", state: fakeInstancesProvider{
			err: sdk.ErrNotFound,
		}, expectedErr: sdk.ErrNotFound},
	}

	for _, test := range tests {
		t.Run(test.desc, func(tt *testing.T) {
			handler := sdk.NewAppInstancesProviderHandler(&test.state)
			s := httptest.NewServer(handler)
			defer s.Close()

			c := NewInstancesProviderClient(s.URL, nil)

			instances, err := c.GetAppInstances(test.id)
			if !errors.Is(err, test.expectedErr) {
				tt.Errorf("\nwanted err: %v\ngot: %v", test.expectedErr, err)
			}

			if test.state.recordedId != test.id {
				tt.Errorf("\nwanted id: %s\ngot: %s", test.id, test.state.recordedId)
			}

			if !cmp.Equal(test.state.instances, instances) {
				tt.Errorf("\ndiff between instances\n%s\n", cmp.Diff(test.state.instances, instances))
			}
		})
	}

}

type fakeInstancesProvider struct {
	err        error
	instances  sdk.AppInstances
	recordedId string
}

func (f *fakeInstancesProvider) GetAppInstances(id string) (sdk.AppInstances, error) {
	f.recordedId = id
	if f.err != nil {
		return sdk.AppInstances{}, f.err
	}

	return f.instances, nil
}
