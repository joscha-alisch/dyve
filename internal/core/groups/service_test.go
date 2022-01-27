package groups

import (
	"errors"
	"github.com/google/go-cmp/cmp"
	"github.com/joscha-alisch/dyve/internal/core/database"
	"github.com/joscha-alisch/dyve/internal/core/fakes"
	"github.com/joscha-alisch/dyve/internal/core/fakes/db"
	"github.com/joscha-alisch/dyve/internal/core/provider"
	"github.com/joscha-alisch/dyve/pkg/provider/sdk"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
)

var someGroup = GroupWithProvider{
	Provider: "group-provider",
	Group: sdk.Group{
		Id:   "group-a",
		Name: "name",
	},
}
var somePagination = sdk.Pagination{
	TotalResults: 123,
	TotalPages:   123,
	PerPage:      123,
	Page:         123,
}
var someErr = errors.New("some error")

type DecodableFunc func(target interface{}) error

func (f DecodableFunc) Decode(dec interface{}) error {
	return f(dec)
}

func TestService_GetGroup(t *testing.T) {
	tests := []struct {
		desc        string
		id          string
		db          *db.RecordingDatabase
		recorded    []db.DatabaseRecord
		expected    sdk.Group
		expectedErr error
	}{
		{
			desc: "gets group",
			id:   "group-a",
			db: &db.RecordingDatabase{
				Return: func(target interface{}) {
					*(target.(*sdk.Group)) = someGroup.Group
				},
			},
			expected: someGroup.Group,
			recorded: []db.DatabaseRecord{{
				Collection: "groups",
				Id:         "group-a",
			}},
		},
		{
			desc: "error while getting group",
			id:   "a",
			db: &db.RecordingDatabase{
				Err: someErr,
			},
			expected:    sdk.Group{},
			expectedErr: someErr,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(tt *testing.T) {
			recorder := &db.DatabaseRecorder{}
			test.db.Recorder = recorder
			s := NewService(test.db, nil)
			res, err := s.GetGroup(test.id)
			if !errors.Is(err, test.expectedErr) {
				tt.Errorf("errors mismatch: %s\n", cmp.Diff(test.expectedErr, err))
			}

			if !cmp.Equal(test.expected, res) {
				tt.Errorf("results mismatch: %s\n", cmp.Diff(test.expected, res))
			}

			if !cmp.Equal(test.recorded, recorder.Records) {
				tt.Errorf("recorded mismatch: %s\n", cmp.Diff(test.recorded, recorder.Records))
			}

		})
	}
}

func TestService_DeleteGroup(t *testing.T) {
	tests := []struct {
		desc        string
		id          string
		db          *db.RecordingDatabase
		recorded    []db.DatabaseRecord
		expectedErr error
	}{
		{
			desc: "deletes group",
			id:   "group-a",
			db:   &db.RecordingDatabase{},
			recorded: []db.DatabaseRecord{{
				Collection: "groups",
				Id:         "group-a",
			}},
		},
		{
			desc: "error while deleting group",
			id:   "a",
			db: &db.RecordingDatabase{
				Err: someErr,
			},
			expectedErr: someErr,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(tt *testing.T) {
			recorder := &db.DatabaseRecorder{}
			test.db.Recorder = recorder
			s := NewService(test.db, nil)
			err := s.DeleteGroup(test.id)
			if !errors.Is(err, test.expectedErr) {
				tt.Errorf("errors mismatch: %s\n", cmp.Diff(test.expectedErr, err))
			}

			if !cmp.Equal(test.recorded, recorder.Records) {
				tt.Errorf("recorded mismatch: %s\n", cmp.Diff(test.recorded, recorder.Records))
			}

		})
	}
}

func TestService_UpdateGroups(t *testing.T) {
	tests := []struct {
		desc        string
		provider    string
		groups      []sdk.Group
		db          *db.RecordingDatabase
		recorded    []db.DatabaseRecord
		expectedErr error
	}{
		{
			desc:     "updates groups",
			provider: "provider-a",
			groups:   []sdk.Group{someGroup.Group},
			db:       &db.RecordingDatabase{},
			recorded: []db.DatabaseRecord{{
				Collection: "groups",
				Provider:   "provider-a",
				Updates: map[string]interface{}{
					someGroup.Group.Id: someGroup.Group,
				},
			}},
		},
		{
			desc:     "error while updating groups",
			provider: "a",
			db: &db.RecordingDatabase{
				Err: someErr,
			},
			expectedErr: someErr,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(tt *testing.T) {
			recorder := &db.DatabaseRecorder{}
			test.db.Recorder = recorder
			s := NewService(test.db, nil)
			err := s.UpdateGroups(test.provider, test.groups)
			if !errors.Is(err, test.expectedErr) {
				tt.Errorf("errors mismatch: %s\n", cmp.Diff(test.expectedErr, err))
			}

			if !cmp.Equal(test.recorded, recorder.Records) {
				tt.Errorf("recorded mismatch: %s\n", cmp.Diff(test.recorded, recorder.Records))
			}

		})
	}
}

