package routing

import (
	"errors"
	"github.com/google/go-cmp/cmp"
	"github.com/joscha-alisch/dyve/internal/core/fakes/db"
	"github.com/joscha-alisch/dyve/pkg/provider/sdk"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
)

var someRouting = sdk.AppRouting{
	Routes: sdk.AppRoutes{
		{
			Host:    "host",
			Path:    "path",
			AppPort: 123,
		},
	},
}
var someErr = errors.New("some error")

type DecodableFunc func(target interface{}) error

func (f DecodableFunc) Decode(dec interface{}) error {
	return f(dec)
}

func TestService_GetRouting(t *testing.T) {
	tests := []struct {
		desc        string
		app         string
		db          *db.Database
		recorded    []db.DatabaseRecord
		expected    sdk.AppRouting
		expectedErr error
	}{
		{
			desc: "gets routes",
			app:  "app-a",
			db: &db.Database{
				Return: func(target interface{}) {
					*(target.(*routeData)) = routeData{
						Id:        "app-a",
						RouteData: someRouting,
					}
				},
			},
			expected: someRouting,
			recorded: []db.DatabaseRecord{{
				Collection: "routing",
				Id:         "app-a",
			}},
		},
		{
			desc: "error while getting group",
			app:  "a",
			db: &db.Database{
				Err: someErr,
			},
			expected:    sdk.AppRouting{},
			expectedErr: someErr,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(tt *testing.T) {
			recorder := &db.DatabaseRecorder{}
			test.db.Recorder = recorder
			s := NewService(test.db)
			res, err := s.GetRoutes(test.app)
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

func TestService_UpdateRouting(t *testing.T) {
	tests := []struct {
		desc        string
		app         string
		routes      sdk.AppRouting
		db          *db.Database
		recorded    []db.DatabaseRecord
		expectedErr error
	}{
		{
			desc:   "updates routes",
			app:    "app-a",
			routes: someRouting,
			db:     &db.Database{},
			recorded: []db.DatabaseRecord{{
				Collection:      "routing",
				Filter:          bson.M{"id": "app-a"},
				CreateIfMissing: true,
				Update: routeData{
					Id:        "app-a",
					RouteData: someRouting,
				},
			}},
		},
		{
			desc: "error while updating routes",
			app:  "a",
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
			err := s.UpdateRoutes(test.app, test.routes)
			if !errors.Is(err, test.expectedErr) {
				tt.Errorf("errors mismatch: %s\n", cmp.Diff(test.expectedErr, err))
			}

			if !cmp.Equal(test.recorded, recorder.Records) {
				tt.Errorf("recorded mismatch: %s\n", cmp.Diff(test.recorded, recorder.Records))
			}
		})
	}
}
