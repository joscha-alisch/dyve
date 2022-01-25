package apps

import (
	"errors"
	"github.com/google/go-cmp/cmp"
	"github.com/joscha-alisch/dyve/internal/core/database"
	"github.com/joscha-alisch/dyve/internal/core/fakes/db"
	"github.com/joscha-alisch/dyve/pkg/provider/sdk"
	"testing"
)

var someApp = App{
	App: sdk.App{
		Id:   "a",
		Name: "name",
	},
	ProviderId: "provider-a",
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

func TestService_GetApp(t *testing.T) {
	tests := []struct {
		desc        string
		id          string
		db          *db.Database
		recorded    []db.DatabaseRecord
		expected    App
		expectedErr error
	}{
		{
			desc: "gets app",
			id:   "a",
			db: &db.Database{
				Return: func(target interface{}) {
					*(target.(*App)) = someApp
				},
			},
			expected: someApp,
			recorded: []db.DatabaseRecord{{
				Collection: "apps",
				Id:         "a",
			}},
		},
		{
			desc: "error while getting app",
			id:   "a",
			db: &db.Database{
				Err: someErr,
			},
			expected:    App{},
			expectedErr: someErr,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(tt *testing.T) {
			recorder := &db.DatabaseRecorder{}
			test.db.Recorder = recorder
			s := NewService(test.db)
			res, err := s.GetApp(test.id)
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

func TestService_ListAppsPaginated(t *testing.T) {
	tests := []struct {
		desc        string
		perPage     int
		page        int
		db          *db.Database
		recorded    []db.DatabaseRecord
		expected    sdk.AppPage
		expectedErr error
	}{
		{
			desc:    "lists apps",
			perPage: 5,
			page:    2,
			db: &db.Database{
				ReturnPagination: func(pagination *sdk.Pagination) {
					*pagination = somePagination
				},
				ReturnEach: func(each func(dec database.Decodable) error) {
					_ = each(DecodableFunc(func(target interface{}) error {
						*(target.(*sdk.App)) = someApp.App
						return nil
					}))
				},
			},
			expected: sdk.AppPage{
				Pagination: somePagination,
				Apps:       []sdk.App{someApp.App},
			},
			recorded: []db.DatabaseRecord{{
				Collection: "apps",
				PerPage:    5,
				Page:       2,
			}},
		},
		{
			desc:    "error while getting app",
			perPage: 5,
			page:    2,
			db: &db.Database{
				Err: someErr,
			},
			expected:    sdk.AppPage{},
			expectedErr: someErr,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(tt *testing.T) {
			recorder := &db.DatabaseRecorder{}
			test.db.Recorder = recorder
			s := NewService(test.db)
			res, err := s.ListAppsPaginated(test.perPage, test.page)
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

func TestService_UpdateApps(t *testing.T) {
	tests := []struct {
		desc        string
		providerId  string
		apps        []sdk.App
		db          *db.Database
		recorded    []db.DatabaseRecord
		expectedErr error
	}{
		{
			desc:       "update apps",
			providerId: "provider-a",
			apps:       []sdk.App{someApp.App},
			db:         &db.Database{},
			recorded: []db.DatabaseRecord{{
				Collection: "apps",
				Provider:   "provider-a",
				Updates:    map[string]interface{}{someApp.Id: someApp.App},
			}},
		},
		{
			desc:       "error while updating apps",
			providerId: "provider-a",
			apps:       []sdk.App{someApp.App},
			db: &db.Database{
				Err: someErr,
			},
			expectedErr: someErr,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(tt *testing.T) {
			recorder := &db.DatabaseRecorder{}
			test.db.Recorder = recorder
			s := NewService(test.db)
			err := s.UpdateApps(test.providerId, test.apps)
			if !errors.Is(err, test.expectedErr) {
				tt.Errorf("errors mismatch: %s\n", cmp.Diff(test.expectedErr, err))
			}

			if !cmp.Equal(test.recorded, recorder.Records) {
				tt.Errorf("recorded mismatch: %s\n", cmp.Diff(test.recorded, recorder.Records))
			}

		})
	}
}

func TestService_UpdateApp(t *testing.T) {
	tests := []struct {
		desc        string
		app         sdk.App
		db          *db.Database
		recorded    []db.DatabaseRecord
		expectedErr error
	}{
		{
			desc: "updates app",
			app:  someApp.App,
			db:   &db.Database{},
			recorded: []db.DatabaseRecord{{
				Collection:      "apps",
				CreateIfMissing: false,
				Id:              someApp.App.Id,
				Update:          someApp.App,
			}},
		},
		{
			desc: "error while updating app",
			app:  someApp.App,
			db: &db.Database{
				Err: someErr,
			},
			expectedErr: someErr,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(tt *testing.T) {
			recorder := &db.DatabaseRecorder{}
			test.db.Recorder = recorder
			s := NewService(test.db)
			err := s.UpdateApp(test.app)
			if !errors.Is(err, test.expectedErr) {
				tt.Errorf("errors mismatch: %s\n", cmp.Diff(test.expectedErr, err))
			}

			if !cmp.Equal(test.recorded, recorder.Records) {
				tt.Errorf("recorded mismatch: %s\n", cmp.Diff(test.recorded, recorder.Records))
			}

		})
	}
}
