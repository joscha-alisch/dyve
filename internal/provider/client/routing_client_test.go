package client

import (
	"errors"
	"github.com/google/go-cmp/cmp"
	"github.com/joscha-alisch/dyve/pkg/provider/sdk"
	"net/http/httptest"
	"testing"
)

func TestGetAppRouting(t *testing.T) {
	tests := []struct {
		desc        string
		id          string
		state       fakeRoutingProvider
		expectedErr error
	}{
		{desc: "returns routing", id: "id-a", state: fakeRoutingProvider{
			routing: sdk.AppRouting{
				Routes: []sdk.AppRoute{{
					Host:    "host",
					Path:    "path",
					AppPort: 123,
				}},
			}},
		},
		{desc: "returns not found err", id: "not-exist", state: fakeRoutingProvider{
			err: sdk.ErrNotFound,
		}, expectedErr: sdk.ErrNotFound},
	}

	for _, test := range tests {
		t.Run(test.desc, func(tt *testing.T) {
			handler := sdk.NewAppRoutingProviderHandler(&test.state)
			s := httptest.NewServer(handler)
			defer s.Close()

			c := NewRoutingProviderClient(s.URL, nil)

			routing, err := c.GetAppRouting(test.id)
			if !errors.Is(err, test.expectedErr) {
				tt.Errorf("\nwanted err: %v\ngot: %v", test.expectedErr, err)
			}

			if test.state.recordedId != test.id {
				tt.Errorf("\nwanted id: %s\ngot: %s", test.id, test.state.recordedId)
			}

			if !cmp.Equal(test.state.routing, routing) {
				tt.Errorf("\ndiff between routing\n%s\n", cmp.Diff(test.state.routing, routing))
			}
		})
	}

}

type fakeRoutingProvider struct {
	err        error
	routing    sdk.AppRouting
	recordedId string
}

func (f *fakeRoutingProvider) GetAppRouting(id string) (sdk.AppRouting, error) {
	f.recordedId = id
	if f.err != nil {
		return sdk.AppRouting{}, f.err
	}

	return f.routing, nil
}