func TestService_ListGroupsPaginated(t *testing.T) {
	tests := []struct {
		desc        string
		perPage     int
		page        int
		db          *db.RecordingDatabase
		recorded    []db.DatabaseRecord
		expected    sdk.GroupPage
		expectedErr error
	}{
		{
			desc:    "lists groups",
			perPage: 5,
			page:    2,
			db: &db.RecordingDatabase{
				ReturnEach: func(each func(decodable database.Decodable) error) {
					_ = each(DecodableFunc(func(target interface{}) error {
						*target.(*sdk.Group) = someGroup.Group
						return nil
					}))
				},
				ReturnPagination: func(pagination *sdk.Pagination) {
					*pagination = somePagination
				},
			},
			expected: sdk.GroupPage{Pagination: somePagination, Groups: []sdk.Group{someGroup.Group}},
			recorded: []db.DatabaseRecord{{
				Collection: "groups",
				PerPage:    5,
				Page:       2,
			}},
		},
		{
			desc:    "error while listing groups",
			perPage: 5,
			page:    2,
			db: &db.RecordingDatabase{
				Err: someErr,
			},
			expectedErr: someErr,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(tt *testing.T) {
			recorder := &db.DatabaseRecorder{}
			test.db.Recorder = recorder
			s := NewService(test.db, nil)
			res, err := s.ListGroupsPaginated(test.perPage, test.page)
			if !errors.Is(err, test.expectedErr) {
				tt.Errorf("errors mismatch: %s\n", cmp.Diff(test.expectedErr, err))
			}

			if !cmp.Equal(test.expected, res) {
				tt.Errorf("result mismatch: %s\n", cmp.Diff(test.expected, res))
			}

			if !cmp.Equal(test.recorded, recorder.Records) {
				tt.Errorf("recorded mismatch: %s\n", cmp.Diff(test.recorded, recorder.Records))
			}

		})
	}
}

func TestService_ListGroupsByProvider(t *testing.T) {
	tests := []struct {
		desc        string
		db          *db.RecordingDatabase
		providers   provider.Service
		recorded    []db.DatabaseRecord
		expected    GroupByProviderMap
		expectedErr error
	}{
		{
			desc: "lists groups by providers",
			db: &db.RecordingDatabase{
				ReturnEach: func(each func(decodable database.Decodable) error) {
					_ = each(DecodableFunc(func(target interface{}) error {
						*target.(*GroupWithProvider) = someGroup
						return nil
					}))
				},
			},
			providers: &fakes.ProviderService{GroupProviders: map[string]sdk.GroupProvider{
				"group-provider": nil,
			}},
			expected: GroupByProviderMap{
				"group-provider": ProviderWithGroups{
					Provider: "group-provider",
					Name:     "group-provider",
					Groups:   []sdk.Group{someGroup.Group},
				},
			},
			recorded: []db.DatabaseRecord{{
				Collection: "groups",
				Filter:     bson.M{},
			}},
		},
		{
			desc: "error while listing groups",
			providers: &fakes.ProviderService{GroupProviders: map[string]sdk.GroupProvider{
				"group-provider": nil,
			}},
			db: &db.RecordingDatabase{
				Err: someErr,
			},
			expectedErr: someErr,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(tt *testing.T) {
			recorder := &db.DatabaseRecorder{}
			test.db.Recorder = recorder
			s := NewService(test.db, test.providers)
			res, err := s.ListGroupsByProvider()
			if !errors.Is(err, test.expectedErr) {
				tt.Errorf("errors mismatch: %s\n", cmp.Diff(test.expectedErr, err))
			}

			if !cmp.Equal(test.expected, res) {
				tt.Errorf("result mismatch: %s\n", cmp.Diff(test.expected, res))
			}

			if !cmp.Equal(test.recorded, recorder.Records) {
				tt.Errorf("recorded mismatch: %s\n", cmp.Diff(test.recorded, recorder.Records))
			}

		})
	}
}
