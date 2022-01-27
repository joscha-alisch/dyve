package instances

import (
	"errors"
	"github.com/google/go-cmp/cmp"
	"github.com/joscha-alisch/dyve/internal/core/fakes/db"
	"github.com/joscha-alisch/dyve/pkg/provider/sdk"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
)

var someInstances = sdk.AppInstances{
	{State: "stopped"},
}
var someErr = errors.New("some error")

type DecodableFunc func(target interface{}) error

func (f DecodableFunc) Decode(dec interface{}) error {
	return f(dec)
}

func TestService_GetInstances(t *testing.T) {
	tests := []struct {
		desc        string
		app         string
		db          *db.RecordingDatabase
		recorded    []db.DatabaseRecord
		expected    sdk.AppInstances
		expectedErr error
	}{
		{
			desc: "gets instances",
			app:  "app-a",
			db: &db.RecordingDatabase{
				Return: func(target interface{}) {
					*(target.(*instancesData)) = instancesData{
						Id:            "app-a",
						InstancesData: someInstances,
					}
				},
			},
			expected: someInstances,
			recorded: []db.DatabaseRecord{{
				Collection: "instances",
				Id:         "app-a",
			}},
		},
		{
			desc: "error while getting group",
			app:  "a",
			db: &db.RecordingDatabase{
				Err: someErr,
			},
			expected:    nil,
			expectedErr: someErr,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(tt *testing.T) {
			recorder := &db.DatabaseRecorder{}
			test.db.Recorder = recorder
			s := NewService(test.db)
			res, err := s.GetInstances(test.app)
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

func TestService_UpdateInstances(t *testing.T) {
	tests := []struct {
		desc        string
		app         string
		instances   sdk.AppInstances
		db          *db.RecordingDatabase
		recorded    []db.DatabaseRecord
		expectedErr error
	}{
		{
			desc:      "updates instances",
			app:       "app-a",
			instances: someInstances,
			db:        &db.RecordingDatabase{},
			recorded: []db.DatabaseRecord{{
				Collection:      "instances",
				Filter:          bson.M{"id": "app-a"},
				CreateIfMissing: true,
				Update: instancesData{
					Id:            "app-a",
					InstancesData: someInstances,
				},
			}},
		},
		{
			desc: "error while updating instances",
			app:  "a",
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
			s := NewService(test.db)
			err := s.UpdateInstances(test.app, test.instances)
			if !errors.Is(err, test.expectedErr) {
				tt.Errorf("errors mismatch: %s\n", cmp.Diff(test.expectedErr, err))
			}

			if !cmp.Equal(test.recorded, recorder.Records) {
				tt.Errorf("recorded mismatch: %s\n", cmp.Diff(test.recorded, recorder.Records))
			}
		})
	}
}
