package client

import (
	"errors"
	"github.com/google/go-cmp/cmp"
	"github.com/joscha-alisch/dyve/pkg/provider/sdk"
	"net/http/httptest"
	"testing"
)

func TestGetGroup(t *testing.T) {
	tests := []struct {
		desc        string
		id          string
		state       fakeGroupProvider
		expectedErr error
	}{
		{desc: "returns group", id: "id-a", state: fakeGroupProvider{
			group: sdk.Group{Id: "id-a", Name: "name-a"}},
		},
		{desc: "returns not found err", id: "not-exist", state: fakeGroupProvider{
			err: sdk.ErrNotFound,
		}, expectedErr: sdk.ErrNotFound},
	}

	for _, test := range tests {
		t.Run(test.desc, func(tt *testing.T) {
			handler := sdk.NewGroupProviderHandler(&test.state)
			s := httptest.NewServer(handler)
			defer s.Close()

			c := NewGroupProviderClient(s.URL, nil)

			group, err := c.GetGroup(test.id)
			if !errors.Is(err, test.expectedErr) {
				tt.Errorf("\nwanted err: %v\ngot: %v", test.expectedErr, err)
			}

			if test.state.recordedId != test.id {
				tt.Errorf("\nwanted id: %s\ngot: %s", test.id, test.state.recordedId)
			}

			if !cmp.Equal(test.state.group, group) {
				tt.Errorf("\ndiff between groups\n%s\n", cmp.Diff(test.state.group, group))
			}
		})
	}

}

func TestListGroups(t *testing.T) {
	tests := []struct {
		desc        string
		state       fakeGroupProvider
		expectedErr error
	}{
		{desc: "returns groups", state: fakeGroupProvider{
			groups: []sdk.Group{
				{Id: "id-a", Name: "name-a"},
				{Id: "id-b", Name: "name-b"},
			},
		},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(tt *testing.T) {
			handler := sdk.NewGroupProviderHandler(&test.state)
			s := httptest.NewServer(handler)
			defer s.Close()

			c := NewGroupProviderClient(s.URL, nil)

			groups, err := c.ListGroups()
			if !errors.Is(err, test.expectedErr) {
				tt.Errorf("\nwanted err: %v\ngot: %v", test.expectedErr, err)
			}

			if !cmp.Equal(test.state.groups, groups) {
				tt.Errorf("\ndiff between pipelines\n%s\n", cmp.Diff(test.state.groups, groups))
			}
		})
	}
}

type fakeGroupProvider struct {
	err        error
	groups     []sdk.Group
	group      sdk.Group
	recordedId string
}

func (f *fakeGroupProvider) ListGroups() ([]sdk.Group, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.groups, nil
}

func (f *fakeGroupProvider) GetGroup(id string) (sdk.Group, error) {
	f.recordedId = id
	if f.err != nil {
		return sdk.Group{}, f.err
	}

	return f.group, nil
}
