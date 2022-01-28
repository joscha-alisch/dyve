package teams

import (
	"errors"
	"github.com/google/go-cmp/cmp"
	"github.com/joscha-alisch/dyve/internal/core/database"
	"github.com/joscha-alisch/dyve/internal/core/fakes/db"
	"github.com/joscha-alisch/dyve/pkg/provider/sdk"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
)

var someTeam = Team{
	Id: "team-a",
	TeamSettings: TeamSettings{
		Name:        "team-name",
		Description: "team desc",
		Access: AccessGroups{
			Admin:  []string{"a"},
			Member: []string{"b"},
			Viewer: []string{"c"},
		},
	},
}
var somePagination = sdk.Pagination{
	TotalResults: 123,
	TotalPages:   51234,
	PerPage:      5123,
	Page:         5123,
}
var someErr = errors.New("some error")

func TestService_GetTeam(t *testing.T) {
	tests := []struct {
		desc        string
		team        string
		db          *db.RecordingDatabase
		recorded    []db.DatabaseRecord
		expected    Team
		expectedErr error
	}{
		{
			desc: "gets team",
			team: "team-a",
			db: &db.RecordingDatabase{
				Return: func(target interface{}) {
					*(target.(*Team)) = someTeam
				},
			},
			expected: someTeam,
			recorded: []db.DatabaseRecord{{
				Collection: "teams",
				Id:         "team-a",
			}},
		},
		{
			desc: "error while getting team",
			team: "team-a",
			db: &db.RecordingDatabase{
				Err: someErr,
			},
			expected:    Team{},
			expectedErr: someErr,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(tt *testing.T) {
			recorder := &db.DatabaseRecorder{}
			test.db.Recorder = recorder
			s := NewService(test.db)
			res, err := s.GetTeam(test.team)
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

func TestService_UpdateTeam(t *testing.T) {
	tests := []struct {
		desc        string
		team        string
		data        TeamSettings
		db          *db.RecordingDatabase
		recorded    []db.DatabaseRecord
		expectedErr error
	}{
		{
			desc: "updates team",
			team: "team-a",
			data: TeamSettings{Name: "new-name"},
			db:   &db.RecordingDatabase{},
			recorded: []db.DatabaseRecord{{
				Collection:      "teams",
				Id:              "team-a",
				CreateIfMissing: false,
				Update:          Team{Id: "team-a", TeamSettings: TeamSettings{Name: "new-name"}},
			}},
		},
		{
			desc: "error while updating team",
			team: "team-a",
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
			err := s.UpdateTeam(test.team, test.data)
			if !errors.Is(err, test.expectedErr) {
				tt.Errorf("errors mismatch: %s\n", cmp.Diff(test.expectedErr, err))
			}

			if !cmp.Equal(test.recorded, recorder.Records) {
				tt.Errorf("recorded mismatch: %s\n", cmp.Diff(test.recorded, recorder.Records))
			}
		})
	}
}

func TestService_DeleteTeam(t *testing.T) {
	tests := []struct {
		desc        string
		team        string
		db          *db.RecordingDatabase
		recorded    []db.DatabaseRecord
		expected    Team
		expectedErr error
	}{
		{
			desc:     "deletes team",
			team:     "team-a",
			db:       &db.RecordingDatabase{},
			expected: someTeam,
			recorded: []db.DatabaseRecord{{
				Collection: "teams",
				Id:         "team-a",
			}},
		},
		{
			desc: "error while deleting team",
			team: "team-a",
			db: &db.RecordingDatabase{
				Err: someErr,
			},
			expected:    Team{},
			expectedErr: someErr,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(tt *testing.T) {
			recorder := &db.DatabaseRecorder{}
			test.db.Recorder = recorder
			s := NewService(test.db)
			err := s.DeleteTeam(test.team)
			if !errors.Is(err, test.expectedErr) {
				tt.Errorf("errors mismatch: %s\n", cmp.Diff(test.expectedErr, err))
			}

			if !cmp.Equal(test.recorded, recorder.Records) {
				tt.Errorf("recorded mismatch: %s\n", cmp.Diff(test.recorded, recorder.Records))
			}

		})
	}
}

func TestService_CreateTeam(t *testing.T) {
	tests := []struct {
		desc        string
		team        string
		data        TeamSettings
		db          *db.RecordingDatabase
		recorded    []db.DatabaseRecord
		expectedErr error
	}{
		{
			desc: "creates team",
			team: "team-a",
			data: TeamSettings{Name: "new-name"},
			db:   &db.RecordingDatabase{},
			recorded: []db.DatabaseRecord{{
				Collection: "teams",
				Filter:     bson.M{"id": "team-a"},
				Data:       Team{Id: "team-a", TeamSettings: TeamSettings{Name: "new-name"}},
			}},
		},
		{
			desc: "error while creating team",
			team: "team-a",
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
			err := s.CreateTeam(test.team, test.data)
			if !errors.Is(err, test.expectedErr) {
				tt.Errorf("errors mismatch: %s\n", cmp.Diff(test.expectedErr, err))
			}

			if !cmp.Equal(test.recorded, recorder.Records) {
				tt.Errorf("recorded mismatch: %s\n", cmp.Diff(test.recorded, recorder.Records))
			}
		})
	}
}

func TestService_ListTeamsPaginated(t *testing.T) {
	tests := []struct {
		desc        string
		perPage     int
		page        int
		db          *db.RecordingDatabase
		recorded    []db.DatabaseRecord
		expected    TeamPage
		expectedErr error
	}{
		{
			desc:    "lists teams",
			perPage: 5,
			page:    2,
			db: &db.RecordingDatabase{
				ReturnPagination: func(pagination *sdk.Pagination) {
					*pagination = somePagination
				},
				ReturnEach: func(each func(decodable database.Decodable) error) {
					_ = each(database.DecodableFunc(func(target interface{}) error {
						*target.(*Team) = someTeam
						return nil
					}))
				},
			},
			expected: TeamPage{
				Pagination: somePagination,
				Teams:      []Team{someTeam},
			},
			recorded: []db.DatabaseRecord{{
				Collection: "teams",
				PerPage:    5,
				Page:       2,
			}},
		},
		{
			desc:    "error while listing teams",
			perPage: 5,
			page:    2,
			db: &db.RecordingDatabase{
				Err: someErr,
			},
			expected:    TeamPage{},
			expectedErr: someErr,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(tt *testing.T) {
			recorder := &db.DatabaseRecorder{}
			test.db.Recorder = recorder
			s := NewService(test.db)
			res, err := s.ListTeamsPaginated(test.perPage, test.page)
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

func TestService_TeamsForGroups(t *testing.T) {
	tests := []struct {
		desc        string
		groups      []string
		db          *db.RecordingDatabase
		recorded    []db.DatabaseRecord
		expected    ByAccess
		expectedErr error
	}{
		{
			desc:   "lists teams",
			groups: []string{"a"},
			db: &db.RecordingDatabase{
				ReturnEach: func(each func(decodable database.Decodable) error) {
					teams := []Team{
						{Id: "a", TeamSettings: TeamSettings{Access: AccessGroups{Admin: []string{"a"}}}},
						{Id: "b", TeamSettings: TeamSettings{Access: AccessGroups{Viewer: []string{"a"}}}},
						{Id: "c", TeamSettings: TeamSettings{Access: AccessGroups{Member: []string{"a"}}}},
					}
					for _, team := range teams {
						_ = each(database.DecodableFunc(func(target interface{}) error {
							*target.(*Team) = team
							return nil
						}))
					}

				},
			},
			expected: ByAccess{
				Admin: []Team{
					{Id: "a", TeamSettings: TeamSettings{Access: AccessGroups{Admin: []string{"a"}}}},
				},
				Member: []Team{
					{Id: "c", TeamSettings: TeamSettings{Access: AccessGroups{Member: []string{"a"}}}},
				},
				Viewer: []Team{
					{Id: "b", TeamSettings: TeamSettings{Access: AccessGroups{Viewer: []string{"a"}}}},
				},
			},
			recorded: []db.DatabaseRecord{{
				Collection: "teams",
				Filter: bson.M{
					"$or": []bson.M{
						{"access.admin": bson.M{"$in": []string{"a"}}},
						{"access.member": bson.M{"$in": []string{"a"}}},
						{"access.viewer": bson.M{"$in": []string{"a"}}},
					},
				},
			}},
		},
		{
			desc:   "error while listing teams",
			groups: []string{"a"},
			db: &db.RecordingDatabase{
				Err: someErr,
			},
			expected:    ByAccess{},
			expectedErr: someErr,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(tt *testing.T) {
			recorder := &db.DatabaseRecorder{}
			test.db.Recorder = recorder
			s := NewService(test.db)
			res, err := s.TeamsForGroups(test.groups)
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
